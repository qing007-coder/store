package main

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"store/internal/component/logger"
	"store/pkg/config"
	"store/pkg/elasticsearch"
	"store/pkg/errors"
	"store/pkg/kafka"
)

func main() {
	co, err := config.NewGlobalConfig()
	if err != nil {
		errors.HandleError(err)
		return
	}

	es, err := elasticsearch.NewClient(context.Background(), fmt.Sprintf("%s:%s", co.Elasticsearch.Addr, co.Elasticsearch.Port), "log")
	if err != nil {
		fmt.Println("err:", err)
		return
	}

	//if err := es.CreateIndex(); err != nil {
	//	errors.HandleError(err)
	//	return
	//}

	conf := sarama.NewConfig()
	conf.Net.MaxOpenRequests = 2
	setting := kafka.Setting{
		Addr: []string{fmt.Sprintf("%s:%s", co.Kafka.Addr, co.Kafka.Port)},
		Conf: conf,
	}
	c, err := kafka.NewConsumer(&setting)
	if err != nil {
		fmt.Println("err:", err)
		return
	}

	worker := logger.NewLogConsumer(es, c)
	if err := worker.Run(); err != nil {
		fmt.Println("err:", err)
		return
	}
}
