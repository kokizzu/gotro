package Ms

import "github.com/meilisearch/meilisearch-go"

type Meilis struct {
	meilisearch.Client
}

func (m *Meilis) Create(index string, documents map[string]interface{}, primaryKey string) (*meilisearch.Task, error) {
	return m.Index(index).AddDocuments(documents, primaryKey)
}

func (m *Meilis) Search(index, search string, request meilisearch.SearchRequest) (*meilisearch.SearchResponse, error) {
	return m.Index(index).Search(search, &request)
}
