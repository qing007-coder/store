package user

import (
	"context"
	"store/internal/proto/user"
	"store/internal/rpc/base"
)

type ReceiverAddress struct {
	*base.Base
}

func NewReceiverAddress(b *base.Base) *ReceiverAddress {
	return &ReceiverAddress{b}
}

func (r *ReceiverAddress) AddReceiverAddress(context.Context, *user.AddReceiverAddressReq, *user.AddReceiverAddressResp) error {

}

func (r *ReceiverAddress) DeleteReceiverAddress(context.Context, *user.DeleteReceiverAddressReq, *user.DeleteReceiverAddressResp) error {

}

func (r *ReceiverAddress) UpdateReceiverAddress(context.Context, *user.UpdateReceiverAddressReq, *user.UpdateReceiverAddressResp) error {

}

func (r *ReceiverAddress) GetReceiverAddress(context.Context, *user.GetReceiverAddressReq, *user.GetReceiverAddressResp) error {

}
