package user

import (
	"context"
	"store/internal/proto/user"
	"store/internal/rpc/base"
)

type Follow struct {
	*base.Base
}

func NewFollow(b *base.Base) *Follow {
	return &Follow{b}
}

func (r *Footprint) FollowMerchant(context.Context, *user.FollowMerchantReq, *user.FollowMerchantResp) error {

}

func (r *Footprint) CancelFollow(context.Context, *user.CancelFollowReq, *user.CancelFollowResp) error {

}

func (r *Footprint) GetFollowList(context.Context, *user.GetFollowListReq, *user.GetFollowListResp) error {

}
