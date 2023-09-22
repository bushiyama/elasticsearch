package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

type Message struct {
	ID        string    `json:"id"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
}

func main() {
	//　公式サンプルだとCloudIDを設定しているが、"http"が使えないのでAddressesを設定しておく
	es, err := elasticsearch.NewTypedClient(elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
	})
	if err != nil {
		panic(err)
	}

	// 生のクエリも扱えるが...
	res, err := es.Search().
		Index("messages").
		Request(&search.Request{
			Query: &types.Query{
				Bool: &types.BoolQuery{
					Must: []types.Query{
						{
							MatchPhrase: map[string]types.MatchPhraseQuery{
								"body": {
									Query: "body",
								},
							},
						},
					},
				},
			},
		}).
		Do(context.TODO())
	if err != nil {
		panic(err)
	}
	m := &Message{}
	err = json.Unmarshal(res.Hits.Hits[0].Source_, m)
	if err != nil {
		panic(err)
	}
	fmt.Println(m)
}
