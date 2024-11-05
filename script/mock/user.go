package mock

import (
	"encoding/json"
	"io"
	"os"
	"store/pkg/model"
)

type User struct {
	*Base
}

func NewUser(b *Base) *User {
	return &User{
		b,
	}
}

func (m *Merchandise) CreateUserMock() error {
	file, err := os.Open("./script/mock/data/user.json")
	if err != nil {
		return err
	}

	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	var users []model.User
	if err := json.Unmarshal(data, &users); err != nil {
		return err
	}

	for _, u := range users {
		if err := m.db.Create(&u).Error; err != nil {
			return err
		}
	}

	return nil
}
