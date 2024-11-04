package gateway

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go-micro.dev/v4/metadata"
	"store/internal/proto/user"
	rsp "store/pkg/constant/response"
	"store/pkg/errors"
	"store/pkg/model"
	"store/pkg/model/request"
	"store/pkg/tools"
)

type UserApi struct {
	ctx    context.Context
	srv    *Service
	client user.UserService
}

func NewUserApi(srv *Service) *UserApi {
	u := new(UserApi)
	u.init(srv)

	return u
}

func (u *UserApi) init(srv *Service) {
	u.ctx = context.Background()
	u.srv = srv
	u.client = user.NewUserService("user", u.srv.Client())
}

func (u *UserApi) UpdatePersonalInfo(ctx *gin.Context) {
	var req request.UpdatePersonalInfoReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		tools.BadRequest(ctx, errors.ShouldBindJsonError.Error())
		return
	}

	c := metadata.Set(u.ctx, "Authorization", ctx.GetString("Authorization"))
	resp, err := u.client.UpdatePersonalInfo(c, &user.UpdatePersonalInfoReq{
		Nickname:     req.Nickname,
		Introduction: req.Introduction,
		Gender:       req.Gender,
		Sign:         req.Sign,
	})

	if err != nil {
		tools.BadRequest(ctx, err.Error())
		return
	}

	if resp.GetCode() != rsp.OK {
		tools.BadRequest(ctx, resp.Message)
		return
	}

	tools.StatusOK(ctx, nil, resp.GetMessage())
}

func (u *UserApi) ModifyPassword(ctx *gin.Context) {
	var req request.ModifyPasswordReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		tools.BadRequest(ctx, errors.ShouldBindJsonError.Error())
		return
	}

	c := metadata.Set(u.ctx, "Authorization", ctx.GetString("Authorization"))
	resp, err := u.client.ModifyPassword(c, &user.ModifyPasswordReq{
		NewPassword:      req.NewPassword,
		VerificationCode: req.VerificationCode,
	})

	if err != nil {
		tools.BadRequest(ctx, err.Error())
		return
	}

	if resp.GetCode() != rsp.OK {
		tools.BadRequest(ctx, resp.Message)
		return
	}

	tools.StatusOK(ctx, nil, resp.GetMessage())
}

func (u *UserApi) AddReceiverAddress(ctx *gin.Context) {
	var req request.AddReceiverAddressReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		tools.BadRequest(ctx, errors.ShouldBindJsonError.Error())
		return
	}

	c := metadata.Set(u.ctx, "Authorization", ctx.GetString("Authorization"))
	resp, err := u.client.AddReceiverAddress(c, &user.AddReceiverAddressReq{
		DetailedAddress: req.DetailedAddress,
		Receiver:        req.Receiver,
		Phone:           req.Phone,
		Label:           req.Label,
	})

	if err != nil {
		tools.BadRequest(ctx, err.Error())
		return
	}

	if resp.GetCode() != rsp.OK {
		tools.BadRequest(ctx, resp.Message)
		return
	}

	tools.StatusOK(ctx, nil, resp.GetMessage())
}

func (u *UserApi) DeleteReceiverAddress(ctx *gin.Context) {
	var req request.DeleteReceiverAddressReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		tools.BadRequest(ctx, errors.ShouldBindJsonError.Error())
		return
	}

	c := metadata.Set(u.ctx, "Authorization", ctx.GetString("Authorization"))
	resp, err := u.client.DeleteReceiverAddress(c, &user.DeleteReceiverAddressReq{
		Id: req.ID,
	})

	if err != nil {
		tools.BadRequest(ctx, err.Error())
		return
	}

	if resp.GetCode() != rsp.OK {
		tools.BadRequest(ctx, resp.Message)
		return
	}

	tools.StatusOK(ctx, nil, resp.GetMessage())
}

