package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go-micro.dev/v4/metadata"
	"go-micro.dev/v4/server"
	"net/http"
	"store/pkg/model"
)

type AuthMiddleware struct {
}

func NewAuthMiddleware() *AuthMiddleware {
	c := new(AuthMiddleware)
	return c
}

func (a *AuthMiddleware) Auth(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, resp interface{}) error {
		token, ok := metadata.Get(ctx, "Authorization")
		if !ok {
			fmt.Println("no-ok")
			return errors.New("need token")
		}

		fmt.Println("token:", token)
		authURL := "http://127.0.0.1:8081/validate_token"
		client := new(http.Client)
		request, err := http.NewRequest("GET", authURL, nil)
		if err != nil {
			fmt.Println("err:", err)
			return err
		}

		request.Header.Set("Authorization", token)
		request.Header.Set("Content-Type", "application/json")

		response, err := client.Do(request)
		if err != nil {
			fmt.Println("err:", err)
			return err
		}

		var res model.ValidateTokenResp
		if err := json.NewDecoder(response.Body).Decode(&res); err != nil {
			fmt.Println("err:", err)
			return err
		}

		if res.Code != 200 {
			fmt.Println(res.Message)
			return errors.New(res.Message)
		}

		ctx = context.WithValue(ctx, "user_id", res.Message)
		return fn(ctx, req, resp)
	}
}
