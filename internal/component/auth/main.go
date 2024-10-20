package main

import (
	"store/internal/component/auth/server"
	"store/pkg/redis"
)

func main() {
	redis.NewClient()
	server.NewServer()
}
