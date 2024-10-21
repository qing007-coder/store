package base

import (
	"gorm.io/gorm"
	"store/pkg/config"
	"store/pkg/elasticsearch"
	"store/pkg/logger"
	"store/pkg/redis"
	"store/pkg/rules"
)

type Base struct {
	db       *gorm.DB
	rdb      *redis.Client
	enforcer *rules.Enforcer
	logger   *logger.Logger
	conf     *config.GlobalConfig
	es       *elasticsearch.Elasticsearch
}

func NewBase(db *gorm.DB, rdb *redis.Client, e *rules.Enforcer, l *logger.Logger, conf *config.GlobalConfig, es *elasticsearch.Elasticsearch) *Base {
	return &Base{
		db:       db,
		rdb:      rdb,
		enforcer: e,
		logger:   l,
		conf:     conf,
		es:       es,
	}
}
