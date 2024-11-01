package user

import (
	"context"
	"store/internal/proto/user"
	"store/internal/rpc/base"
)

type Personage struct {
	*base.Base
}

func NewPersonage(b *base.Base) *Personage {
	return &Personage{b}
}

func (p *Personage) UpdatePersonalInfo(context.Context, *user.UpdatePersonalInfoReq, *user.UpdatePersonalInfoResp) error {

}

func (p *Personage) ModifyPassword(context.Context, *user.ModifyPasswordReq, *user.ModifyPasswordResp) error {

}
