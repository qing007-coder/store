package main

import (
	"fmt"
	"store/internal/component/auth/server"
	"store/pkg/config"
	"store/pkg/mysql"
	"store/pkg/redis"
	"store/pkg/rules"
	server_ "store/pkg/sso/server"
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
	srv := server_.NewServer(rdb, db, conf, enforcer)
	auth := server.NewAuthApi(srv, db, enforcer)
	router := server.NewRouter(auth)
	if err := router.Run(); err != nil {
		fmt.Println(err)
	}
}
