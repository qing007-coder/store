package user

import (
	"context"
	"store/internal/proto/user"
	"store/internal/rpc/base"
)

type Favourites struct {
	*base.Base
}

func NewFavourites(b *base.Base) *Favourites {
	return &Favourites{b}
}

func (r *Favourites) AddFavourites(context.Context, *user.AddFavouritesReq, *user.AddFavouritesResp) error {

}

func (r *Favourites) DeleteFavourites(context.Context, *user.DeleteFavouritesReq, *user.DeleteFavouritesResp) error {

}

func (r *Favourites) GetFavouritesList(context.Context, *user.GetFavouritesListReq, *user.GetFavouritesListResp) error {

}
