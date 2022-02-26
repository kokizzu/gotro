package Ms

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/kokizzu/gotro/L"
	"github.com/meilisearch/meilisearch-go"
	"github.com/ory/dockertest"
	"github.com/stretchr/testify/assert"
)

type Query struct {
	Index    string
	Document meilisearch.SearchRequest
}

type UpsertOne struct {
	Space      string
	row        interface{}
	primaryKey string
}

type MigrateMeilisearch struct {
	Space         string
	Id            string
	rangkingRules []string
}

var meili *Meili

func TestMain(m *testing.M) {
	// connect to dockertest
	globalPool, err := dockertest.NewPool(``)
	L.PanicIf(err, `Could not connect to docker daemon`)

	resource, err := globalPool.RunWithOptions(&dockertest.RunOptions{
		Repository: `getmeili/meilisearch`,
		Tag:        `v0.21.0`,
		Env: []string{
			"MEILI_NO_ANALYTICS=true",
			"MEILI_NO_SENTRY=true",
			"MEILI_MASTER_KEY=test_api_key",
		},
	})
	L.PanicIf(err, "Could not start resource")

	if err := globalPool.Retry(func() error {

		client := meilisearch.NewClient(meilisearch.Config{
			Host: fmt.Sprintf("http://127.0.0.1:%s", resource.GetPort("7700/tcp")),
		})

		return client.Health().Get()
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

var documents = []map[string]interface{}{
	{
		"id":   287947,
		"name": "apple",
	},
}

func TestMeiliFlow(t *testing.T) {
	t.Run(`insert`, func(t *testing.T) {
		in := &UpsertOne{
			Space:      "apple",
			row:        documents,
			primaryKey: "buah",
		}
		out, _ := meili.UpsertOne(in.Space, in.row, in.primaryKey)
		assert.Equal(t, out, "apple")
	})
	t.Run(`migrate`, func(t *testing.T) {
		in := &MigrateMeilisearch{
			Space: "apple",
			Id:    "fruit",
			rangkingRules: []string{
				"",
			},
		}
		out := meili.MigrateMeilisearch(in.Space, in.Id, in.rangkingRules)
		assert.Empty(t, out.Error)
	})
	// t.Run(`query`, func(t *testing.T) {
	// 	in := &Query{
	// 		Index:    "apple",
	// 		Document: meilisearch.SearchRequest{},
	// 	}
	// 	out, _ := meili.Query(in.Index, &in.Document)
	// 	assert.Equal(t, out, "apple")
	// })
}
