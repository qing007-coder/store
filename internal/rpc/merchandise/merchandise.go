package merchandise

import (
	"context"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
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

type Merchandise struct {
	*base.Base
	*MerchandiseStyle
	*Picture
}

func NewMerchandise(b *base.Base) *Merchandise {
	return &Merchandise{
		b,
		NewMerchandiseStyle(b),
		NewPicture(b),
	}
}

func (m *Merchandise) PutAwayMerchandise(ctx context.Context, req *merchandise.PutAwayMerchandiseReq, resp *merchandise.PutAwayMerchandiseResp) error {
	uid := ctx.Value("user_id").(string)

	if err := m.ES[store.MERCHANDISE].CreateDocument(&model.Merchandise{
		ID:          req.GetId(),
		Name:        req.GetName(),
		Info:        req.GetInfo(),
		PictureList: req.GetPictureList(),
		MerchantID:  uid,
		Delivery:    req.GetDelivery(),
		Category:    req.GetCategory(),
		CreateAt:    time.Now().Unix(),
		Views:       0,
		SalesVolume: 0,
	}, req.GetId()); err != nil {
		m.Logger.Error(errors.EsCreateError.Error(), resource.MERCHANDISEMODULE)
		return err
	}

	m.DB.Create(&model.MerchantRecord{
		ID:     tools.CreateID(),
		Time:   time.Now(),
		UserID: uid,
		Action: rules.ADD,
		Source: store.MERCHANDISE,
	})

	resp.Code = rsp.OK
	resp.Message = rsp.CREATESUCCESS

	return nil
}

func (m *Merchandise) RemoveMerchandise(ctx context.Context, req *merchandise.RemoveMerchandiseReq, resp *merchandise.RemoveMerchandiseResp) error {
	uid := ctx.Value("user_id").(string)
	if err := m.ES[store.MERCHANDISE].DeleteDocument(req.Id); err != nil {
		m.Logger.Error(errors.EsDeleteError.Error(), resource.MERCHANDISEMODULE)
		return err
	}

	m.DB.Create(&model.MerchantRecord{
		ID:     tools.CreateID(),
		Time:   time.Now(),
		UserID: uid,
		Action: rules.DELETE,
		Source: store.MERCHANDISE,
	})

	resp.Code = rsp.OK
	resp.Message = rsp.DELETESUCCESS

	return nil
}

func (m *Merchandise) UpdateMerchandise(ctx context.Context, req *merchandise.UpdateMerchandiseReq, resp *merchandise.UpdateMerchandiseResp) error {
	uid := ctx.Value("user_id").(string)

	queries := make(map[string]interface{})

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

	queries["picture_list"] = req.GetPictureList()
	queries["updated_at"] = time.Now().Unix()

	if err := m.ES[store.MERCHANDISE].Update(req.Id, queries); err != nil {
		m.Logger.Error(errors.EsUpdateError.Error(), resource.MERCHANDISEMODULE)
		return err
	}

	m.DB.Create(&model.MerchantRecord{
		ID:     tools.CreateID(),
		Time:   time.Now(),
		UserID: uid,
		Action: rules.UPDATE,
		Source: store.MERCHANDISE,
	})

	resp.Code = rsp.OK
	resp.Message = rsp.UPDATESUCCESS

	return nil
}

func (m *Merchandise) GetMerchandiseDetails(ctx context.Context, req *merchandise.GetMerchandiseDetailsReq, resp *merchandise.GetMerchandiseDetailsResp) error {
	data, err := m.ES[store.MERCHANDISE].GetDocumentByID(req.GetId())
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

	switch req.GetTime() {
	case 0:
		//sort = append(sort, map[string]interface{}{"field": "create_at", "order": "desc"})
		sort = append(sort, map[string]interface{}{
			"create_at": map[string]interface{}{
				"order": "desc",
			},
		})
	case 1:
		//sort = append(sort, map[string]interface{}{"field": "create_at", "order": "asc"})
		sort = append(sort, map[string]interface{}{
			"create_at": map[string]interface{}{
				"order": "asc",
			},
		})
	default:
		m.Logger.Error(errors.UndefinedValue("time").Error(), resource.MERCHANDISEMODULE)
		return errors.New("未知time值")
	}

	switch req.GetSales() {
	case 0:
		//sort = append(sort, map[string]interface{}{"field": "sales_volume", "order": "desc"})
		sort = append(sort, map[string]interface{}{
			"sales_volume": map[string]interface{}{
				"order": "desc",
			},
		})
	case 1:
		//sort = append(sort, map[string]interface{}{"field": "sales_volume", "order": "asc"})
		sort = append(sort, map[string]interface{}{
			"sales_volume": map[string]interface{}{
				"order": "asc",
			},
		})
	default:
		m.Logger.Error(errors.UndefinedValue("sales").Error(), resource.MERCHANDISEMODULE)
		return errors.New("未知sales值")
	}

	//if req.GetTime() == 0 {
	//	sort = append(sort, map[string]interface{}{
	//		"field": "create_at",
	//		"order": "desc",
	//	})
	//} else if req.GetSales() == 1 {
	//	sort = append(sort, map[string]interface{}{
	//		"field": "create_at",
	//		"order": "asc",
	//	})
	//} else {
	//	m.Logger.Error(errors.UndefinedValue("time").Error(), resource.MERCHANDISEMODULE)
	//	return errors.New("未知time值")
	//}
	//
	//if req.GetSales() == 0 {
	//	sort = append(sort, map[string]interface{}{
	//		"field": "sales_volume",
	//		"order": "desc",
	//	})
	//} else if req.GetSales() == 1 {
	//	sort = append(sort, map[string]interface{}{
	//		"field": "sales_volume",
	//		"order": "asc",
	//	})
	//} else {
	//	m.Logger.Error(errors.UndefinedValue("sales").Error(), resource.MERCHANDISEMODULE)
	//	return errors.New("未知sales值")
	//}

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

	response, err := m.ES[store.MERCHANDISE].Search(nil, shouldQueries, sort, int(req.Req*req.Size), int(req.Size))
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

	switch req.GetTime() {
	case 0:
		//sort = append(sort, map[string]interface{}{"field": "create_at", "order": "desc"})
		sort = append(sort, map[string]interface{}{
			"create_at": map[string]interface{}{
				"order": "desc",
			},
		})
	case 1:
		//sort = append(sort, map[string]interface{}{"field": "create_at", "order": "asc"})
		sort = append(sort, map[string]interface{}{
			"create_at": map[string]interface{}{
				"order": "asc",
			},
		})
	default:
		m.Logger.Error(errors.UndefinedValue("time").Error(), resource.MERCHANDISEMODULE)
		return errors.New("未知time值")
	}

	switch req.GetSales() {
	case 0:
		//sort = append(sort, map[string]interface{}{"field": "sales_volume", "order": "desc"})
		sort = append(sort, map[string]interface{}{
			"sales_volume": map[string]interface{}{
				"order": "desc",
			},
		})
	case 1:
		//sort = append(sort, map[string]interface{}{"field": "sales_volume", "order": "asc"})
		sort = append(sort, map[string]interface{}{
			"sales_volume": map[string]interface{}{
				"order": "asc",
			},
		})
	default:
		m.Logger.Error(errors.UndefinedValue("sales").Error(), resource.MERCHANDISEMODULE)
		return errors.New("未知sales值")
	}

	//if req.GetTime() == 0 {
	//	sort = append(sort, map[string]interface{}{
	//		"field": "create_at",
	//		"order": "desc",
	//	})
	//} else if req.GetSales() == 1 {
	//	sort = append(sort, map[string]interface{}{
	//		"field": "create_at",
	//		"order": "asc",
	//	})
	//} else {
	//	m.Logger.Error(errors.UndefinedValue("time").Error(), resource.MERCHANDISEMODULE)
	//	return errors.New("未知time值")
	//}
	//
	//if req.GetSales() == 0 {
	//	sort = append(sort, map[string]interface{}{
	//		"field": "sales_volume",
	//		"order": "desc",
	//	})
	//} else if req.GetSales() == 1 {
	//	sort = append(sort, map[string]interface{}{
	//		"field": "sales_volume",
	//		"order": "asc",
	//	})
	//} else {
	//	m.Logger.Error(errors.UndefinedValue("sales").Error(), resource.MERCHANDISEMODULE)
	//	return errors.New("未知sales值")
	//}

	queries := []types.Query{
		{
			MatchPhrase: map[string]types.MatchPhraseQuery{
				"category": {Query: req.GetCategory()},
			},
		},
	}

	response, err := m.ES[store.MERCHANDISE].Search(queries, nil, sort, int(req.Req*req.Size), int(req.Size))
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
