package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-micro.dev/v4/metadata"
	"io"
	"store/internal/proto/merchandise"
	rsp "store/pkg/constant/response"
	"store/pkg/constant/rules"
	"store/pkg/constant/store"
	"store/pkg/errors"
	"store/pkg/model"
	"store/pkg/model/request"
	"store/pkg/tools"
	"strconv"
)

type MerchandiseApi struct {
	ctx    context.Context
	srv    *Service
	client merchandise.MerchandiseService
}

func NewMerchantApi(srv *Service) *MerchandiseApi {
	m := new(MerchandiseApi)
	m.init(srv)

	return m
}

func (m *MerchandiseApi) init(srv *Service) {
	m.ctx = context.Background()
	m.srv = srv
	m.client = merchandise.NewMerchandiseService("merchandise", m.srv.Client())
}

func (m *MerchandiseApi) PutAwayMerchandise(ctx *gin.Context) {
	c := metadata.Set(m.ctx, "Authorization", ctx.GetString("Authorization"))

	name := ctx.PostForm("name")
	info := ctx.PostForm("info")
	delivery := ctx.PostForm("delivery")
	category := ctx.PostForm("category")
	total, err := strconv.Atoi(ctx.PostForm("picture_total"))
	if err != nil {
		tools.BadRequest(ctx, err.Error())
		return
	}
	form, err := ctx.MultipartForm()
	if err != nil {
		tools.BadRequest(ctx, err.Error())
		return
	}

	pictures := make(map[string][]byte)
	var paths []string
	id := tools.CreateID()

	for i := 0; i < total; i++ {
		file, err := form.File[strconv.Itoa(i)][0].Open()
		if err != nil {
			tools.BadRequest(ctx, err.Error())
			return
		}

		data, err := io.ReadAll(file)
		if err != nil {
			tools.BadRequest(ctx, err.Error())
			return
		}

		path := fmt.Sprintf("%s/%d", id, i)
		paths = append(paths, path)
		pictures[path] = data
	}

	errChan := make(chan error)
	go func() {
		defer close(errChan)
		stream, err := m.client.HandlePicture(c)
		if err != nil {
			errChan <- err
			return
		}

		for path, data := range pictures {
			chunks := tools.SplitBytes(data, rules.CHUNK)

			for _, chunk := range chunks {
				if err := stream.Send(&merchandise.HandlePicturesReq{
					Bucket: store.MERCHANDISE,
					Path:   path,
					Data:   chunk,
				}); err != nil {
					fmt.Println("err:", err)
					errChan <- err
					return
				}
			}
		}

		if err := stream.CloseSend(); err != nil {
			fmt.Println("err1", err)
			errChan <- err
			return
		}

		var resp merchandise.HandlePicturesResp
		if err := stream.RecvMsg(&resp); err != nil {
			errChan <- err
			return
		}

		errChan <- nil
	}()

	resp, err := m.client.PutAwayMerchandise(c, &merchandise.PutAwayMerchandiseReq{
		Id:          id,
		Name:        name,
		Info:        info,
		Delivery:    delivery,
		PictureList: paths,
		Category:    category,
	})

	if err != nil {
		tools.BadRequest(ctx, err.Error())
		return
	}

	if resp.GetCode() != rsp.OK {
		tools.BadRequest(ctx, resp.GetMessage())
		return
	}

	select {
	case err := <-errChan:
		if err != nil {
			tools.BadRequest(ctx, err.Error())
			return
		}
		tools.StatusOK(ctx, nil, resp.GetMessage())
	}
}

func (m *MerchandiseApi) RemoveMerchandise(ctx *gin.Context) {
	c := metadata.Set(m.ctx, "Authorization", ctx.GetString("Authorization"))
	var req request.RemoveMerchandiseReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		tools.BadRequest(ctx, errors.ShouldBindJsonError.Error())
		return
	}

	resp, err := m.client.RemoveMerchandise(c, &merchandise.RemoveMerchandiseReq{
		Id: req.ID,
	})

	if err != nil {
		tools.BadRequest(ctx, err.Error())
		return
	}

	if resp.GetCode() != rsp.OK {
		tools.BadRequest(ctx, resp.GetMessage())
		return
	}

	tools.StatusOK(ctx, nil, resp.GetMessage())
}

