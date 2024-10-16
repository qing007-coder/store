package kafka

import (
	"github.com/IBM/sarama"
)

type Producer struct {
	producer      sarama.AsyncProducer
	handleSuccess func(*sarama.ProducerMessage)
	handleError   func(error)
}

func NewProducer(addr []string, conf *sarama.Config) (*Producer, error) {
	p := new(Producer)
	if err := p.init(addr, conf); err != nil {
		return nil, err
	}

	return p, nil
}

func (p *Producer) init(addr []string, conf *sarama.Config) error {
	producer, err := sarama.NewAsyncProducer(addr, conf)
	if err != nil {
		return err
	}

	p.producer = producer
	go p.handle()
	return nil
}

func (p *Producer) SetSuccessHandler(f func(*sarama.ProducerMessage)) {
	p.handleSuccess = f
}

func (p *Producer) SetErrorHandler(f func(error)) {
	p.handleError = f
}

func (p *Producer) Publish(message *sarama.ProducerMessage) {
	worker := p.producer.Input()
	worker <- message
}

func (p *Producer) handle() {
	for {
		select {
		case success := <-p.producer.Successes():
			p.handleSuccess(success)
		case err := <-p.producer.Errors():
			p.handleError(err.Err)
		}
	}
}

func (p *Producer) Close() error {
	return p.producer.Close()
}
