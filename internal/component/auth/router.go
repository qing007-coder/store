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

	auth := r.engine.Group("/api/user")
	{
		auth.GET("auth", r.auth.Auth)
		auth.GET("token", r.auth.ExchangeToken)
		auth.POST("register_client", r.auth.RegisterClient)
		auth.GET("validate_token", r.auth.ValidateToken)
		auth.POST("register", r.auth.Register)
		auth.POST("login_by_password", r.auth.LoginByPassword)
		auth.POST("login_by_verification_code", r.auth.LoginByVerificationCode)
		auth.POST("send_email", r.auth.SendEmail)
	}

	return r.engine.Run(":8081")
}
