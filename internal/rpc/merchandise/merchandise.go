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
	"store/pkg/constant"
	"store/pkg/constant/resource"
	rsp "store/pkg/constant/response"
	"store/pkg/errors"
	"store/pkg/model"
	"store/pkg/tools"
	"time"
)

type Merchandise struct {
	*base.Base
	*MerchandiseStyle
}

func NewMerchandise(b *base.Base) *Merchandise {
	return &Merchandise{
		b,
		NewMerchandiseStyle(b),
	}
}

func (m *Merchandise) PutAwayMerchandise(ctx context.Context, stream merchandise.MerchandiseService_PutAwayMerchandiseStream) error {
	uid := ctx.Value("user_id").(string)
	var req *merchandise.PutAwayMerchandiseReq
	pictures := make(map[string]*bytes.Buffer)

	for {
		data, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		req = &merchandise.PutAwayMerchandiseReq{
			Name:     data.GetName(),
			Info:     data.GetInfo(),
			Delivery: data.GetDelivery(),
			Category: data.GetCategory(),
		}

		chunk := data.GetChunk()
		_, err = pictures[chunk.GetPictureID()].Write(chunk.GetData())
		if err != nil {
			return err
		}
	}

	id := tools.CreateID()
	var pictureList []string

	for pictureID, data := range pictures {
		path := fmt.Sprintf("%s/%s", id, pictureID)
		pictureList = append(pictureList, path)

		_, err := m.MC.PutObject(m.Ctx, constant.MERCHANDISE, path, data, int64(data.Len()), minio.PutObjectOptions{})
		if err != nil {
			return err
		}
	}

	if err := m.ES[constant.MERCHANDISE].CreateDocument(&model.Merchandise{
		ID:          id,
		Name:        req.GetName(),
		Info:        req.GetInfo(),
		PictureList: pictureList,
		MerchantID:  uid,
		Delivery:    req.GetDelivery(),
		Category:    req.GetCategory(),
		CreateAt:    time.Now().Unix(),
		Views:       0,
		SalesVolume: 0,
	}, id); err != nil {
		m.Logger.Error(errors.EsCreateError.Error(), resource.MERCHANDISEMODULE)
		return err
	}

	m.DB.Create(&model.MerchantRecord{
		ID:     tools.CreateID(),
		Time:   time.Now(),
		UserID: uid,
		Action: constant.ADD,
		Source: constant.MERCHANDISE,
	})

	return stream.SendMsg(&merchandise.PutAwayMerchandiseResp{
		Code:    rsp.OK,
		Message: rsp.CREATESUCCESS,
	})
}

func (m *Merchandise) RemoveMerchandise(ctx context.Context, req *merchandise.RemoveMerchandiseReq, resp *merchandise.RemoveMerchandiseResp) error {
	uid := ctx.Value("user_id").(string)
	if err := m.ES[constant.MERCHANDISE].DeleteDocument(req.Id); err != nil {
		m.Logger.Error(errors.EsDeleteError.Error(), resource.MERCHANDISEMODULE)
		return err
	}

	m.DB.Create(&model.MerchantRecord{
		ID:     tools.CreateID(),
		Time:   time.Now(),
		UserID: uid,
		Action: constant.DELETE,
		Source: constant.MERCHANDISE,
	})

	resp.Code = rsp.OK
	resp.Message = rsp.DELETESUCCESS

	return nil
}

