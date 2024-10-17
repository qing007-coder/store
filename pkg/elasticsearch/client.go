package elasticsearch

import (
	"context"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/update"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

type Elasticsearch struct {
	Ctx    context.Context
	Client *elasticsearch.TypedClient
	Index  string
}

func NewClient(ctx context.Context, address string, index string) (*Elasticsearch, error) {
	client := new(Elasticsearch)
	err := client.init(ctx, address, index)
	return client, err
}

func (es *Elasticsearch) init(ctx context.Context, address string, index string) error {
	conf := elasticsearch.Config{
		Addresses: []string{
			address,
		},
	}

	client, err := elasticsearch.NewTypedClient(conf)
	if err != nil {
		return err
	}

	es.Client = client
	es.Index = index
	es.Ctx = ctx
	return nil
}

func (es *Elasticsearch) CreateIndex() error {
	_, err := es.Client.Indices.Create(es.Index).Do(es.Ctx)
	return err
}

func (es *Elasticsearch) CreateDocument(document interface{}, id string) error {
	_, err := es.Client.Index(es.Index).
		Id(id).
		Document(document).
		Do(es.Ctx)

	return err
}

func (es *Elasticsearch) GetDocumentByID(id string) ([]byte, error) {
	resp, err := es.Client.Get(es.Index, id).Do(es.Ctx)
	return resp.Source_, err
}

func (es *Elasticsearch) Search(mustQueries []types.Query, shouldQueries []types.Query, sort []types.SortCombinations, from, size int) (*search.Response, error) {
	return es.Client.Search().
		Index(es.Index).
		Request(&search.Request{
			Query: &types.Query{
				Bool: &types.BoolQuery{
					Must:   mustQueries,   // 每个都要匹配
					Should: shouldQueries, // 有一个匹配就可以
				},
			},
			Sort: sort,
			Size: &size,
			From: &from,
		}).Do(es.Ctx)
}

func (es *Elasticsearch) GetList(size, from int) (*search.Response, error) {
	return es.Client.Search().
		Index(es.Index).
		Request(&search.Request{
			Query: &types.Query{
				MatchAll: &types.MatchAllQuery{},
			},

			Size: &size,
			From: &from,
		}).Do(es.Ctx)
}

func (es *Elasticsearch) Update(id string, query map[string]interface{}) error {
	data, _ := json.Marshal(&query)
	_, err := es.Client.Update(es.Index, id).
		Request(&update.Request{
			Doc: data,
		}).Do(context.TODO())

	return err
}

// UpdateDownloads 自增
func (es *Elasticsearch) UpdateDownloads(id string) error {
	script := "ctx._source.download_total += 1"
	_, err := es.Client.Update(es.Index, id).Script(&types.Script{
		Source: &script,
	}).Do(es.Ctx)

	return err
}

func (es *Elasticsearch) DeleteDocument(id string) error {
	_, err := es.Client.Delete(es.Index, id).Do(es.Ctx)
	return err
}
