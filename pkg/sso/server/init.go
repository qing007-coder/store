package server

import (
	"context"
	"fmt"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/generates"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/server"
	redis_ "github.com/go-oauth2/redis/v4"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
	"log"
	"net/http"
	"store/pkg/config"
	"store/pkg/redis"
	"store/pkg/rules"
	"store/pkg/tools"
	"time"
)

type Server struct {
	ctx      context.Context
	Srv      *server.Server
	rdb      *redis.Client
	db       *gorm.DB
	conf     *config.GlobalConfig
	enforcer *rules.Enforcer
}

func NewServer(r *redis.Client, db *gorm.DB, conf *config.GlobalConfig, e *rules.Enforcer) *Server {
	srv := &Server{
		ctx:      context.Background(),
		rdb:      r,
		db:       db,
		conf:     conf,
		enforcer: e,
	}
	srv.init()
	return srv
}

func (s *Server) init() {
	//manager := manage.NewDefaultManager()  自动初始化一些默认的存储器，比如客户端存储、令牌存储、授权码存储等
	manager := manage.NewManager()
	manager.SetAuthorizeCodeTokenCfg(&manage.Config{
		AccessTokenExp:    time.Hour * 6,
		RefreshTokenExp:   time.Hour * 24 * 2,
		IsGenerateRefresh: true,
	}) // 设置token的配置

	//tokenStore := NewTokeStore(s.rdb, s.conf)
	//manager.MapTokenStorage(tokenStore) // 设置存储token的配置
	manager.MapTokenStorage(redis_.NewRedisStore(s.rdb.GetConfig(s.conf))) // 设置存储token的配置
	manager.MapAccessGenerate(generates.NewJWTAccessGenerate("qing", []byte("密钥miyao"), jwt.SigningMethodHS256))
	manager.MapAuthorizeGenerate(generates.NewAuthorizeGenerate())

	cliStore := NewStore(s.db)
	manager.MapClientStorage(cliStore) // 设置client信息存储

	s.Srv = server.NewServer(server.NewConfig(), manager)
	s.Srv.SetUserAuthorizationHandler(s.UserAuthorizeHandler)
	s.Srv.SetClientScopeHandler(s.ClientScopeHandler) // 设置验证client的scope函数

	// s.Srv.SetClientAuthorizedHandler(s.ClientAuthorizedHandler) 设置某些client用固定的授权方式
	s.Srv.SetInternalErrorHandler(func(err error) *errors.Response {
		log.Println("Internal Error:", err.Error())
		return nil
	})

	s.Srv.SetResponseErrorHandler(func(resp *errors.Response) {
		log.Println("Response Error:", resp.Error.Error())
	})
}

// UserAuthorizeHandler 这个函数是在用户授权后执行的，判断用户是否登录本开放平台
func (s *Server) UserAuthorizeHandler(w http.ResponseWriter, r *http.Request) (string, error) {
	token := r.Header.Get("Authorization")
	if token == "" {
		return "", errors.ErrInvalidRequest
	}

	return tools.ParseToken(token[7:], []byte(s.conf.JWT.SecretKey))
}

// ClientScopeHandler 验证client的scope是否合规
func (s *Server) ClientScopeHandler(tgr *oauth2.TokenGenerateRequest) (bool, error) {
	fmt.Println(tgr.Scope)
	if err := s.enforcer.Enforce(tgr.ClientID, tgr.Scope, "access"); err != nil {
		return false, err
	}
	return true, nil
}

func (s *Server) HandleAuthorizeRequest(w http.ResponseWriter, r *http.Request) error {
	return s.Srv.HandleAuthorizeRequest(w, r)
}

func (s *Server) HandleTokenRequest(w http.ResponseWriter, r *http.Request) error {
	return s.Srv.HandleTokenRequest(w, r)
}

func (s *Server) ValidationBearerToken(r *http.Request) (oauth2.TokenInfo, error) {
	return s.Srv.ValidationBearerToken(r)
}
