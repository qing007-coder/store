package auth

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"store/pkg/config"
	"store/pkg/errors"
	"store/pkg/model"
	"store/pkg/redis"
	"store/pkg/rules"
	"store/pkg/sso/server"
	"store/pkg/tools"
)

type AuthApi struct {
	srv  *server.Server
	db   *gorm.DB
	rdb  *redis.Client
	e    *rules.Enforcer
	conf *config.GlobalConfig
}

func NewAuthApi(srv *server.Server, db *gorm.DB, e *rules.Enforcer, conf *config.GlobalConfig) *AuthApi {
	return &AuthApi{
		srv:  srv,
		db:   db,
		e:    e,
		conf: conf,
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

func (a *AuthApi) Login(ctx *gin.Context) {
	var req model.LoginReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		tools.BadRequest(ctx, err.Error())
		return
	}

	var user model.User
	if err := a.db.Where("account = ?", req.Account).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			tools.BadRequest(ctx, errors.RecordNotFound.Error())
		} else {
			tools.BadRequest(ctx, errors.OtherError.Error())
			return
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			tools.BadRequest(ctx, errors.New("密码错误").Error())
			return
		} else {
			tools.BadRequest(ctx, errors.OtherError.Error())
			return
		}
	}

	accessToken, err := tools.CreateToken(user.ID, a.conf.JWT.AccessExpiry, []byte(a.conf.JWT.SecretKey))
	if err != nil {
		tools.BadRequest(ctx, err.Error())
		return
	}

	refreshToken, err := tools.CreateToken(user.ID, a.conf.JWT.RefreshExpiry, []byte(a.conf.JWT.SecretKey))
	if err != nil {
		tools.BadRequest(ctx, err.Error())
		return
	}

	tools.StatusOK(ctx, model.Data{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Expiry:       a.conf.JWT.AccessExpiry,
		TokenType:    "Bearer",
	}, "登录成功")
}

func (a *AuthApi) Register(ctx *gin.Context) {

}

func (a *AuthApi) SendEmail(ctx *gin.Context) {

}
