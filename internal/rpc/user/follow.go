package user

import (
	"context"
	"store/internal/proto/user"
	"store/internal/rpc/base"
	"store/pkg/constant/resource"
	rsp "store/pkg/constant/response"
	"store/pkg/errors"
	"store/pkg/model"
)

type Follow struct {
	*base.Base
}

func NewFollow(b *base.Base) *Follow {
	return &Follow{b}
}

func (r *Footprint) FollowMerchant(ctx context.Context, req *user.FollowMerchantReq, resp *user.FollowMerchantResp) error {
	uid := ctx.Value("user_id").(string)

	if err := r.DB.Create(&model.Follow{
		UserID:     uid,
		MerchantID: req.GetMerchantID(),
	}).Error; err != nil {
		r.Logger.Error(errors.DBCreateError.Error(), resource.USERMODULE)
		return errors.DBCreateError
	}

	resp.Code = rsp.OK
	resp.Message = rsp.CREATESUCCESS

	return nil
}

func (r *Footprint) CancelFollow(ctx context.Context, req *user.CancelFollowReq, resp *user.CancelFollowResp) error {
	uid := ctx.Value("user_id").(string)

	if err := r.DB.Where("user_id = ? AND merchant_id = ?", uid, req.GetMerchantID()).Delete(&model.Follow{}).Error; err != nil {
		r.Logger.Error(errors.DBUpdateError.Error(), resource.USERMODULE)
		return errors.DBUpdateError
	}

	resp.Code = rsp.OK
	resp.Message = rsp.UPDATESUCCESS

	return nil
}

func (r *Footprint) GetFollowList(ctx context.Context, req *user.GetFollowListReq, resp *user.GetFollowListResp) error {
	uid := ctx.Value("user_id").(string)

	var followers []model.Follow
	if err := r.DB.Joins("JOIN  ON ").Where("user_id = ?", uid).Find(&followers).Error; err != nil {
		r.Logger.Error(errors.DBQueryError.Error(), resource.USERMODULE)
		return errors.DBQueryError
	}

	resp.Code = rsp.OK
	resp.Message = rsp.SEARCHSUCCESS

	return nil
}
