package merchandise

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/minio/minio-go/v7"
	"io"
	"store/internal/proto/merchandise"
	"store/internal/rpc/base"
	"store/pkg/constant/resource"
	rsp "store/pkg/constant/response"
	"store/pkg/constant/rules"
	"store/pkg/constant/store"
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

func (m *MerchandiseStyle) AddMerchandiseStyle(ctx context.Context, stream merchandise.MerchandiseService_AddMerchandiseStyleStream) error {
	uid := ctx.Value("user_id").(string)

	var req *merchandise.AddMerchandiseStyleReq
	var picture bytes.Buffer

	for {
		data, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		req = &merchandise.AddMerchandiseStyleReq{
			MerchandiseID: data.GetMerchandiseID(),
			Name:          data.GetName(),
			Info:          data.GetInfo(),
			Price:         data.GetPrice(),
			Stock:         data.GetStock(),
			Status:        data.GetStatus(),
		}

		chunk := data.GetChunk()
		_, err = picture.Write(chunk.GetData())
		if err != nil {
			return err
		}
	}

	id := tools.CreateID()
	path := fmt.Sprintf("%s", id)
	_, err := m.MC.PutObject(m.Ctx, store.MERCHANDISESTYLE, path, &picture, int64(picture.Len()), minio.PutObjectOptions{})
	if err != nil {
		return err
	}

	if err := m.ES[store.MERCHANDISESTYLE].CreateDocument(&model.MerchandiseStyle{
		ID:            id,
		MerchandiseID: req.GetMerchandiseID(),
		Name:          req.GetName(),
		Info:          req.GetInfo(),
		Picture:       path,
		Price:         req.GetPrice(),
		Stock:         req.GetStock(),
		Status:        req.GetStatus(),
		CreatedAt:     time.Now().Unix(),
	}, id); err != nil {
		m.Logger.Error(errors.EsCreateError.Error(), resource.MERCHANDISEMODULE)
		return err
	}

	m.DB.Create(&model.MerchantRecord{
		ID:     tools.CreateID(),
		Time:   time.Now(),
		UserID: uid,
		Action: rules.ADD,
		Source: store.MERCHANDISESTYLE,
	})

	return stream.SendMsg(&merchandise.AddMerchandiseStyleResp{
		Code:    rsp.OK,
		Message: rsp.CREATESUCCESS,
	})
}

func (m *MerchandiseStyle) RemoveMerchandiseStyle(ctx context.Context, req *merchandise.RemoveMerchandiseStyleReq, resp *merchandise.RemoveMerchandiseStyleResp) error {
	uid := ctx.Value("user_id").(string)
	if err := m.ES[store.MERCHANDISESTYLE].DeleteDocument(req.GetId()); err != nil {
		m.Logger.Error(errors.EsDeleteError.Error(), resource.MERCHANDISEMODULE)
		return err
	}

	m.DB.Create(&model.MerchantRecord{
		ID:     tools.CreateID(),
		Time:   time.Now(),
		UserID: uid,
		Action: rules.DELETE,
		Source: store.MERCHANDISESTYLE,
	})

	resp.Code = rsp.OK
	resp.Message = rsp.CREATESUCCESS

	return nil
}

func (m *MerchandiseStyle) UpdateMerchandiseStyle(ctx context.Context, stream merchandise.MerchandiseService_UpdateMerchandiseStyleStream) error {
	uid := ctx.Value("user_id").(string)

	var req *merchandise.UpdateMerchandiseStyleReq
	var picture bytes.Buffer

	for {
		data, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		req = &merchandise.UpdateMerchandiseStyleReq{
			Id:     data.GetId(),
			Name:   data.GetName(),
			Info:   data.GetInfo(),
			Price:  data.GetPrice(),
			Stock:  data.GetStock(),
			Status: data.GetStatus(),
		}

		chunk := data.GetChunk()
		_, err = picture.Write(chunk.GetData())
		if err != nil {
			return err
		}
	}

	if picture.Len() != 0 {
		path := fmt.Sprintf("%s", req.GetId())
		_, err := m.MC.PutObject(m.Ctx, store.MERCHANDISESTYLE, path, &picture, int64(picture.Len()), minio.PutObjectOptions{})
		if err != nil {
			return err
		}
	}

	var queries map[string]interface{}
	if req.GetName() != "" {
		queries["name"] = req.GetName()
	}

	if req.GetInfo() != "" {
		queries["info"] = req.GetInfo()
	}

	if req.GetPrice() != -1 {
		queries["price"] = req.GetPrice()
	}

	if req.GetStatus() != "" {
		queries["status"] = req.GetStatus()
	}

	if req.GetStock() != -1 {
		if req.GetStock() < 0 {
			m.Logger.Error(errors.UndefinedValue("stock").Error(), resource.MERCHANDISEMODULE)
		}
		queries["stock"] = req.GetStock()
	}

	queries["updated_at"] = time.Now().Unix()

	if err := m.ES[store.MERCHANDISESTYLE].Update(req.GetId(), queries); err != nil {
		m.Logger.Error(errors.EsUpdateError.Error(), resource.MERCHANDISEMODULE)
		return err
	}

	m.DB.Create(&model.MerchantRecord{
		ID:     tools.CreateID(),
		Time:   time.Now(),
		UserID: uid,
		Action: rules.UPDATE,
		Source: store.MERCHANDISESTYLE,
	})

	return stream.SendMsg(&merchandise.UpdateMerchandiseStyleResp{
		Code:    rsp.OK,
		Message: rsp.UPDATESUCCESS,
	})
}

func (m *MerchandiseStyle) GetMerchandiseStyleList(ctx context.Context, req *merchandise.GetMerchandiseStyleListReq, resp *merchandise.GetMerchandiseStyleListResp) error {
	must := []types.Query{
		{
			MatchPhrase: map[string]types.MatchPhraseQuery{
				"merchandise_id": {Query: req.GetMerchandiseID()},
			},
		},
	}
	response, err := m.ES[store.MERCHANDISESTYLE].Search(must, nil, nil, int(req.GetReq()*req.GetSize()), int(req.GetSize()))
	if err != nil {
		m.Logger.Error(errors.EsSearchError.Error(), resource.MERCHANDISEMODULE)
		return err
	}

	var ms []model.MerchandiseStyle
	for _, v := range response.Hits.Hits {
		var style model.MerchandiseStyle
		if err := json.Unmarshal(v.Source_, &style); err != nil {
			m.Logger.Error(errors.JsonUnmarshalError.Error(), resource.MERCHANDISEMODULE)
			return err
		}

		ms = append(ms, style)
	}

	data, err := json.Marshal(&ms)
	if err != nil {
		m.Logger.Error(errors.JsonMarshalError.Error(), resource.MERCHANDISEMODULE)
		return err
	}

	resp.Code = rsp.OK
	resp.Data = data
	resp.Message = rsp.SEARCHSUCCESS

	return nil
}

func (m *MerchandiseStyle) GetMerchandiseStyleDetails(ctx context.Context, req *merchandise.GetMerchandiseStyleDetailsReq, resp *merchandise.GetMerchandiseStyleDetailsResp) error {
	data, err := m.ES[store.MERCHANDISESTYLE].GetDocumentByID(req.GetId())
	if err != nil {
		m.Logger.Error(errors.EsSearchError.Error(), resource.MERCHANDISEMODULE)
		return err
	}

	resp.Code = rsp.OK
	resp.Data = data
	resp.Message = rsp.SEARCHSUCCESS

	return nil
}
