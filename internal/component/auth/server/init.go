package server

import (
	"fmt"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/generates"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/server"
	redis_ "github.com/go-oauth2/redis/v4"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"net/http"
	"store/pkg/model"
	"store/pkg/redis"
	"time"
)

type Server struct {
	Srv *server.Server
	rdb *redis.Client
}

func NewServer(r *redis.Client) *Server {
	srv := &Server{
		rdb: r,
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

	cliStore := NewStore(s.rdb)
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
	//tokenString := r.Header.Get("Authorization")
	//if tokenString == "" {
	//	return "", errors.ErrInvalidRequest
	//}
	//
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

	return "", nil
}

// ClientScopeHandler 验证client的scope是否合规
func (s *Server) ClientScopeHandler(tgr *oauth2.TokenGenerateRequest) (bool, error) {
	data, err := s.rdb.Get(context.Background(), tgr.ClientID)
	if err != nil {
		fmt.Println("err:", err)
		return false, err
	}

	var client model.Client
	if err := json.Unmarshal([]byte(data), &client); err != nil {
		fmt.Println("err:", err)
		return false, err
	}

	ok := slices.Contains(client.Scope, tgr.Scope)
	if ok {
		return ok, nil
	} else {
		fmt.Println("未授权scope")
		return !ok, errors.New("未授权scope")
	}
}
