package gateway

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go-micro.dev/v4/metadata"
	"io"
	"store/internal/proto/merchandise"
	rsp "store/pkg/constant/response"
	"store/pkg/constant/rules"
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
	m.client = merchandise.NewMerchandiseService("merchant", m.srv.Client())
}

func (m *MerchandiseApi) PutAwayMerchandise(ctx *gin.Context) {
	name := ctx.PostForm("name")
	info := ctx.PostForm("info")
	delivery := ctx.PostForm("delivery")
	category := ctx.PostForm("category")

	files, err := ctx.MultipartForm()
	if err != nil {
		tools.BadRequest(ctx, err.Error())
		return
	}
	c := metadata.Set(m.ctx, "Authorization", ctx.GetString("Authorization"))

	stream, err := m.client.PutAwayMerchandise(c)
	if err != nil {
		tools.BadRequest(ctx, err.Error())
		return
	}

	for id, file := range files.File["pictures"] {
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

		chunks := tools.SplitBytes(data, rules.CHUNK)
		for _, chunk := range chunks {
			if err := stream.Send(&merchandise.PutAwayMerchandiseReq{
				Name:     name,
				Info:     info,
				Delivery: delivery,
				Category: category,
				Chunk: &merchandise.Chunk{
					PictureID: strconv.Itoa(id),
					Data:      chunk,
				},
			}); err != nil {
				tools.BadRequest(ctx, err.Error())
				return
			}
		}
	}

	if err := stream.CloseSend(); err != nil {
		tools.BadRequest(ctx, err.Error())
		return
	}

	var resp merchandise.PutAwayMerchandiseResp
	if err := stream.RecvMsg(&resp); err != nil {
		tools.BadRequest(ctx, err.Error())
		return
	}

	if resp.GetCode() != rsp.OK {
		tools.BadRequest(ctx, resp.GetMessage())
		return
	}

	tools.StatusOK(ctx, nil, resp.GetMessage())
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

}

func (m *MerchandiseApi) RemoveMerchandiseStyle(ctx *gin.Context) {

}

func (m *MerchandiseApi) UpdateMerchandiseStyle(ctx *gin.Context) {

}

func (m *MerchandiseApi) GetMerchandiseStyleList(ctx *gin.Context) {

}

func (m *MerchandiseApi) GetMerchandiseStyleDetails(ctx *gin.Context) {

}
