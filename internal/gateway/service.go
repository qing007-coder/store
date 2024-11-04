package gateway

import (
	"fmt"
	grpcc "github.com/go-micro/plugins/v4/client/grpc"
	"github.com/go-micro/plugins/v4/registry/consul"
	grpcs "github.com/go-micro/plugins/v4/server/grpc"
	"go-micro.dev/v4"
	"go-micro.dev/v4/client"
	"go-micro.dev/v4/registry"
	"store/pkg/config"
)

type Service struct {
	srv micro.Service
}

func NewService(conf *config.GlobalConfig) *Service {
	srv := new(Service)
	srv.init(conf)

	return srv
}

func (s *Service) init(conf *config.GlobalConfig) {
	consulRegister := consul.NewRegistry(
		registry.Addrs(fmt.Sprintf("%s:%s", conf.Consul.Addr, conf.Consul.Port)),
	)
	s.srv = micro.NewService(
		micro.Server(grpcs.NewServer()), // 使用 gRPC server
		micro.Client(grpcc.NewClient()), // 使用 gRPC client
		micro.Registry(consulRegister),
	)
}

func (s *Service) Client() client.Client {
	return s.srv.Client()
}
