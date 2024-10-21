package main

import (
	"fmt"
	"store/internal/component/auth"
	"store/pkg/config"
	"store/pkg/mysql"
	"store/pkg/redis"
	"store/pkg/rules"
	"store/pkg/sso/server"
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
	authApi := auth.NewAuthApi(srv, db, enforcer)
	router := auth.NewRouter(authApi)
	if err := router.Run(); err != nil {
		fmt.Println(err)
	}
}
