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

	var favourites []model.Favourites
	if err := r.DB.Where("user_id = ? AND category = ?", uid, req.GetCategory()).Find(&favourites).Error; err != nil {
		r.Logger.Error(errors.DBQueryError.Error(), resource.USERMODULE)
		return errors.DBQueryError
	}

	resp.Code = rsp.OK
	resp.Message = rsp.SEARCHSUCCESS
	return nil
}
