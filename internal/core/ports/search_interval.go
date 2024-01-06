package ports

import "OpenSearchAdvancedProxy/internal/core/models"

type SearchInternalParser interface {
	Parse(src *models.DateHistogram, dest interface{}) error
}
