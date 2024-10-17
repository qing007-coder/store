package main

import (
	"fmt"
	"github.com/IBM/sarama"
	"store/pkg/kafka"
	"store/pkg/logger"
)

func main() {
	conf := sarama.NewConfig()
	conf.Net.MaxOpenRequests = 2
	setting := kafka.Setting{
		Addr: []string{"192.168.152.128:9092"},
		Conf: conf,
		SuccessHandler: func(msg *sarama.ProducerMessage) {
			fmt.Println("success:", msg.Value)
		},
		ErrorHandler: func(err error) {
			fmt.Println("err:", err)
		},
	}
	p, err := kafka.NewProducer(&setting)
	if err != nil {
		fmt.Println("err:", err)
		return
	}

	l := logger.NewLogger(p)
	l.Info("id:12,action:access", "gateway")
	fmt.Println("success")
	select {}
}
