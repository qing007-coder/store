package main

import (
	"store/internal/gateway"
	"store/pkg/config"
	"store/pkg/errors"
)

func main() {
	conf, err := config.NewGlobalConfig()
	if err != nil {
		errors.HandleError(err)
		return
	}

	router := gateway.NewRouter(conf)
	if err := router.Run(); err != nil {
		errors.HandleError(err)
		return
	}
}