func (u *UserApi) UpdateReceiverAddress(ctx *gin.Context) {
	var req request.UpdateReceiverAddressReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		tools.BadRequest(ctx, errors.ShouldBindJsonError.Error())
		return
	}

	c := metadata.Set(u.ctx, "Authorization", ctx.GetString("Authorization"))
	resp, err := u.client.UpdateReceiverAddress(c, &user.UpdateReceiverAddressReq{
		Id:              req.ID,
		DetailedAddress: req.DetailedAddress,
		Receiver:        req.Receiver,
		Phone:           req.Phone,
		Label:           req.Label,
	})

	if err != nil {
		tools.BadRequest(ctx, err.Error())
		return
	}

	if resp.GetCode() != rsp.OK {
		tools.BadRequest(ctx, resp.Message)
		return
	}

	tools.StatusOK(ctx, nil, resp.GetMessage())
}

func (u *UserApi) GetReceiverAddress(ctx *gin.Context) {
	c := metadata.Set(u.ctx, "Authorization", ctx.GetString("Authorization"))
	resp, err := u.client.GetReceiverAddress(c, &user.GetReceiverAddressReq{})

	if err != nil {
		tools.BadRequest(ctx, err.Error())
		return
	}

	if resp.GetCode() != rsp.OK {
		tools.BadRequest(ctx, resp.Message)
		return
	}

	var list []model.ReceiverAddress
	if err := json.Unmarshal(resp.Data, &list); err != nil {
		tools.BadRequest(ctx, err.Error())
		return
	}

	tools.StatusOK(ctx, gin.H{
		"list": list,
	}, resp.GetMessage())
}

func (u *UserApi) AddFavourites(ctx *gin.Context) {
	var req request.AddFavouritesReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		tools.BadRequest(ctx, errors.ShouldBindJsonError.Error())
		return
	}

	c := metadata.Set(u.ctx, "Authorization", ctx.GetString("Authorization"))
	resp, err := u.client.AddFavourites(c, &user.AddFavouritesReq{
		TargetID: req.TargetID,
		Category: req.Category,
	})

	if err != nil {
		tools.BadRequest(ctx, err.Error())
		return
	}

	if resp.GetCode() != rsp.OK {
		tools.BadRequest(ctx, resp.Message)
		return
	}

	tools.StatusOK(ctx, nil, resp.GetMessage())
}

func (u *UserApi) DeleteFavourites(ctx *gin.Context) {
	var req request.DeleteFavouritesReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		tools.BadRequest(ctx, errors.ShouldBindJsonError.Error())
		return
	}

	c := metadata.Set(u.ctx, "Authorization", ctx.GetString("Authorization"))
	resp, err := u.client.DeleteFavourites(c, &user.DeleteFavouritesReq{
		TargetID: req.TargetID,
	})

	if err != nil {
		tools.BadRequest(ctx, err.Error())
		return
	}

	if resp.GetCode() != rsp.OK {
		tools.BadRequest(ctx, resp.Message)
		return
	}

	tools.StatusOK(ctx, nil, resp.GetMessage())
}

func (u *UserApi) GetFavouritesList(ctx *gin.Context) {
	category := ctx.Query("category")

	c := metadata.Set(u.ctx, "Authorization", ctx.GetString("Authorization"))
	resp, err := u.client.GetFavouritesList(c, &user.GetFavouritesListReq{
		Category: category,
	})

	if err != nil {
		tools.BadRequest(ctx, err.Error())
		return
	}

	if resp.GetCode() != rsp.OK {
		tools.BadRequest(ctx, resp.Message)
		return
	}

	var list []model.Favourites
	if err := json.Unmarshal(resp.Data, &list); err != nil {
		tools.BadRequest(ctx, err.Error())
		return
	}

	tools.StatusOK(ctx, gin.H{
		"list": list,
	}, resp.GetMessage())
}

