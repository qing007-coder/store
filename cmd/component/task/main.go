package main

import (
	"fmt"
	"store/internal/component/task"
	"store/pkg/config"
	"store/pkg/email"
	"store/pkg/redis"
)

func main() {
	conf, err := config.NewGlobalConfig()
	if err != nil {
		fmt.Println(err)
		return
	}
	rdb := redis.NewClient(conf)
	e := email.NewServer(conf)
	server := task.NewServer(conf, e, rdb)

	server.Run()
}
