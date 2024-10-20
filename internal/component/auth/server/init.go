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
	"slices"
	"store/pkg/model"
	"store/pkg/redis"
	"strings"
	"time"
)

type Server struct {
	ctx context.Context
	Srv *server.Server
	rdb *redis.Client
	db  *gorm.DB
}

func NewServer(r *redis.Client, db *gorm.DB) *Server {
	srv := &Server{
		rdb: r,
		db:  db,
	}
	srv.init()
	return srv
}

func (s *Server) init() {
	//manager := manage.NewDefaultManager()  自动初始化一些默认的存储器，比如客户端存储、令牌存储、授权码存储等
	manager := manage.NewManager()
	manager.SetAuthorizeCodeTokenCfg(&manage.Config{
		AccessTokenExp:    time.Hour * 6,
		RefreshTokenExp:   time.Hour * 24 * 3,
		IsGenerateRefresh: true,
	}) // 设置token的配置

	manager.MapTokenStorage(redis_.NewRedisStore(s.rdb.GetConfig())) // 设置存储token的配置
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

	//// 验证 JWT
	//token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
	//	// 验证签名方法
	//	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
	//		return nil, errors.New("invalidToken")
	//	}
	//	return []byte("密钥miyao"), nil
	//})
	//
	//if err != nil || !token.Valid {
	//	return "", errors.New("invalidToken")
	//}
	//
	//// 提取用户信息
	//claims, ok := token.Claims.(jwt.MapClaims)
	//if !ok || !token.Valid {
	//	return "", errors.New("invalidToken")
	//}
	//
	//userID := claims["userID"].(string) // 假设用户 ID 存储在 claims 中
	//return userID, nil                  // 返回用户 ID

	//s.rdb.Get(s.ctx, )

	return "31432342342", nil
}

// ClientScopeHandler 验证client的scope是否合规
func (s *Server) ClientScopeHandler(tgr *oauth2.TokenGenerateRequest) (bool, error) {
	var client model.Client
	if err := s.db.Where("id = ?", tgr.ClientID).First(&client).Error; err != nil {
		return false, err
	}

	scopes := strings.Split(client.Scope, " ")
	ok := slices.Contains(scopes, tgr.Scope)
	if ok {
		return ok, nil
	} else {
		fmt.Println("未授权scope")
		return !ok, errors.New("未授权scope")
	}
}

func (s *Server) HandleAuthorizeRequest(w http.ResponseWriter, r *http.Request) error {
	return s.Srv.HandleAuthorizeRequest(w, r)
}

func (s *Server) HandleTokenRequest(w http.ResponseWriter, r *http.Request) error {
	return s.Srv.HandleTokenRequest(w, r)
}
