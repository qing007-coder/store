package user

import (
	"store/internal/rpc/base"
	"store/pkg/sso/client"
)

type Server struct {
	*base.Base
	client *client.Client
}

func (s *Server) Register() {

}
