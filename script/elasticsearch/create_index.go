package elasticsearch

import (
	"store/pkg/elasticsearch"
)

func CreateIndex(es map[string]*elasticsearch.Elasticsearch) error {
	for _, e := range es {
		if err := e.CreateIndex(); err != nil {
			return err
		}
	}

	return nil
}

func DeleteIndex(es map[string]*elasticsearch.Elasticsearch) error {
	for _, e := range es {
		if err := e.DeleteIndex(); err != nil {
			return err
		}
	}

	return nil
}
