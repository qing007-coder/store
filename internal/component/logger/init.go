package logger

import (
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"log"
	"store/pkg/constant"
	"store/pkg/elasticsearch"
	"store/pkg/kafka"
	"store/pkg/model"
)

type LogConsumer struct {
	es       *elasticsearch.Elasticsearch
	consumer *kafka.Consumer
}

func NewLogConsumer(es *elasticsearch.Elasticsearch, c *kafka.Consumer) *LogConsumer {
	l := &LogConsumer{
		es:       es,
		consumer: c,
	}
	l.consumer.SetMessageHandler(l.HandleLog)
	l.consumer.SetErrorHandler(l.HandleError)

	return l
}

func (l *LogConsumer) Run() error {
	return l.consumer.Subscribe(constant.LOGTOPIC, 0, sarama.OffsetNewest)
}

func (l *LogConsumer) HandleLog(message *sarama.ConsumerMessage) {
	var lg model.Log
	fmt.Println(string(message.Value))
	if err := json.Unmarshal(message.Value, &lg); err != nil {
		fmt.Println("err:", err)
		return
	}
	fmt.Println(lg)
	if err := l.es.CreateDocument(&lg, lg.ID); err != nil {
		fmt.Println("err:", err)
		return
	}
}

func (l *LogConsumer) HandleError(err error) {
	log.Println("err:", err)
}
