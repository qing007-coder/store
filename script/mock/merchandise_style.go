package mock

import (
	"encoding/json"
	"io"
	"os"
	"store/pkg/constant"
	"store/pkg/model"
)

type MerchandiseStyle struct {
	*Base
}

func NewMerchandiseStyle(b *Base) *MerchandiseStyle {
	return &MerchandiseStyle{
		b,
	}
}

func (m *Merchandise) CreateMerchandiseStyleMock() error {
	file, err := os.Open("./script/mock/data/merchandise_style.json")
	if err != nil {
		return err
	}

	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	var ms []model.MerchandiseStyle
	if err := json.Unmarshal(data, &ms); err != nil {
		return err
	}

	for _, style := range ms {
		if err := m.es[constant.MERCHANDISESTYLE].CreateDocument(&style, style.ID); err != nil {
			return err
		}
	}

	return nil
}
