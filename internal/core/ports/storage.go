package ports

import "github.com/moaddib666/OpenSearch-Advanced-Proxy/internal/core/models"

// Storage is the interface that represents the external log storage.
type Storage interface {
	Name() string
	Fields() *models.Fields
	Search(r *models.SearchRequest) (*models.SearchResult, error)
}

// SearchFunc is a function that searches the storage.
type SearchFunc func(r *models.SearchRequest) (*models.SearchResult, error)

type StorageFactory interface {
	FromConfig(name string, config *models.SubConfig) (Storage, error)
}

type DistributedRequestsProcessor interface {
	MakeRequest(searchRequest *models.SearchRequest) (string, <-chan *models.SearchResult)
	AnswerRequest(id string, result *models.SearchResult)
	ResponseExpected(requestId string, answerRequires int)
}
