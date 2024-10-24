package task_queue

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hibiken/asynq"
	"store/pkg/config"
	"store/pkg/email"
)

type Client struct {
	client *asynq.Client
	email  *email.Server
}

func NewClient(conf *config.GlobalConfig, e *email.Server) *Client {
	c := &Client{
		email: e,
	}
	c.init(conf)

	return c
}

func (c *Client) init(conf *config.GlobalConfig) {
	c.client = asynq.NewClient(asynq.RedisClientOpt{
		Addr: fmt.Sprintf("%s:%s", conf.Redis.Addr, conf.Redis.Port),
		DB:   conf.Redis.DB,
	})
}

func (c *Client) SendTask(taskType string, payload interface{}) error {
	switch taskType {
	case EMAIL:
		return c.sendEmail(payload)
	default:
		return errors.New("unknown task type")
	}
}

func (c *Client) sendEmail(payload interface{}) error {
	data, err := json.Marshal(&payload)
	if err != nil {
		return err
	}
	task := asynq.NewTask(EMAILDELIVERY, data)
	_, err = c.client.Enqueue(task, asynq.Queue(DEFAULT))
	if err != nil {
		return err
	}

	return nil
}