func (m *Merchandise) UpdateMerchandise(ctx context.Context, stream merchandise.MerchandiseService_UpdateMerchandiseStream) error {
	uid := ctx.Value("user_id").(string)

	var req *merchandise.UpdateMerchandiseReq
	pictures := make(map[string]*bytes.Buffer)

	for {
		data, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		req = &merchandise.UpdateMerchandiseReq{
			Id:       data.GetId(),
			Name:     data.GetName(),
			Info:     data.GetInfo(),
			Delivery: data.GetDelivery(),
			Category: data.GetCategory(),
		}

		chunk := data.GetChunk()
		_, err = pictures[chunk.GetPictureID()].Write(chunk.GetData())
		if err != nil {
			return err
		}
	}

	id := tools.CreateID()
	var pictureList []string
	for pictureID, data := range pictures {
		path := fmt.Sprintf("%s/%s", id, pictureID)
		pictureList = append(pictureList, path)

		if data.Len() == 0 {
			continue
		}
		_, err := m.MC.PutObject(m.Ctx, constant.MERCHANDISE, path, data, int64(data.Len()), minio.PutObjectOptions{})
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

	if req.GetCategory() != "" {
		queries["category"] = req.GetCategory()
	}

	if req.GetDelivery() != "" {
		queries["delivery"] = req.GetDelivery()
	}

	queries["picture_list"] = pictureList
	queries["updated_at"] = time.Now().Unix()

	if err := m.ES[constant.MERCHANDISE].Update(req.Id, queries); err != nil {
		m.Logger.Error(errors.EsUpdateError.Error(), resource.MERCHANDISEMODULE)
		return err
	}

	m.DB.Create(&model.MerchantRecord{
		ID:     tools.CreateID(),
		Time:   time.Now(),
		UserID: uid,
		Action: constant.UPDATE,
		Source: constant.MERCHANDISE,
	})

	return stream.SendMsg(&merchandise.UpdateMerchandiseResp{
		Code:    rsp.OK,
		Message: rsp.UPDATESUCCESS,
	})
}

func (m *Merchandise) GetMerchandiseDetails(ctx context.Context, req *merchandise.GetMerchandiseDetailsReq, resp *merchandise.GetMerchandiseDetailsResp) error {
	data, err := m.ES[constant.MERCHANDISE].GetDocumentByID(req.GetId())
	if err != nil {
		m.Logger.Error(errors.EsSearchError.Error(), resource.MERCHANDISEMODULE)
		return err
	}

	resp.Code = rsp.OK
	resp.Data = data
	resp.Message = rsp.SEARCHSUCCESS
	return nil
}

func (m *Merchandise) Search(ctx context.Context, req *merchandise.SearchReq, resp *merchandise.SearchResp) error {
	var sort []types.SortCombinations
	if req.GetTime() == 0 {
		sort = append(sort, map[string]interface{}{
			"field": "create_at",
			"order": "desc",
		})
	} else if req.GetSales() == 1 {
		sort = append(sort, map[string]interface{}{
			"field": "create_at",
			"order": "asc",
		})
	} else {
		m.Logger.Error(errors.UndefinedValue("time").Error(), resource.MERCHANDISEMODULE)
		return errors.New("未知time值")
	}

	if req.GetSales() == 0 {
		sort = append(sort, map[string]interface{}{
			"field": "sales_volume",
			"order": "desc",
		})
	} else if req.GetSales() == 1 {
		sort = append(sort, map[string]interface{}{
			"field": "sales_volume",
			"order": "asc",
		})
	} else {
		m.Logger.Error(errors.UndefinedValue("sales").Error(), resource.MERCHANDISEMODULE)
		return errors.New("未知sales值")
	}

	shouldQueries := []types.Query{
		{
			Match: map[string]types.MatchQuery{
				"name": {
					Query:     req.GetContent(),
					Fuzziness: "AUTO",
				},
			},
		},
		{
			Match: map[string]types.MatchQuery{
				"info": {
					Query:     req.GetContent(),
					Fuzziness: "AUTO",
				},
			},
		},
		{
			Match: map[string]types.MatchQuery{
				"category": {
					Query:     req.GetContent(),
					Fuzziness: "AUTO",
				},
			},
		},
	}

	response, err := m.ES[constant.MERCHANDISE].Search(nil, shouldQueries, sort, int(req.Req*req.Size), int(req.Size))
	if err != nil {
		m.Logger.Error(errors.EsSearchError.Error(), resource.MERCHANDISEMODULE)
		return err
	}

	var ms []model.Merchandise
	for _, hit := range response.Hits.Hits {
		var mer model.Merchandise
		if err := json.Unmarshal(hit.Source_, &mer); err != nil {
			m.Logger.Error(errors.JsonUnmarshalError.Error(), resource.MERCHANDISEMODULE)
			return err
		}

		ms = append(ms, mer)
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

func (m *Merchandise) SearchByCategory(ctx context.Context, req *merchandise.SearchByCategoryReq, resp *merchandise.SearchByCategoryResp) error {
	var sort []types.SortCombinations
	if req.GetTime() == 0 {
		sort = append(sort, map[string]interface{}{
			"field": "create_at",
			"order": "desc",
		})
	} else if req.GetSales() == 1 {
		sort = append(sort, map[string]interface{}{
			"field": "create_at",
			"order": "asc",
		})
	} else {
		m.Logger.Error(errors.UndefinedValue("time").Error(), resource.MERCHANDISEMODULE)
		return errors.New("未知time值")
	}

	if req.GetSales() == 0 {
		sort = append(sort, map[string]interface{}{
			"field": "sales_volume",
			"order": "desc",
		})
	} else if req.GetSales() == 1 {
		sort = append(sort, map[string]interface{}{
			"field": "sales_volume",
			"order": "asc",
		})
	} else {
		m.Logger.Error(errors.UndefinedValue("sales").Error(), resource.MERCHANDISEMODULE)
		return errors.New("未知sales值")
	}

	queries := []types.Query{
		{
			MatchPhrase: map[string]types.MatchPhraseQuery{
				"category": {Query: req.GetCategory()},
			},
		},
	}

	response, err := m.ES[constant.MERCHANDISE].Search(queries, nil, sort, int(req.Req*req.Size), int(req.Size))
	if err != nil {
		m.Logger.Error(errors.EsSearchError.Error(), resource.MERCHANDISEMODULE)
		return err
	}

	var ms []model.Merchandise
	for _, v := range response.Hits.Hits {
		var mer model.Merchandise
		if err := json.Unmarshal(v.Source_, &mer); err != nil {
			m.Logger.Error(errors.JsonUnmarshalError.Error(), resource.MERCHANDISEMODULE)
			return err
		}

		ms = append(ms, mer)
	}

	data, err := json.Marshal(&ms)
	if err != nil {
		m.Logger.Error(errors.JsonMarshalError.Error(), resource.MERCHANDISEMODULE)
	}

	resp.Code = rsp.OK
	resp.Data = data
	resp.Message = rsp.SEARCHSUCCESS

	return nil
}
