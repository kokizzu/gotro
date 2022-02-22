package Es_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/hexops/autogold"
	"github.com/kokizzu/gotro/D/Es"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/M"
	"github.com/olivere/elastic/v7"
	"github.com/ory/dockertest/v3"
)

var adapter Es.Adapter

// go test -update -run TestElasticSearch

func TestMain(m *testing.M) {

	// connect to dockertest
	globalPool, err := dockertest.NewPool(``)
	L.PanicIf(err, `Could not connect to docker daemon`)

	resource, err := globalPool.RunWithOptions(&dockertest.RunOptions{
		Repository: `docker.elastic.co/elasticsearch/elasticsearch`,
		Tag:        `7.17.0`,
		Env: []string{
			"discovery.type=single-node",
			"xpack.security.enabled=false",
			"ES_JAVA_OPTS=-Xmx1g",
		},
	})
	L.PanicIf(err, "Could not start resource")

	const connStr = `http://127.0.0.1:%s`
	const esPort = `9200/tcp`

	// try connect to elasticsearch
	if err := globalPool.Retry(func() error {
		var err error

		s := fmt.Sprintf(connStr, resource.GetPort(esPort))

		opts := []elastic.ClientOptionFunc{
			elastic.SetURL(s),
			elastic.SetSniff(false),
		}
		opts = append(opts,
			elastic.SetErrorLog(log.New(os.Stderr, "ES-ERR ", log.LstdFlags)),
			elastic.SetInfoLog(log.New(os.Stdout, "ES-INFO ", log.LstdFlags)),
		)
		esClient, err := elastic.NewClient(opts...)
		if err != nil {
			return err
		}
		adapter = Es.Adapter{
			Client: esClient,
			Reconnect: func() *elastic.Client {
				esClient, err := elastic.NewClient(opts...)
				L.IsError(err, `elastic.NewClient %v`, opts)
				return esClient
			},
		}
		return nil
	}); err != nil {
		log.Printf("Cannot connect to spawned docker: %s\n", err)
		return
	}

	// run tests
	code := m.Run()
	defer func() {
		os.Exit(code)
	}()

	// cleanup dockertest
	if err := globalPool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
}

func TestElasticSearch(t *testing.T) {

	const indexName = `index1`
	ctx := context.Background()

	t.Run(`insertOne`, func(t *testing.T) {
		_, err := adapter.Index().Index(indexName).Id(`1`).
			BodyJson(M.SX{
				`type1`: `A`,
				`type2`: `X`,
			}).Refresh("true").Do(ctx)
		L.PanicIf(err, `adapter.Index.Do`)
	})

	t.Run(`insertMany`, func(t *testing.T) {
		bulkReq := adapter.Bulk()
		bulkReq.Add(elastic.NewBulkIndexRequest().Index(indexName).Id(`1`).Doc(M.SX{
			`type1`: `B`,
			`type2`: `Y`,
		}))
		bulkReq.Add(elastic.NewBulkIndexRequest().Index(indexName).Id(`2`).Doc(M.SX{
			`type1`: `C`,
			`type2`: `Z`,
		}))
		bulkReq.Add(elastic.NewBulkIndexRequest().Index(indexName).Id(`3`).Doc(M.SX{
			`type1`: `B`,
			`type2`: `Z`,
		}))
		bulkReq.Add(elastic.NewBulkIndexRequest().Index(indexName).Id(`4`).Doc(M.SX{
			`type1`: `C`,
			`type2`: `X`,
		}))
		_, err := bulkReq.Refresh("true").Do(ctx)
		L.PanicIf(err, `bulkReq.Do`)
	})

	t.Run(`querySimple`, func(t *testing.T) {
		res := []string{}
		adapter.QueryRaw(indexName, M.SX{
			`query`: M.SX{
				`bool`: M.SX{
					`must`: []interface{}{
						M.SX{`match`: M.SX{`type2`: `X`}},
						M.SX{`match`: M.SX{`type1`: `C`}},
					},
				},
			},
		}, func(id string, rawJson []byte) (exitEarly bool) {
			res = append(res, id)
			return false
		})
		fmt.Println(res)
		want := autogold.Want(`querySimple`, []string{"4"})
		want.Equal(t, res)
	})

	t.Run(`queryOr`, func(t *testing.T) {

		res := []string{}
		adapter.QueryRaw(indexName, M.SX{
			`query`: M.SX{
				`bool`: M.SX{
					`should`: []interface{}{
						M.SX{`match`: M.SX{`type2`: `Z`}},
						M.SX{`match`: M.SX{`type1`: `B`}},
						M.SX{`match`: M.SX{`type1`: `C`}},
					},
				},
			},
		}, func(id string, rawJson []byte) (exitEarly bool) {
			res = append(res, id)
			return false
		})
		fmt.Println(res)
		want := autogold.Want(`queryOr`, []string{"2", "3", "1", "4"})
		want.Equal(t, res)
	})
}
