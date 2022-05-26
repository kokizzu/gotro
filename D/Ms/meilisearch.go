package Ms

import (
	"github.com/francoispqt/onelog"
	"github.com/kokizzu/gotro/A"
	"github.com/kokizzu/gotro/L"
	"github.com/meilisearch/meilisearch-go"
)

type Meili struct {
	meilisearch.ClientInterface
	Log *onelog.Logger
}

func (m *Meili) Query(space string, searchReq *meilisearch.SearchRequest) (*meilisearch.SearchResponse, error) {
	return m.Query(space, searchReq)
}

func (m *Meili) UpsertOne(space string, row interface{}, primaryKey string) (interface{}, error) {
	return m.UpsertOne(space, A.X{row}, primaryKey)
}

func (m *Meili) MigrateMeilisearch(space string, id string, rankingRules []string) error {
	_, err := m.CreateIndex(&meilisearch.IndexConfig{
		Uid:        space,
		PrimaryKey: id,
	})
	if err != nil {
		merr, ok := err.(*meilisearch.Error)
		if !ok || merr.Error() != `Index `+space+` already exists` {
			L.IsError(err, `failed create index`)
			return err
		}
	}
	_, err = m.Index(space).UpdateRankingRules(&rankingRules)
	L.IsError(err, `failed update rankingRules`)
	return err
}
