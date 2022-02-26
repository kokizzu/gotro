package Ms

import (
	"github.com/francoispqt/onelog"
	"github.com/meilisearch/meilisearch-go"
)

type Meili struct {
	meilisearch.ClientInterface
	Log *onelog.Logger
}

func (m *Meili) Query(space string, searchReq *meilisearch.SearchRequest) (*meilisearch.SearchResponse, error) {
	// return m.Search(space).Search(*searchReq)
	return nil, nil
}

func (m *Meili) UpsertOne(space string, row interface{}, primaryKey string) (interface{}, error) {
	// return m.Documents(space).AddOrReplaceWithPrimaryKey(A.X{row}, primaryKey)
	return nil, nil
}

func (m *Meili) MigrateMeilisearch(space string, id string, rankingRules []string) error {
	// _, err := m.Indexes().Create(meilisearch.CreateIndexRequest{
	// 	UID:        space,
	// 	PrimaryKey: id,
	// })
	// if err != nil {
	// 	merr, ok := err.(*meilisearch.Error)
	// 	if !ok || merr.MeilisearchMessage != `Index `+space+` already exists` {
	// 		L.IsError(err, `failed create index`)
	// 		return err
	// 	}
	// }
	// _, err = m.Settings(space).UpdateRankingRules(rankingRules)
	// L.IsError(err, `failed update rankingRules`)
	return nil
}
