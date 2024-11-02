package cmd

import (
	"store/internal/rpc/base"
	"store/internal/rpc/middleware"
	"store/internal/rpc/user"
	"store/pkg/config"
	"store/pkg/constant"
	"store/pkg/errors"
)

func NewUserServer() {
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

	server, err := user.NewServer(conf, auth, b)
	if err != nil {
		errors.HandleError(err)
		return
	}

	if err := server.Run(); err != nil {
		errors.HandleError(err)
		return
	}
}
