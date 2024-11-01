package user

import (
	"context"
	"store/internal/proto/user"
	"store/internal/rpc/base"
)

type Footprint struct {
	*base.Base
}

func NewFootprint(b *base.Base) *Footprint {
	return &Footprint{b}
}

func (r *Footprint) AddFootprint(context.Context, *user.AddFootprintReq, *user.AddFootprintResp) error {

}

func (r *Footprint) DeleteFootprint(context.Context, *user.DeleteFootprintReq, *user.DeleteFootprintResp) error {

}

func (r *Footprint) GetFootprintList(context.Context, *user.GetFootprintListReq, *user.GetFootprintListResp) error {

}
