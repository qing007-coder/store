package user

import (
	"context"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"store/internal/proto/user"
	"store/internal/rpc/base"
	"store/pkg/constant/resource"
	rsp "store/pkg/constant/response"
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
	if err := p.DB.Where("id = ?", uid).Updates(&model.User{
		Nickname:     req.GetNickname(),
		Introduction: req.GetIntroduction(),
		Gender:       req.GetGender(),
		Sign:         req.GetSign(),
		UpdatedAt:    time.Now(),
	}).Error; err != nil {
		p.Logger.Error(errors.DBCreateError.Error(), resource.USERMODULE)
		return errors.DBCreateError
	}

	resp.Code = rsp.OK
	resp.Message = rsp.UPDATESUCCESS

	return nil
}

func (p *Personage) ModifyPassword(ctx context.Context, req *user.ModifyPasswordReq, resp *user.ModifyPasswordResp) error {
	uid := ctx.Value("user_id").(string)
	var u model.User
	if err := p.DB.Where("id = ?", uid).First(&u).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			p.Logger.Error(errors.RecordNotFound.Error(), resource.USERMODULE)
			return errors.RecordNotFound
		} else {
			p.Logger.Error(errors.OtherError.Error(), resource.USERMODULE)
			return errors.OtherError
		}
	}

	code, err := p.RDB.Get(ctx, u.Email)
	if err != nil {
		resp.Code = rsp.ERROR
		resp.Message = "verification code is expiry "
		return nil
	}

	if code != req.GetVerificationCode() {
		resp.Code = rsp.ERROR
		resp.Message = "wrong verification code"
		return nil
	}

	password, err := bcrypt.GenerateFromPassword([]byte(req.GetNewPassword()), bcrypt.DefaultCost)
	if err != nil {
		p.Logger.Error(err.Error(), resource.USERMODULE)
		return err
	}

	p.DB.Where("id = ?", uid).Updates(&model.User{
		Password: string(password),
	})

	resp.Code = rsp.OK
	resp.Message = rsp.UPDATESUCCESS

	return nil
}
