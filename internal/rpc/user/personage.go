package user

import (
	"context"
	"store/internal/proto/user"
	"store/internal/rpc/base"
	"store/pkg/model"
	"time"
)

type Personage struct {
	*base.Base
}

func NewPersonage(b *base.Base) *Personage {
	return &Personage{b}
}

func (p *Personage) UpdatePersonalInfo(ctx context.Context, req *user.UpdatePersonalInfoReq, resp *user.UpdatePersonalInfoResp) error {
	uid := ctx.Value("user_id").(string)
	p.DB.Where("id = ?", uid).Updates(&model.User{
		Nickname:     req.GetNickname(),
		Introduction: req.GetIntroduction(),
		Gender:       req.GetGender(),
		Sign:         req.GetSign(),
		UpdatedAt:    time.Now(),
	})

	return nil
}

func (p *Personage) ModifyPassword(ctx context.Context, req *user.ModifyPasswordReq, resp *user.ModifyPasswordResp) error {
	uid := ctx.Value("user_id").(string)
	p.RDB.Get()
}
