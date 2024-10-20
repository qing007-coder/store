package main

import (
	"fmt"
	"store/internal/component/auth/server"
	"store/pkg/config"
	"store/pkg/mysql"
	"store/pkg/redis"
	"store/pkg/rules"
)

func main() {
	conf, err := config.NewGlobalConfig()
	if err != nil {
		fmt.Println(err)
		return
	}
	rdb := redis.NewClient(conf)
	db, err := mysql.NewClient(conf)
	if err != nil {
		fmt.Println(err)
		return
	}
	enforcer := rules.NewEnforcer(db)
	srv := server.NewServer(rdb, db, conf, enforcer)
	auth := server.NewAuthApi(srv, db, enforcer)
	router := server.NewRouter(auth)
	if err := router.Run(); err != nil {
		fmt.Println(err)
	}
}
