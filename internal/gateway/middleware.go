package gateway

import "github.com/gin-gonic/gin"

type Middleware struct {
}

func NewMiddleware() *Middleware {
	return &Middleware{}
}

func (m *Middleware) GetAuthorizationHeader() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")

		ctx.Set("Authorization", token)
	}
}
