package user

import (
	"context"
	"encoding/json"
	"store/internal/proto/user"
	"store/internal/rpc/base"
	"store/pkg/constant/resource"
	rsp "store/pkg/constant/response"
	"store/pkg/errors"
	"store/pkg/model"
)

type Favourites struct {
	*base.Base
}

func NewFavourites(b *base.Base) *Favourites {
	return &Favourites{b}
}

func (r *Favourites) AddFavourites(ctx context.Context, req *user.AddFavouritesReq, resp *user.AddFavouritesResp) error {
	uid := ctx.Value("user_id").(string)

	if err := r.DB.Create(&model.Favourites{
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

func (r *Favourites) DeleteFavourites(ctx context.Context, req *user.DeleteFavouritesReq, resp *user.DeleteFavouritesResp) error {
	uid := ctx.Value("user_id").(string)

	if err := r.DB.Where("user_id = ? AND target_id = ?", uid, req.GetTargetID()).Delete(&model.Favourites{}).Error; err != nil {
		r.Logger.Error(errors.DBDeleteError.Error(), resource.USERMODULE)
		return errors.DBDeleteError
	}

	resp.Code = rsp.OK
	resp.Message = rsp.DELETESUCCESS

	return nil
}

func (r *Favourites) GetFavouritesList(ctx context.Context, req *user.GetFavouritesListReq, resp *user.GetFavouritesListResp) error {
	uid := ctx.Value("user_id").(string)

	var count int64
	var favourites []model.Favourites

	if err := r.DB.Where("user_id = ? AND category = ?", uid, req.GetCategory()).Limit(int(req.GetSize())).Offset(int((req.GetReq() - 1) * req.GetSize())).Find(&favourites).Count(&count).Error; err != nil {
		r.Logger.Error(errors.DBQueryError.Error(), resource.USERMODULE)
		return errors.DBQueryError
	}

	//if req.GetCategory() == constant.MERCHANDISE {
	//	var merchandises []model.Merchandise
	//	for _, f := range favourites {
	//		m, err := tools.GetMerchandise(f.TargetID, r.ES[constant.MERCHANDISE])
	//		if err != nil {
	//			r.Logger.Error(err.Error(), resource.USERMODULE)
	//			return err
	//		}
	//		merchandises = append(merchandises, m)
	//	}
	//
	//	data, err := json.Marshal(&merchandises)
	//	if err != nil {
	//		r.Logger.Error(errors.JsonMarshalError.Error(), resource.USERMODULE)
	//		return errors.JsonMarshalError
	//	}
	//
	//	resp.Data = data
	//}

	data, err := json.Marshal(&favourites)
	if err != nil {
		r.Logger.Error(errors.JsonMarshalError.Error(), resource.USERMODULE)
		return errors.JsonMarshalError
	}

	resp.Code = rsp.OK
	resp.Message = rsp.SEARCHSUCCESS
	resp.Data = data
	resp.Total = count
	return nil
}
