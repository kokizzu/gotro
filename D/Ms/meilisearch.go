package Ms

import (
	"github.com/francoispqt/onelog"
	"github.com/meilisearch/meilisearch-go"

	"github.com/kokizzu/gotro/A"
	"github.com/kokizzu/gotro/L"
)

type Meili struct {
	meilisearch.Client
	Log *onelog.Logger
}

func (m *Meili) Query(index, str string, searchReq *meilisearch.SearchRequest) (*meilisearch.SearchResponse, error) {
	return m.Client.Index(index).Search(str, searchReq)
}

func (m *Meili) UpsertOne(index string, rows A.MSX) (*meilisearch.Task, error) {
	task, err := m.Client.Index(index).AddDocuments(rows)
	if L.IsError(err, `Ms.Client.Index.AddDocuments`) {
		return nil, err
	}
	return m.Client.WaitForTask(task)
}

func (m *Meili) MigrateMeilisearch(space string, primaryKey string, rankingRules []string) error {
	task, err := m.CreateIndex(&meilisearch.IndexConfig{
		Uid:        space,
		PrimaryKey: primaryKey,
	})
	if L.IsError(err, `Ms.Client.CreateIndex`) {
		return err
	}
	_, err = m.Client.WaitForTask(task)
	if L.IsError(err, `Ms.Client.WaitForTask.CreateIndex`) {
		return err
	}
	task, err = m.Index(space).UpdateRankingRules(&rankingRules)
	if L.IsError(err, `Ms.Index.UpdateRankingRules`) {
		return err
	}
	_, err = m.Client.WaitForTask(task)
	L.IsError(err, `Ms.Client.WaitForTask.UpdateRankingRules`)
	return err
}
