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

type ReceiverAddress struct {
	*base.Base
}

func NewReceiverAddress(b *base.Base) *ReceiverAddress {
	return &ReceiverAddress{b}
}

func (r *ReceiverAddress) AddReceiverAddress(ctx context.Context, req *user.AddReceiverAddressReq, resp *user.AddReceiverAddressResp) error {
	uid := ctx.Value("user_id").(string)

	if err := r.DB.Create(&model.ReceiverAddress{
		DetailedAddress: req.GetDetailedAddress(),
		ReceiverName:    req.GetReceiver(),
		PhoneNumber:     req.GetPhone(),
		Label:           req.GetLabel(),
		UserID:          uid,
	}).Error; err != nil {
		r.Logger.Error(errors.DBCreateError.Error(), resource.USERMODULE)
		return errors.DBCreateError
	}

	resp.Code = rsp.OK
	resp.Message = rsp.CREATESUCCESS

	return nil
}

func (r *ReceiverAddress) DeleteReceiverAddress(ctx context.Context, req *user.DeleteReceiverAddressReq, resp *user.DeleteReceiverAddressResp) error {
	if err := r.DB.Where("id = ?", req.GetId()).Delete(&model.ReceiverAddress{}).Error; err != nil {
		r.Logger.Error(errors.DBDeleteError.Error(), resource.USERMODULE)
		return errors.DBDeleteError
	}

	resp.Code = rsp.OK
	resp.Message = rsp.DELETESUCCESS

	return nil
}

func (r *ReceiverAddress) UpdateReceiverAddress(ctx context.Context, req *user.UpdateReceiverAddressReq, resp *user.UpdateReceiverAddressResp) error {
	if err := r.DB.Where("id = ?", req.GetId()).Updates(&model.ReceiverAddress{
		DetailedAddress: req.GetDetailedAddress(),
		ReceiverName:    req.GetReceiver(),
		PhoneNumber:     req.GetPhone(),
		Label:           req.GetLabel(),
	}).Error; err != nil {
		r.Logger.Error(errors.DBUpdateError.Error(), resource.USERMODULE)
		return errors.DBUpdateError
	}

	resp.Code = rsp.OK
	resp.Message = rsp.UPDATESUCCESS

	return nil
}

func (r *ReceiverAddress) GetReceiverAddress(ctx context.Context, req *user.GetReceiverAddressReq, resp *user.GetReceiverAddressResp) error {
	uid := ctx.Value("user_id").(string)

	var count int64
	var addresses []model.ReceiverAddress
	r.DB.Where("user_id = ?", uid).Limit(int(req.GetSize())).Offset(int((req.GetReq() - 1) * req.GetSize())).Find(&addresses).Count(&count)

	data, err := json.Marshal(&addresses)
	if err != nil {
		r.Logger.Error(errors.JsonMarshalError.Error(), resource.USERMODULE)
		return errors.JsonMarshalError
	}

	resp.Code = rsp.OK
	resp.Data = data
	resp.Message = rsp.SEARCHSUCCESS
	resp.Total = count

	return nil
}
