package user

import (
	"context"
	"gorm.io/gorm"
	"store/internal/proto/user"
	"store/internal/rpc/base"
	"store/pkg/errors"
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
	var u model.User
	if err := p.DB.Where("id = ?", uid).First(&u).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			p.Logger.Error(errors.RecordNotFound.Error(), constant.)
		}
	}
	p.RDB.Get()
}
