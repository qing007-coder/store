package gateway

import (
	"context"
	"github.com/gin-gonic/gin"
	"store/internal/proto/merchandise"
)

type MerchantApi struct {
	ctx    context.Context
	srv    *Service
	client merchandise.MerchandiseService
}

func NewMerchantApi(srv *Service) *MerchantApi {
	m := new(MerchantApi)
	m.init(srv)

	return m
}

func (m *MerchantApi) init(srv *Service) {
	m.ctx = context.Background()
	m.srv = srv
	m.client = merchandise.NewMerchandiseService("merchant", m.srv.Client())
}

func (m *MerchantApi) PutAwayMerchandise(ctx *gin.Context) {

}

func (m *MerchantApi) RemoveMerchandise(ctx *gin.Context) {

}
