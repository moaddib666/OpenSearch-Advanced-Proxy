package search

import (
	"OpenSearchAdvancedProxy/internal/core/models"
	"encoding/json"
)

type DistributedJsonSearchProtocol struct {
}

func (d *DistributedJsonSearchProtocol) MarshallSearchRequest(request *models.DistributedSearchRequest) []byte {
	raw, _ := json.Marshal(request)
	return raw
}

func (d *DistributedJsonSearchProtocol) UnmarshallSearchRequest(request []byte) (*models.DistributedSearchRequest, error) {
	result := &models.DistributedSearchRequest{}
	err := json.Unmarshal(request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (d *DistributedJsonSearchProtocol) MarshallSearchResult(result *models.DistributedSearchResult) []byte {
	raw, _ := json.Marshal(result)
	return raw
}

func (d *DistributedJsonSearchProtocol) UnmarshallSearchResult(result []byte) (*models.DistributedSearchResult, error) {
	searchResult := &models.DistributedSearchResult{}
	err := json.Unmarshal(result, searchResult)
	if err != nil {
		return nil, err
	}
	return searchResult, nil
}

func NewDistributedJsonSearchProtocol() *DistributedJsonSearchProtocol {
	return &DistributedJsonSearchProtocol{}
}
