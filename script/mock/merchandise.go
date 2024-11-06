package mock

import (
	"encoding/json"
	"io"
	"os"
	"store/pkg/constant/store"
	"store/pkg/model"
)

type Merchandise struct {
	*Base
}

func NewMerchandise(b *Base) *Merchandise {
	return &Merchandise{
		b,
	}
}

func (m *Merchandise) CreateMerchandiseMock() error {
	file, err := os.Open("./script/mock/data/merchandise.json")
	if err != nil {
		return err
	}

	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	var ms []model.Merchandise
	if err := json.Unmarshal(data, &ms); err != nil {
		return err
	}

	for _, merchandise := range ms {
		if err := m.es[store.MERCHANDISE].CreateDocument(&merchandise, merchandise.ID); err != nil {
			return err
		}
	}

	return nil
}
