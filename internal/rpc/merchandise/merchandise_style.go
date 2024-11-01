package merchandise

import (
	"context"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"store/internal/proto/merchandise"
	"store/internal/rpc/base"
	"store/pkg/constant"
	rsp "store/pkg/constant/response"
	"store/pkg/errors"
	"store/pkg/model"
	"store/pkg/tools"
	"time"
)

type MerchandiseStyle struct {
	*base.Base
}

func NewMerchandiseStyle(b *base.Base) *MerchandiseStyle {
	return &MerchandiseStyle{b}
}

func (m *MerchandiseStyle) AddMerchandiseStyle(ctx context.Context, req *merchandise.AddMerchandiseStyleReq, resp *merchandise.AddMerchandiseStyleResp) error {
	uid := ctx.Value("user_id").(string)
	id := tools.CreateID()
	if err := m.ES[constant.MERCHANDISESTYLE].CreateDocument(&model.MerchandiseStyle{
		ID:            id,
		MerchandiseID: req.GetMerchandiseID(),
		Name:          req.GetName(),
		Info:          req.GetInfo(),
		Picture:       req.GetPicture(),
		Price:         req.GetPrice(),
		Stock:         req.GetStock(),
		Status:        req.GetStatus(),
		CreatedAt:     time.Now().Unix(),
	}, id); err != nil {
		m.Logger.Error(errors.EsCreateError.Error(), constant.MERCHANDISESTYLE)
		return err
	}

	m.DB.Create(&model.MerchantRecord{
		ID:     tools.CreateID(),
		Time:   time.Now(),
		UserID: uid,
		Action: constant.ADD,
		Source: constant.MERCHANDISESTYLE,
	})

	resp.Code = rsp.OK
	resp.Message = rsp.CREATESUCCESS

	return nil
}

func (m *MerchandiseStyle) RemoveMerchandiseStyle(ctx context.Context, req *merchandise.RemoveMerchandiseStyleReq, resp *merchandise.RemoveMerchandiseStyleResp) error {
	uid := ctx.Value("user_id").(string)
	if err := m.ES[constant.MERCHANDISESTYLE].DeleteDocument(req.GetId()); err != nil {
		m.Logger.Error(errors.EsDeleteError.Error(), constant.MERCHANDISESTYLE)
		return err
	}

	m.DB.Create(&model.MerchantRecord{
		ID:     tools.CreateID(),
		Time:   time.Now(),
		UserID: uid,
		Action: constant.DELETE,
		Source: constant.MERCHANDISESTYLE,
	})

	resp.Code = rsp.OK
	resp.Message = rsp.CREATESUCCESS

	return nil
}

func (m *MerchandiseStyle) UpdateMerchandiseStyle(ctx context.Context, req *merchandise.UpdateMerchandiseStyleReq, resp *merchandise.UpdateMerchandiseStyleResp) error {
	uid := ctx.Value("user_id").(string)

	var queries map[string]interface{}
	if req.GetName() != "" {
		queries["name"] = req.GetName()
	}

	if req.GetInfo() != "" {
		queries["info"] = req.GetInfo()
	}

	if req.GetPicture() != "" {
		queries["picture"] = req.GetPicture()
	}

	if req.GetPrice() != -1 {
		queries["price"] = req.GetPrice()
	}

	if req.GetStatus() != "" {
		queries["status"] = req.GetStatus()
	}

	if req.GetStock() != -1 {
		if req.GetStock() < 0 {
			m.Logger.Error(errors.UndefinedValue("stock").Error(), constant.MERCHANDISESTYLE)
		}
		queries["stock"] = req.GetStock()
	}

	queries["updated_at"] = time.Now().Unix()

	if err := m.ES[constant.MERCHANDISESTYLE].Update(req.GetId(), queries); err != nil {
		m.Logger.Error(errors.EsUpdateError.Error(), constant.MERCHANDISESTYLE)
		return err
	}

	m.DB.Create(&model.MerchantRecord{
		ID:     tools.CreateID(),
		Time:   time.Now(),
		UserID: uid,
		Action: constant.UPDATE,
		Source: constant.MERCHANDISESTYLE,
	})

	resp.Code = rsp.OK
	resp.Message = rsp.UPDATESUCCESS

	return nil
}

func (m *MerchandiseStyle) GetMerchandiseStyleList(ctx context.Context, req *merchandise.GetMerchandiseStyleListReq, resp *merchandise.GetMerchandiseStyleListResp) error {
	must := []types.Query{
		{
			MatchPhrase: map[string]types.MatchPhraseQuery{
				"merchandise_id": {Query: req.GetMerchandiseID()},
			},
		},
	}
	response, err := m.ES[constant.MERCHANDISESTYLE].Search(must, nil, nil, int(req.GetReq()*req.GetSize()), int(req.GetSize()))
	if err != nil {
		m.Logger.Error(errors.EsSearchError.Error(), constant.MERCHANDISESTYLE)
		return err
	}

	var ms []model.MerchandiseStyle
	for _, v := range response.Hits.Hits {
		var style model.MerchandiseStyle
		if err := json.Unmarshal(v.Source_, &style); err != nil {
			m.Logger.Error(errors.JsonUnmarshalError.Error(), constant.MERCHANDISESTYLE)
			return err
		}

		ms = append(ms, style)
	}

	data, err := json.Marshal(&ms)
	if err != nil {
		m.Logger.Error(errors.JsonMarshalError.Error(), constant.MERCHANDISESTYLE)
		return err
	}

	resp.Code = rsp.OK
	resp.Data = data
	resp.Message = rsp.SEARCHSUCCESS

	return nil
}

func (m *MerchandiseStyle) GetMerchandiseStyleDetails(ctx context.Context, req *merchandise.GetMerchandiseStyleDetailsReq, resp *merchandise.GetMerchandiseStyleDetailsResp) error {
	data, err := m.ES[constant.MERCHANDISESTYLE].GetDocumentByID(req.GetId())
	if err != nil {
		m.Logger.Error(errors.EsSearchError.Error(), constant.MERCHANDISESTYLE)
		return err
	}

	resp.Code = rsp.OK
	resp.Data = data
	resp.Message = rsp.SEARCHSUCCESS

	return nil
}
