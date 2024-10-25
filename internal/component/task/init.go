package task

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"store/pkg/config"
	"store/pkg/email"
	"store/pkg/model"
	"store/pkg/redis"
	"store/pkg/task_queue"
	"time"
)

type Server struct {
	srv  *asynq.Server
	e    *email.Server
	conf *config.GlobalConfig
	rdb  *redis.Client
}

func NewServer(conf *config.GlobalConfig, e *email.Server, rdb *redis.Client) *Server {
	srv := &Server{
		e:    e,
		conf: conf,
		rdb:  rdb,
	}
	srv.init(conf)

	return srv
}

func (s *Server) init(conf *config.GlobalConfig) {
	s.srv = asynq.NewServer(
		asynq.RedisClientOpt{Addr: fmt.Sprintf("%s:%s", conf.Redis.Addr, conf.Redis.Port)},
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				task_queue.CRITICAL: 6,
				task_queue.DEFAULT:  3,
				task_queue.LOW:      1,
			},
		},
	)
}

func (s *Server) Run() {
	mux := asynq.NewServeMux()
	mux.HandleFunc(task_queue.EMAILDELIVERY, s.emailDelivery)

	if err := s.srv.Run(mux); err != nil {
		fmt.Println("err:", err)
	}
}

func (s *Server) emailDelivery(ctx context.Context, task *asynq.Task) error {
	var e model.EmailTask
	if err := json.Unmarshal(task.Payload(), &e); err != nil {
		return err
	}
	code, err := s.e.SendVerificationCode(e.Email)
	if err != nil {
		fmt.Println("err:", err)
		return err
	}

	if err := s.rdb.Set(ctx, e.Email, code, time.Minute); err != nil {
		return err
	}
	if err := s.rdb.Set(ctx, e.Email+".send", 1, time.Minute*10); err != nil {
		return err
	}
	return nil
}