func (m *MerchandiseApi) UpdateMerchandise(ctx *gin.Context) {
	c := metadata.Set(m.ctx, "Authorization", ctx.GetString("Authorization"))

	id := ctx.PostForm("id")
	name := ctx.PostForm("name")
	info := ctx.PostForm("info")
	delivery := ctx.PostForm("delivery")
	category := ctx.PostForm("category")
	total, err := strconv.Atoi(ctx.PostForm("picture_total"))
	if err != nil {
		tools.BadRequest(ctx, err.Error())
		return
	}
	form, err := ctx.MultipartForm()
	if err != nil {
		tools.BadRequest(ctx, err.Error())
		return
	}

	pictures := make(map[string][]byte)
	var paths []string

	for i := 0; i < total; i++ {
		val, ok := form.File[strconv.Itoa(i)]
		if !ok {
			paths = append(paths, strconv.Itoa(i))
			continue
		}

		file, err := val[0].Open()
		if err != nil {
			tools.BadRequest(ctx, err.Error())
			return
		}

		data, err := io.ReadAll(file)
		if err != nil {
			tools.BadRequest(ctx, err.Error())
			return
		}

		path := fmt.Sprintf("%s/%d", id, i)
		paths = append(paths, path)
		pictures[path] = data
	}

	errChan := make(chan error)
	go func() {
		defer close(errChan)
		stream, err := m.client.HandlePicture(c)
		if err != nil {
			errChan <- err
			return
		}

		for path, data := range pictures {
			chunks := tools.SplitBytes(data, rules.CHUNK)

			for _, chunk := range chunks {
				if err := stream.Send(&merchandise.HandlePicturesReq{
					Bucket: store.MERCHANDISE,
					Path:   path,
					Data:   chunk,
				}); err != nil {
					errChan <- err
					return
				}
			}
		}

		if err := stream.CloseSend(); err != nil {
			errChan <- err
			return
		}

		var resp merchandise.HandlePicturesResp
		if err := stream.RecvMsg(&resp); err != nil {
			errChan <- err
			return
		}

		errChan <- nil
	}()

	resp, err := m.client.UpdateMerchandise(c, &merchandise.UpdateMerchandiseReq{
		Id:          id,
		Name:        name,
		Info:        info,
		Delivery:    delivery,
		PictureList: paths,
		Category:    category,
	})

	if err != nil {
		tools.BadRequest(ctx, err.Error())
		return
	}

	if resp.GetCode() != rsp.OK {
		tools.BadRequest(ctx, resp.GetMessage())
		return
	}

	select {
	case err := <-errChan:
		if err != nil {
			tools.BadRequest(ctx, err.Error())
			return
		}
		tools.StatusOK(ctx, nil, resp.GetMessage())
	}
}

func (m *MerchandiseApi) GetMerchandiseDetails(ctx *gin.Context) {
	id := ctx.Query("id")
	c := metadata.Set(m.ctx, "Authorization", ctx.GetString("Authorization"))

	resp, err := m.client.GetMerchandiseDetails(c, &merchandise.GetMerchandiseDetailsReq{
		Id: id,
	})

	if err != nil {
		tools.BadRequest(ctx, err.Error())
		return
	}

	if resp.GetCode() != rsp.OK {
		tools.BadRequest(ctx, resp.GetMessage())
		return
	}

	var md model.Merchandise
	if err := json.Unmarshal(resp.GetData(), &md); err != nil {
		tools.BadRequest(ctx, errors.JsonUnmarshalError.Error())
		return
	}

	tools.StatusOK(ctx, gin.H{
		"data": md,
	}, resp.GetMessage())
}

func (m *MerchandiseApi) Search(ctx *gin.Context) {
	content := ctx.Query("content")
	time, _ := strconv.Atoi(ctx.Query("time"))
	sales, _ := strconv.Atoi(ctx.Query("sales"))
	req, _ := strconv.Atoi(ctx.Query("req"))
	size, _ := strconv.Atoi(ctx.Query("size"))

	c := metadata.Set(m.ctx, "Authorization", ctx.GetString("Authorization"))
	resp, err := m.client.Search(c, &merchandise.SearchReq{
		Content: content,
		Time:    int32(time),
		Sales:   int32(sales),
		Size:    int32(size),
		Req:     int32(req),
	})

	if err != nil {
		tools.BadRequest(ctx, err.Error())
		return
	}

	if resp.GetCode() != rsp.OK {
		tools.BadRequest(ctx, resp.GetMessage())
		return
	}
	var ms []model.Merchandise
	if err := json.Unmarshal(resp.GetData(), &ms); err != nil {
		tools.BadRequest(ctx, errors.JsonUnmarshalError.Error())
		return
	}

	tools.StatusOK(ctx, gin.H{
		"data": ms,
	}, resp.GetMessage())
}

