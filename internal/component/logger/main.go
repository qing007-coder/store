package main

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"store/internal/component/logger/consumer"
	"store/pkg/elasticsearch"
	"store/pkg/kafka"
)

func main() {
	es, err := elasticsearch.NewClient(context.Background(), "http://192.168.152.128:9200", "log")
	_ = es.CreateIndex()
	if err != nil {
		fmt.Println("err:", err)
		return
	}

	conf := sarama.NewConfig()
	conf.Net.MaxOpenRequests = 2
	setting := kafka.Setting{
		Addr: []string{"192.168.152.128:9092"},
		Conf: conf,
	}
	c, err := kafka.NewConsumer(&setting)
	if err != nil {
		fmt.Println("err:", err)
		return
	}

	worker := consumer.NewLogConsumer(es, c)
	if err := worker.Run(); err != nil {
		fmt.Println("err:", err)
		return
	}
}
