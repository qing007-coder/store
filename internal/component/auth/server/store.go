package server

import (
	"context"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/models"
	"gorm.io/gorm"
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
func (c *Store) GetByID(ctx context.Context, id string) (oauth2.ClientInfo, error) {

	return &models.Client{
		ID:     client.ID,
		UserID: client.UserID,
		Secret: client.Secret,
		Domain: client.Domain,
	}, nil
}