func (m *MerchandiseApi) SearchByCategory(ctx *gin.Context) {
	category := ctx.Query("category")
	time, _ := strconv.Atoi(ctx.Query("time"))
	sales, _ := strconv.Atoi(ctx.Query("sales"))
	req, _ := strconv.Atoi(ctx.Query("req"))
	size, _ := strconv.Atoi(ctx.Query("size"))

	c := metadata.Set(m.ctx, "Authorization", ctx.GetString("Authorization"))
	resp, err := m.client.SearchByCategory(c, &merchandise.SearchByCategoryReq{
		Category: category,
		Time:     int32(time),
		Sales:    int32(sales),
		Size:     int32(size),
		Req:      int32(req),
	})

	if err != nil {
		tools.BadRequest(ctx, err.Error())
		return
	}

	if resp.GetCode() != rsp.OK {
		tools.BadRequest(ctx, resp.GetMessage())
		return
	}

	var ms []model.Merchandise
	if err := json.Unmarshal(resp.GetData(), &ms); err != nil {
		tools.BadRequest(ctx, errors.JsonUnmarshalError.Error())
		return
	}

	tools.StatusOK(ctx, gin.H{
		"data": ms,
	}, resp.GetMessage())
}

func (m *MerchandiseApi) AddMerchandiseStyle(ctx *gin.Context) {
	c := metadata.Set(m.ctx, "Authorization", ctx.GetString("Authorization"))
	merchandiseID := ctx.PostForm("merchandise_id")
	name := ctx.PostForm("name")
	info := ctx.PostForm("info")
	price, err := strconv.ParseFloat(ctx.PostForm("price"), 32)
	if err != nil {
		tools.BadRequest(ctx, err.Error())
		return
	}
	stock, err := strconv.Atoi(ctx.PostForm("stock"))
	if err != nil {
		tools.BadRequest(ctx, err.Error())
		return
	}
	status := ctx.PostForm("status")

	file, err := ctx.FormFile("picture")
	if err != nil {
		tools.BadRequest(ctx, err.Error())
		return
	}
	f, _ := file.Open()
	data, err := io.ReadAll(f)
	if err != nil {
		tools.BadRequest(ctx, err.Error())
		return
	}

	id := tools.CreateID()
	errChan := make(chan error)

	go func() {
		chunks := tools.SplitBytes(data, rules.CHUNK)
		stream, err := m.client.HandlePicture(c)
		if err != nil {
			errChan <- err
			return
		}
		for _, chunk := range chunks {
			if err := stream.Send(&merchandise.HandlePicturesReq{
				Bucket: store.MERCHANDISESTYLE,
				Path:   id,
				Data:   chunk,
			}); err != nil {
				errChan <- err
				return
			}
		}

		if err := stream.CloseSend(); err != nil {
			errChan <- err
			return
		}

		var resp merchandise.HandlePicturesResp
		if err := stream.RecvMsg(&resp); err != nil {
			errChan <- err
			return
		}

		errChan <- nil
	}()

	resp, err := m.client.AddMerchandiseStyle(c, &merchandise.AddMerchandiseStyleReq{
		Id:            id,
		MerchandiseID: merchandiseID,
		Name:          name,
		Info:          info,
		Picture:       id,
		Price:         float32(price),
		Status:        status,
		Stock:         uint32(stock),
	})

	if err != nil {
		tools.BadRequest(ctx, err.Error())
		return
	}

	if resp.GetCode() != rsp.OK {
		tools.BadRequest(ctx, resp.GetMessage())
		return
	}

	select {
	case err := <-errChan:
		if err != nil {
			tools.BadRequest(ctx, err.Error())
			return
		}
		tools.StatusOK(ctx, nil, resp.GetMessage())
	}
}

func (m *MerchandiseApi) RemoveMerchandiseStyle(ctx *gin.Context) {
	var req request.RemoveMerchandiseStyleReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		tools.BadRequest(ctx, errors.ShouldBindJsonError.Error())
		return
	}
	c := metadata.Set(m.ctx, "Authorization", ctx.GetString("Authorization"))

	resp, err := m.client.RemoveMerchandiseStyle(c, &merchandise.RemoveMerchandiseStyleReq{
		Id: req.ID,
	})

	if err != nil {
		tools.BadRequest(ctx, err.Error())
		return
	}

	if resp.GetCode() != rsp.OK {
		tools.BadRequest(ctx, resp.GetMessage())
		return
	}

	tools.StatusOK(ctx, nil, resp.GetMessage())
}

