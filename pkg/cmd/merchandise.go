package cmd

import (
	"store/internal/rpc/base"
	"store/internal/rpc/merchandise"
	"store/internal/rpc/middleware"
	"store/pkg/config"
	"store/pkg/constant"
	"store/pkg/errors"
)

func NewMerchandiseServer() {
	b, err := base.NewBase([]string{constant.MERCHANDISE, constant.MERCHANDISESTYLE})
	if err != nil {
		errors.HandleError(err)
		return
	}

	conf, err := config.NewGlobalConfig()
	if err != nil {
		errors.HandleError(err)
		return
	}

	auth := middleware.NewAuthMiddleware()

	server, err := merchandise.NewServer(conf, auth, b)
	if err != nil {
		errors.HandleError(err)
		return
	}

	if err := server.Run(); err != nil {
		errors.HandleError(err)
		return
	}
}