func (u *UserApi) AddFootprint(ctx *gin.Context) {
	var req request.AddFootprintReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		tools.BadRequest(ctx, errors.ShouldBindJsonError.Error())
		return
	}

	c := metadata.Set(u.ctx, "Authorization", ctx.GetString("Authorization"))
	resp, err := u.client.AddFootprint(c, &user.AddFootprintReq{
		TargetID: req.TargetID,
		Category: req.Category,
	})

	if err != nil {
		tools.BadRequest(ctx, err.Error())
		return
	}

	if resp.GetCode() != rsp.OK {
		tools.BadRequest(ctx, resp.Message)
		return
	}

	tools.StatusOK(ctx, nil, resp.GetMessage())
}

func (u *UserApi) DeleteFootprint(ctx *gin.Context) {
	var req request.DeleteFootprintReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		tools.BadRequest(ctx, errors.ShouldBindJsonError.Error())
		return
	}

	c := metadata.Set(u.ctx, "Authorization", ctx.GetString("Authorization"))
	resp, err := u.client.DeleteFootprint(c, &user.DeleteFootprintReq{
		TargetID: req.TargetID,
	})

	if err != nil {
		tools.BadRequest(ctx, err.Error())
		return
	}

	if resp.GetCode() != rsp.OK {
		tools.BadRequest(ctx, resp.Message)
		return
	}

	tools.StatusOK(ctx, nil, resp.GetMessage())
}

func (u *UserApi) GetFootprintList(ctx *gin.Context) {
	category := ctx.Query("category")

	c := metadata.Set(u.ctx, "Authorization", ctx.GetString("Authorization"))
	resp, err := u.client.GetFootprintList(c, &user.GetFootprintListReq{
		Category: category,
	})

	if err != nil {
		tools.BadRequest(ctx, err.Error())
		return
	}

	if resp.GetCode() != rsp.OK {
		tools.BadRequest(ctx, resp.Message)
		return
	}

	var list []model.Footprint
	if err := json.Unmarshal(resp.Data, &list); err != nil {
		tools.BadRequest(ctx, err.Error())
		return
	}

	tools.StatusOK(ctx, gin.H{
		"list": list,
	}, resp.GetMessage())
}

func (u *UserApi) FollowMerchant(ctx *gin.Context) {
	var req request.FollowMerchantReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		tools.BadRequest(ctx, errors.ShouldBindJsonError.Error())
		return
	}

	c := metadata.Set(u.ctx, "Authorization", ctx.GetString("Authorization"))
	resp, err := u.client.FollowMerchant(c, &user.FollowMerchantReq{
		MerchantID: req.MerchantID,
	})

	if err != nil {
		tools.BadRequest(ctx, err.Error())
		return
	}

	if resp.GetCode() != rsp.OK {
		tools.BadRequest(ctx, resp.Message)
		return
	}

	tools.StatusOK(ctx, nil, resp.GetMessage())
}

func (u *UserApi) CancelFollow(ctx *gin.Context) {
	var req request.CancelFollowReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		tools.BadRequest(ctx, errors.ShouldBindJsonError.Error())
		return
	}

	c := metadata.Set(u.ctx, "Authorization", ctx.GetString("Authorization"))
	resp, err := u.client.CancelFollow(c, &user.CancelFollowReq{
		MerchantID: req.MerchantID,
	})

	if err != nil {
		tools.BadRequest(ctx, err.Error())
		return
	}

	if resp.GetCode() != rsp.OK {
		tools.BadRequest(ctx, resp.Message)
		return
	}

	tools.StatusOK(ctx, nil, resp.GetMessage())
}

func (u *UserApi) GetFollowList(ctx *gin.Context) {
	c := metadata.Set(u.ctx, "Authorization", ctx.GetString("Authorization"))
	resp, err := u.client.GetFollowList(c, &user.GetFollowListReq{})

	if err != nil {
		tools.BadRequest(ctx, err.Error())
		return
	}

	if resp.GetCode() != rsp.OK {
		tools.BadRequest(ctx, resp.Message)
		return
	}

	var users []model.User
	if err := json.Unmarshal(resp.Data, &users); err != nil {
		tools.BadRequest(ctx, err.Error())
		return
	}

	tools.StatusOK(ctx, gin.H{
		"list": users,
	}, resp.GetMessage())
}
