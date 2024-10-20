package server

import (
	"context"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/models"
	"gorm.io/gorm"
	"store/pkg/model"
)

type Store struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) *Store {
	return &Store{
		db: db,
	}
}

// GetByID 实现oauth2.0的接口
func (s *Store) GetByID(ctx context.Context, id string) (oauth2.ClientInfo, error) {
	var cli model.Client
	if err := s.db.Where("id = ?", id).First(&cli).Error; err != nil {
		return nil, err
	}

	return &models.Client{
		ID:     cli.ID,
		UserID: cli.UserID,
		Secret: cli.Secret,
		Domain: cli.Domain,
	}, nil
}
