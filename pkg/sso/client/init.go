package client

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-session/session"
	"golang.org/x/oauth2"
	"net/http"
	"store/pkg/tools"
)

type Client struct {
	conf oauth2.Config
}

func NewClient() *Client {
	client := &Client{}

	client.init()
	return client
}

func (c *Client) init() {

	serverURL := "http://localhost:8080"

	c.conf = oauth2.Config{
		ClientID:     "1828425459575033856",
		ClientSecret: "kGRWnrQG1ORxUf-Pc9BYBljylLsgzVaky5z4VVqJlwQ",
		RedirectURL:  "http://localhost:8081/callback",
		Scopes:       []string{"table"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  serverURL + "/auth",
			TokenURL: serverURL + "/token",
		},
	}
}

func (c *Client) Login(ctx *gin.Context) {
	state := tools.CreateID()
	store, err := session.Start(context.Background(), ctx.Writer, ctx.Request)
	if err != nil {
		fmt.Println("err:", err)
		ctx.JSON(200, err.Error())
		return
	}

	store.Set("state", state)
	if err := store.Save(); err != nil {
		fmt.Println("err:", err)
		ctx.JSON(200, err.Error())
		return
	}

	url := c.conf.AuthCodeURL(state)
	ctx.Redirect(http.StatusFound, url)
}

func (c *Client) Callback(ctx *gin.Context) {
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
