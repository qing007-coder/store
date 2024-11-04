package user

import (
	"context"
	"encoding/json"
	"store/internal/proto/user"
	"store/internal/rpc/base"
	"store/internal/rpc/tools"
	"store/pkg/constant"
	"store/pkg/constant/resource"
	rsp "store/pkg/constant/response"
	"store/pkg/errors"
	"store/pkg/model"
)

type Footprint struct {
	*base.Base
}

func NewFootprint(b *base.Base) *Footprint {
	return &Footprint{b}
}

func (r *Footprint) AddFootprint(ctx context.Context, req *user.AddFootprintReq, resp *user.AddFootprintResp) error {
	uid := ctx.Value("user_id").(string)

	if err := r.DB.Create(&model.Footprint{
		UserID:   uid,
		TargetID: req.GetTargetID(),
		Category: req.GetCategory(),
	}).Error; err != nil {
		r.Logger.Error(errors.DBCreateError.Error(), resource.USERMODULE)
		return errors.DBCreateError
	}

	resp.Code = rsp.OK
	resp.Message = rsp.CREATESUCCESS

	return nil
}

func (r *Footprint) DeleteFootprint(ctx context.Context, req *user.DeleteFootprintReq, resp *user.DeleteFootprintResp) error {
	uid := ctx.Value("user_id").(string)

	if err := r.DB.Where("user_id = ? AND target_id", uid, req.GetTargetID()).Delete(&model.Footprint{}).Error; err != nil {
		r.Logger.Error(errors.DBDeleteError.Error(), resource.USERMODULE)
		return errors.EsDeleteError
	}

	resp.Code = rsp.OK
	resp.Message = rsp.DELETESUCCESS
	return nil
}

func (r *Footprint) GetFootprintList(ctx context.Context, req *user.GetFootprintListReq, resp *user.GetFootprintListResp) error {
	uid := ctx.Value("user_id").(string)

	if req.GetCategory() == constant.MERCHANT {
		var u []model.User
		if err := r.DB.Joins("JOIN user ON user.id = footprint.target_id").Where("footprint.user_id = ? AND footprint.category = ?", uid, req.GetCategory()).Find(&u).Error; err != nil {
			r.Logger.Error(errors.DBQueryError.Error(), resource.USERMODULE)
			return errors.DBQueryError
		}

		data, err := json.Marshal(&u)
		if err != nil {
			r.Logger.Error(errors.JsonMarshalError.Error(), resource.USERMODULE)
			return errors.JsonMarshalError
		}

		resp.Data = data
	}

	if req.GetCategory() == constant.MERCHANDISE {
		var footprintList []model.Footprint
		if err := r.DB.Where("user_id = ? AND category = ?", uid, req.GetCategory()).Find(&footprintList).Error; err != nil {
			r.Logger.Error(errors.DBQueryError.Error(), resource.USERMODULE)
			return errors.DBQueryError
		}

		var ms []model.Merchandise

		for _, f := range footprintList {
			m, err := tools.GetMerchandise(f.TargetID, r.ES[constant.MERCHANDISE])
			if err != nil {
				r.Logger.Error(err.Error(), resource.USERMODULE)
				return err
			}

			ms = append(ms, m)
		}

		data, err := json.Marshal(&ms)
		if err != nil {
			r.Logger.Error(errors.JsonMarshalError.Error(), resource.USERMODULE)
			return errors.JsonMarshalError
		}

		resp.Data = data
	}

	resp.Code = rsp.OK
	resp.Message = rsp.SEARCHSUCCESS

	return nil
}
