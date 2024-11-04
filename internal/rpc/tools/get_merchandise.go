package tools

import (
	"encoding/json"
	"store/pkg/elasticsearch"
	"store/pkg/errors"
	"store/pkg/model"
)

func GetMerchandise(id string, es *elasticsearch.Elasticsearch) (model.Merchandise, error) {
	var m model.Merchandise

	data, err := es.GetDocumentByID(id)
	if err != nil {
		return m, errors.EsSearchError
	}

	if err := json.Unmarshal(data, &m); err != nil {
		return m, errors.JsonUnmarshalError
	}

	return m, nil
}
