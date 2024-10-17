package kafka

import (
	"github.com/IBM/sarama"
)

type Consumer struct {
	consumer      sarama.Consumer
	handleMessage func(*sarama.ConsumerMessage)
	handleError   func(error)
}

func NewConsumer(s *Setting) (*Consumer, error) {
	c := new(Consumer)

	c.SetMessageHandler(s.MessageHandler)
	c.SetErrorHandler(s.ErrorHandler)
	if err := c.init(s.Addr, s.Conf); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Consumer) init(addr []string, conf *sarama.Config) error {
	consumer, err := sarama.NewConsumer(addr, conf)
	if err != nil {
		return err
	}

	c.consumer = consumer
	return nil
}

func (c *Consumer) SetMessageHandler(f func(*sarama.ConsumerMessage)) {
	c.handleMessage = f
}

func (c *Consumer) SetErrorHandler(f func(error)) {
	c.handleError = f
}

func (c *Consumer) Subscribe(topic string, partition int32, offset int64) error {
	worker, err := c.consumer.ConsumePartition(topic, partition, offset)
	//defer worker.Close()
	if err != nil {
		return err
	}

	c.Handler(worker)
	return nil
}

func (c *Consumer) Handler(worker sarama.PartitionConsumer) {
	for {
		select {
		case message := <-worker.Messages():
			c.handleMessage(message)
		case err := <-worker.Errors():
			c.handleError(err.Err)
		}
	}
}
