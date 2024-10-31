package elasticsearch

import "store/pkg/elasticsearch"

func CreateIndex(es map[string]*elasticsearch.Elasticsearch) error {
	for _, e := range es {
		return e.CreateIndex()
	}

	return nil
}