package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

type AuthMiddleware struct {
	conf oauth2.Config
}

func NewAuthMiddleware() *AuthMiddleware {
	c := new(AuthMiddleware)
	c.init()
	return c
}

func (a *AuthMiddleware) init() {

	serverURL := "http://localhost:8081"

	a.conf = oauth2.Config{
		ClientID:     "1828425459575033856",
		ClientSecret: "kGRWnrQG1ORxUf-Pc9BYBljylLsgzVaky5z4VVqJlwQ",
		RedirectURL:  "http://localhost:8080/callback",
		Scopes:       []string{"table"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  serverURL + "/auth",
			TokenURL: serverURL + "/token",
		},
	}
}

func (a *AuthMiddleware) Callback(ctx *gin.Context) {
	code := ctx.Query("code") // 授权码
	//returnState := ctx.Query("state")
	//
	//store, err := session.Start(context.Background(), ctx.Writer, ctx.Request)
	//if err != nil {
	//	fmt.Println("err:", err)
	//	ctx.JSON(200, err.Error())
	//	return
	//}
	//
	//savedState, ok := store.Get("state")
	//if !ok {
	//	fmt.Println("session没有state")
	//	ctx.JSON(200, "session没有state")
	//	return
	//}
	//
	//state, ok := savedState.(string)
	//if !ok {
	//	fmt.Println("断言失败")
	//	ctx.JSON(200, "断言失败")
	//	return
	//}
	//
	//if state != returnState {
	//	fmt.Println("csrf")
	//	ctx.JSON(200, "csrf")
	//	return
	//}

	token, err := c.conf.Exchange(context.Background(), code)
	if err != nil {
		fmt.Println("err:", err)
		ctx.JSON(200, err.Error())
		return
	}

	ctx.JSON(200, gin.H{
		"token": token,
	})
}

func (c *Client) RefreshToken(ctx *gin.Context) {
	access := ctx.Query("access_token")
	refresh := ctx.Query("refresh_token")

	token, err := c.conf.TokenSource(context.Background(), &oauth2.Token{
		AccessToken:  access,
		RefreshToken: refresh,
	}).Token()

	if err != nil {
		fmt.Println("err:", err)
		return
	}

	ctx.JSON(200, gin.H{
		"new_token": token,
	})
}
