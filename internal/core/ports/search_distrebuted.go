package ports

import "github.com/moaddib666/OpenSearch-Advanced-Proxy/internal/core/models"

type DistributedSearchProtocol interface {
	MarshallSearchRequest(request *models.DistributedSearchRequest) []byte
	UnmarshallSearchRequest(request []byte) (*models.DistributedSearchRequest, error)
	MarshallSearchResult(result *models.DistributedSearchResult) []byte
	UnmarshallSearchResult(result []byte) (*models.DistributedSearchResult, error)
}
