package Ms

import (
	"fmt"
	"github.com/kokizzu/gotro/M"
	"log"
	"os"
	"testing"

	"github.com/kokizzu/gotro/L"
	"github.com/meilisearch/meilisearch-go"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/assert"
)

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
			Host:   fmt.Sprintf("http://127.0.0.1:%s", resource.GetPort("7700/tcp")),
			APIKey: `test_api_key`,
		})

		err := client.Health().Get()
		if err != nil && err.Error() == `unaccepted status code found: 200 expected: [204], message from api: '', request: empty request (path "GET /health" with method "Health.Get")` {
			return nil
		}
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

func TestMeilisearch(t *testing.T) {
	const tableName = `fruits`
	t.Run(`migrate`, func(t *testing.T) {
		out := meili.MigrateMeilisearch(tableName, `id`, []string{
			`desc(created)`,
			`typo`,
			`words`,
			`proximity`,
			`attribute`,
			`wordsPosition`,
			`exactness`,
		})
		assert.Empty(t, out.Error)
	})
	t.Run(`insert`, func(t *testing.T) {
		in := &M.SX{
			`id`:   1,
			"name": "apple",
		}
		out, err := meili.UpsertOne(tableName, in, in.GetStr(`id`))
		assert.Nil(t, err)
		assert.Equal(t, out, "apple")
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
