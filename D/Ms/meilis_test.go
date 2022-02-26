package Ms

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/meilisearch/meilisearch-go"
	"github.com/stretchr/testify/assert"

	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/M"
	"github.com/ory/dockertest/v3"
)

func TestMain(m *testing.M) {
	// connect to dockertest
	globalPool, err := dockertest.NewPool(``)
	L.PanicIf(err, `Could not connect to docker daemon`)

	resource, err := globalPool.RunWithOptions(&dockertest.RunOptions{
		Repository: `getmeili/meilisearch`,
		Tag:        `v0.25.0`,
		Env: []string{
			"MEILI_NO_ANALYTICS=true",
			"MEILI_NO_SENTRY=true",
			"MEILI_MASTER_KEY=test_api_key",
		},
	})
	L.PanicIf(err, "Could not start resource")

	if err := globalPool.Retry(func() error {

		client := meilisearch.NewClient(meilisearch.ClientConfig{
			Host: fmt.Sprintf("http://127.0.0.1:%s", resource.GetPort("7700/tcp")),
		})

		_, err := client.Health()
		return err
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

var m *Meilis

type Create struct {
	Index      string
	Documents  map[string]interface{}
	PrimaryKey string
}

type Search struct {
	Index     string
	Search    string
	MeilisReq meilisearch.SearchRequest
}

func TestMeilis(t *testing.T) {
	t.Run(`insert`, func(t *testing.T) {
		in := &Create{
			Index: "title",
			Documents: M.SX{
				"title":        "Kung Fu Panda",
				"genre":        "Children's Animation",
				"release-year": 2008,
				"cast": []map[string]interface{}{
					{"Jack Black": "Po"},
					{"Jackie Chan": "Monkey"},
				},
			},
			PrimaryKey: "release-year",
		}
		out, _ := m.Create(in.Index, in.Documents, in.PrimaryKey)
		assert.Empty(t, out.Error)
	})
	t.Run(`search`, func(t *testing.T) {
		in := &Search{
			Index:  "movies",
			Search: "Kung Fu Panda",
			MeilisReq: meilisearch.SearchRequest{
				Limit: 1,
			},
		}
		out, _ := m.Search(in.Index, in.Search, in.MeilisReq)
		assert.Empty(t, out.Limit)
	})
}