func (m *MerchandiseApi) UpdateMerchandiseStyle(ctx *gin.Context) {
	c := metadata.Set(m.ctx, "Authorization", ctx.GetString("Authorization"))

	id := ctx.PostForm("id")
	name := ctx.PostForm("name")
	info := ctx.PostForm("info")

	var price float64

	if p := ctx.PostForm("price"); p != "" {
		var err error
		price, err = strconv.ParseFloat(p, 32)
		if err != nil {
			tools.BadRequest(ctx, err.Error())
			return
		}
	}

	var stock int

	if s := ctx.PostForm("stock"); s != "" {
		var err error
		stock, err = strconv.Atoi(ctx.PostForm("stock"))
		if err != nil {
			tools.BadRequest(ctx, err.Error())
			return
		}
	}

	status := ctx.PostForm("status")

	file, err := ctx.FormFile("picture")
	if err == nil {
		f, err := file.Open()
		if err != nil {
			tools.BadRequest(ctx, err.Error())
			return
		}
		data, err := io.ReadAll(f)
		if err != nil {
			tools.BadRequest(ctx, err.Error())
			return
		}

		stream, err := m.client.HandlePicture(c)
		if err != nil {
			tools.BadRequest(ctx, err.Error())
			return
		}

		chunks := tools.SplitBytes(data, rules.CHUNK)
		for _, chunk := range chunks {
			if err := stream.Send(&merchandise.HandlePicturesReq{
				Bucket: store.MERCHANDISESTYLE,
				Path:   id,
				Data:   chunk,
			}); err != nil {
				tools.BadRequest(ctx, err.Error())
				return
			}
		}

		if err := stream.CloseSend(); err != nil {
			tools.BadRequest(ctx, err.Error())
			return
		}
	}

	resp, err := m.client.UpdateMerchandiseStyle(c, &merchandise.UpdateMerchandiseStyleReq{
		Id:      id,
		Name:    name,
		Info:    info,
		Picture: id,
		Price:   float32(price),
		Status:  status,
		Stock:   int32(stock),
	})

	if err != nil {
		tools.BadRequest(ctx, err.Error())
		return
	}

	if resp.GetCode() != rsp.OK {
		tools.BadRequest(ctx, resp.GetMessage())
		return
	}

	tools.StatusOK(ctx, nil, resp.GetMessage())
}

func (m *MerchandiseApi) GetMerchandiseStyleList(ctx *gin.Context) {
	c := metadata.Set(m.ctx, "Authorization", ctx.GetString("Authorization"))

	merchandiseID := ctx.Query("merchandise_id")
	req, err := strconv.Atoi(ctx.Query("req"))
	if err != nil {
		tools.BadRequest(ctx, err.Error())
		return
	}
	size, err := strconv.Atoi(ctx.Query("size"))
	if err != nil {
		tools.BadRequest(ctx, err.Error())
		return
	}

	resp, err := m.client.GetMerchandiseStyleList(c, &merchandise.GetMerchandiseStyleListReq{
		MerchandiseID: merchandiseID,
		Req:           int32(req),
		Size:          int32(size),
	})

	if err != nil {
		tools.BadRequest(ctx, err.Error())
		return
	}

	if resp.GetCode() != rsp.OK {
		tools.BadRequest(ctx, resp.GetMessage())
		return
	}

	var ms []model.MerchandiseStyle
	if err := json.Unmarshal(resp.GetData(), &ms); err != nil {
		tools.BadRequest(ctx, errors.JsonUnmarshalError.Error())
		return
	}

	tools.StatusOK(ctx, gin.H{
		"data": ms,
	}, resp.GetMessage())
}

func (m *MerchandiseApi) GetMerchandiseStyleDetails(ctx *gin.Context) {
	id := ctx.Query("id")
	c := metadata.Set(m.ctx, "Authorization", ctx.GetString("Authorization"))

	resp, err := m.client.GetMerchandiseStyleDetails(c, &merchandise.GetMerchandiseStyleDetailsReq{
		Id: id,
	})

	if err != nil {
		tools.BadRequest(ctx, err.Error())
		return
	}

	if resp.GetCode() != rsp.OK {
		tools.BadRequest(ctx, resp.GetMessage())
		return
	}

	var ms model.MerchandiseStyle
	if err := json.Unmarshal(resp.GetData(), &ms); err != nil {
		tools.BadRequest(ctx, errors.JsonUnmarshalError.Error())
		return
	}

	tools.StatusOK(ctx, gin.H{
		"data": ms,
	}, resp.GetMessage())
}
