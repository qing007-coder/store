package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"store/pkg/tools"
)

type AuthApi struct {
	srv *Server
}

func NewAuthApi(srv *Server) *AuthApi {
	return &AuthApi{
		srv: srv,
	}
}

func (a *AuthApi) Auth(ctx *gin.Context) {
	if err := a.srv.HandleAuthorizeRequest(ctx.Writer, ctx.Request); err != nil {
		fmt.Println("err:", err)
		tools.BadRequest(ctx, err.Error())
		return
	}
	fmt.Println("success")
}

func (a *AuthApi) ExchangeToken(ctx *gin.Context) {
	if err := a.srv.HandleTokenRequest(ctx.Writer, ctx.Request); err != nil {
		fmt.Println("err:", err)
		tools.BadRequest(ctx, err.Error())
		return
	}
	fmt.Println("success")
}
