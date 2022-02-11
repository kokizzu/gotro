package Es

import (
	"context"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/M"
	"github.com/olivere/elastic/v7"
	"log"
	"os"
)

type Adapter struct {
	*elastic.Client
	Reconnect func() *elastic.Client
}

func Connect1(url string, debug bool) *elastic.Client {
	opts := []elastic.ClientOptionFunc{
		elastic.SetURL(url),
		elastic.SetSniff(false),
	}
	if debug {
		opts = append(opts,
			elastic.SetErrorLog(log.New(os.Stderr, "ES-ERR ", log.LstdFlags)),
			elastic.SetInfoLog(log.New(os.Stdout, "ES-INFO ", log.LstdFlags)),
		)
	}
	esClient, err := elastic.NewClient(opts...)
	L.PanicIf(err, `elastic.NewClient`)
	_, _, err = esClient.Ping(url).Do(context.Background())
	L.PanicIf(err, `elastic.Ping`)
	return esClient
}

/* example parameter:

{
	"query": {
		"bool": {
			"must": [
				{"match": { "cards.token": TOKEN_HERE }}
			]
		}
	}
}

became:

M.SX{
	`query`: M.SX{
		`bool`: M.SX{
			`must`: []interface{}{
				M.SX{`match`: M.SX{`cards.token`: Token}},
			},
		},
	},
}

and inside you must deserialize it yourself:

_ = a.QueryRaw(index, query, func(id string, rawJson []byte) (exitEarly bool) {
	row := MyStruct{}
	err := json.Unmarshal(rawJson, &v)
	if L.IsError(err, `json.Unmarshal`) {
		return true
	}
	... // do something with row
	return
})
*/

func (a *Adapter) QueryRaw(index string, query M.SX, eachRowFunc func(id string, rawJson []byte) (exitEarly bool)) (res *elastic.SearchResult) {
	var err error
	res, err = a.Search(index).Source(query).Do(context.Background())
	if L.IsError(err, `Es.QueryRaw.Search`) {
		return
	}
	for _, hit := range res.Hits.Hits {
		if eachRowFunc(hit.Id, hit.Source) {
			return
		}
	}
	return
}

// possible generic 1.18 implementation:
// but go 1.18 method doesn't allow this:
// https://github.com/golang/go/issues/49085

//func (a *Adapter) Query[T any](index string, query M.SX, eachRowFunc func(id string, row *T) (exitEarly bool)) (res *elastic.SearchResult) {
//	var err error
//	res, err = a.Search(index).Source(query).Do(context.Background())
//	if L.IsError(err, `Es.Search`) {
//		return
//	}
//	for _, hit := range res.Hits.Hits {
//		var v T
//		err := json.Unmarshal(hit.Source, &v)
//		if L.IsError(err, `json.Unmarshal`) {
//			return
//		}
//		if eachRowFunc(hit.Id, v) {
//			return
//		}
//	}
//	return
//}

// so we might end up with for 1.18:

//func Query[T any](a *Adapter, index string, query M.SX, eachRowFunc func(id string, row *T) (exitEarly bool)) (res *elastic.SearchResult) {
//	var err error
//	res, err = a.Search(index).Source(query).Do(context.Background())
//	if L.IsError(err, `Es.Search`) {
//		return
//	}
//	for _, hit := range res.Hits.Hits {
//		var v T
//		err := json.Unmarshal(hit.Source, &v)
//		if L.IsError(err, `json.Unmarshal`) {
//			return
//		}
//		if eachRowFunc(hit.Id, v) {
//			return
//		}
//	}
//	return
//}

// or use template in the struct but not using it as data member like this:

// type Adapter[T any] struct { ... }

//func (a *Adapter[T]) Query(index string, query M.SX, eachRowFunc func(id string, row *T) (exitEarly bool)) (res *elastic.SearchResult) {
//	var err error
//	res, err = a.Search(index).Source(query).Do(context.Background())
//	if L.IsError(err, `Es.Search`) {
//		return
//	}
//	for _, hit := range res.Hits.Hits {
//		var v T
//		err := json.Unmarshal(hit.Source, &v)
//		if L.IsError(err, `json.Unmarshal`) {
//			return
//		}
//		if eachRowFunc(hit.Id, v) {
//			return
//		}
//	}
//	return
//}
