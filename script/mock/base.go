package mock

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"store/pkg/config"
	"store/pkg/elasticsearch"
	"store/pkg/mysql"
)

type Base struct {
	ctx context.Context
	db  *gorm.DB
	es  map[string]*elasticsearch.Elasticsearch
}

func NewBase(index []string) (*Base, error) {
	ctx := context.Background()
	conf, err := config.NewGlobalConfig()
	if err != nil {
		return nil, err
	}

	db, err := mysql.NewClient(conf)
	if err != nil {
		return nil, err
	}
	es := make(map[string]*elasticsearch.Elasticsearch)
	for _, i := range index {
		client, err := elasticsearch.NewClient(ctx, fmt.Sprintf("%s:%s", conf.Elasticsearch.Addr, conf.Elasticsearch.Port), i)
		if err != nil {
			return nil, err
		}
		es[i] = client
	}

	return &Base{
		ctx: ctx,
		db:  db,
		es:  es,
	}, nil
}
