package user

import (
	"fmt"
	grpcc "github.com/go-micro/plugins/v4/client/grpc"
	"github.com/go-micro/plugins/v4/registry/consul"
	grpcs "github.com/go-micro/plugins/v4/server/grpc"
	"go-micro.dev/v4"
	"go-micro.dev/v4/registry"
	"store/internal/proto/user"
	"store/internal/rpc/base"
	"store/internal/rpc/middleware"
	"store/pkg/config"
)

type Server struct {
	m    *middleware.AuthMiddleware
	conf *config.GlobalConfig
	srv  micro.Service
	user *User
}

func NewServer(conf *config.GlobalConfig, m *middleware.AuthMiddleware, b *base.Base) (*Server, error) {
	s := &Server{
		conf: conf,
		m:    m,
		user: NewUser(b),
	}

	if err := s.init(conf); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Server) init(conf *config.GlobalConfig) error {
	c := consul.NewRegistry(
		registry.Addrs(fmt.Sprintf("%s:%s", conf.Consul.Addr, conf.Consul.Port)),
	)

	s.srv = micro.NewService(
		micro.Server(grpcs.NewServer()),
		micro.Client(grpcc.NewClient()), // 使用 gRPC client
		micro.Name("user"),
		micro.Version("latest"),
		micro.Registry(c),           // 必须放底下哎，不然注册中心的优先级会变的
		micro.WrapHandler(s.m.Auth), // 这个也是 顺序不能变
		micro.Address(":54182"),
	)

	s.srv.Init()

	return user.RegisterUserServiceHandler(s.srv.Server(), s.user)
}

func (s *Server) Run() error {
	return s.srv.Run()
}
