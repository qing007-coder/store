package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"store/pkg/model"
	"store/pkg/rules"
	"store/pkg/tools"
)

type AuthApi struct {
	srv *Server
	db  *gorm.DB
	e   *rules.Enforcer
}

func NewAuthApi(srv *Server, db *gorm.DB, e *rules.Enforcer) *AuthApi {
	return &AuthApi{
		srv: srv,
		db:  db,
		e:   e,
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

func (a *AuthApi) ValidateToken(ctx *gin.Context) {
	tokenInfo, err := a.srv.ValidationBearerToken(ctx.Request)
	if err != nil {
		ctx.JSON(http.StatusOK, model.ValidateTokenResp{
			Code:    400,
			Message: err.Error(),
		})

		return
	}

	uid := tokenInfo.GetUserID()
	ctx.JSON(http.StatusOK, model.ValidateTokenResp{
		Code:    200,
		Message: uid,
	})
}

func (a *AuthApi) RegisterClient(ctx *gin.Context) {
	var req model.RegisterClientReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		tools.BadRequest(ctx, err.Error())
		return
	}

	id := tools.CreateID()
	secret, err := tools.GenerateSecret(32)
	if err != nil {
		tools.BadRequest(ctx, err.Error())
		return
	}

	var count int64
	a.db.Where("user_id = ?", req.UserID).Count(&count)
	if count > 0 {
		tools.BadRequest(ctx, "你已有该账户")
		return
	}
	if err := a.db.Transaction(func(tx *gorm.DB) error {
		a.db.Create(&model.Client{
			ID:     id,
			Secret: secret,
			Domain: "",
			UserID: req.UserID,
		})

		if err := a.e.AddGroup(id, req.Scope...); err != nil {
			return err
		}

		return nil
	}); err != nil {
		tools.BadRequest(ctx, err.Error())
		return
	}

	tools.StatusOK(ctx, gin.H{
		"client_id":     id,
		"client_secret": secret,
	}, "注册成功")
}
