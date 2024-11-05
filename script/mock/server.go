package mock

import (
	"store/pkg/constant"
	"store/pkg/model"
	"store/script/elasticsearch"
)

type Server struct {
	*User
	*MerchandiseStyle
	*Merchandise
	*Base
}

func NewServer() (*Server, error) {
	b, err := NewBase([]string{constant.MERCHANDISE, constant.MERCHANDISESTYLE})
	if err != nil {
		return nil, err
	}

	return &Server{
		NewUser(b),
		NewMerchandiseStyle(b),
		NewMerchandise(b),
		b,
	}, nil
}

func (s *Server) Run() error {
	var err error

	err = s.CreateUserMock()
	err = s.CreateMerchandiseStyleMock()
	err = s.CreateMerchandiseMock()

	if err != nil {
		return s.Clear()
	}

	return nil
}

func (s *Server) Clear() error {
	if err := elasticsearch.DeleteIndex(s.es); err != nil {
		return err
	}

	if err := elasticsearch.CreateIndex(s.es); err != nil {
		return err
	}

	return s.db.Unscoped().Where("1 = 1").Delete(&model.User{}).Error
}
