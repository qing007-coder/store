package merchandise

import (
	"bytes"
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"io"
	"store/internal/proto/merchandise"
	"store/internal/rpc/base"
	"store/pkg/constant/resource"
	rsp "store/pkg/constant/response"
	"store/pkg/errors"
)

type Picture struct {
	*base.Base
}

func NewPicture(b *base.Base) *Picture {
	return &Picture{b}
}

func (p *Picture) HandlePicture(ctx context.Context, stream merchandise.MerchandiseService_HandlePictureStream) error {
	pictures := make(map[string]*bytes.Buffer)
	var bucket string

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}

		bucket = req.GetBucket()

		if pictures[req.GetPath()] == nil {
			pictures[req.GetPath()] = new(bytes.Buffer)
			fmt.Println(req.GetPath())
			fmt.Println(bucket)
		}

		pictures[req.GetPath()].Write(req.GetData())
	}

	for path, data := range pictures {
		_, err := p.MC.PutObject(ctx, bucket, path, data, int64(data.Len()), minio.PutObjectOptions{})
		if err != nil {
			fmt.Println("Err:", err)
			p.Logger.Error(errors.MCCreateError.Error(), resource.MERCHANDISEMODULE)
			return errors.MCCreateError
		}
	}

	return stream.SendMsg(&merchandise.HandlePicturesResp{
		Code:    rsp.OK,
		Message: rsp.CREATESUCCESS,
	})
}
