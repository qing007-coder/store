package auth

import (
	"github.com/gin-gonic/gin"
)

type Router struct {
	auth   *AuthApi
	engine *gin.Engine
}

func NewRouter(auth *AuthApi) *Router {
	return &Router{
		auth:   auth,
		engine: gin.Default(),
	}
}

func (r *Router) Run() error {
	r.engine.GET("auth", r.auth.Auth)
	r.engine.GET("token", r.auth.ExchangeToken)
	r.engine.POST("register_client", r.auth.RegisterClient)
	r.engine.GET("validate_token", r.auth.ValidateToken)
	r.engine.POST("register", r.auth.Register)
	r.engine.POST("login", r.auth.Login)
	r.engine.POST("send_email", r.auth.SendEmail)

	return r.engine.Run(":8081")
}
