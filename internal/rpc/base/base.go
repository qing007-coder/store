package base

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"gorm.io/gorm"
	"store/pkg/config"
	"store/pkg/elasticsearch"
	"store/pkg/kafka"
	"store/pkg/logger"
	"store/pkg/mysql"
	"store/pkg/redis"
	"store/pkg/rules"
)

type Base struct {
	Ctx      context.Context
	DB       *gorm.DB
	RDB      *redis.Client
	Enforcer *rules.Enforcer
	Logger   *logger.Logger
	Conf     *config.GlobalConfig
	ES       map[string]*elasticsearch.Elasticsearch
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

	rdb := redis.NewClient(conf)
	e := rules.NewEnforcer(db)
	c := sarama.NewConfig()
	c.Net.MaxOpenRequests = 2
	setting := kafka.Setting{
		Addr: []string{fmt.Sprintf("%s:%s", conf.Kafka.Addr, conf.Kafka.Port)},
		Conf: c,
		SuccessHandler: func(msg *sarama.ProducerMessage) {
			fmt.Println("success:", msg.Value)
		},
		ErrorHandler: func(err error) {
			fmt.Println("err:", err)
		},
	}
	p, err := kafka.NewProducer(&setting)
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

	l := logger.NewLogger(p)
	return &Base{
		Ctx:      ctx,
		DB:       db,
		RDB:      rdb,
		Enforcer: e,
		Logger:   l,
		Conf:     conf,
		ES:       es,
	}, nil
}
