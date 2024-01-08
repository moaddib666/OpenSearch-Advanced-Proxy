package ports

import "github.com/moaddib666/OpenSearch-Advanced-Proxy/internal/core/models"

type SearchInternalParser interface {
	Parse(src *models.DateHistogram, dest interface{}) error
}
