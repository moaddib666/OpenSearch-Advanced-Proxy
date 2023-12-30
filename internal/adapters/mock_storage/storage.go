package mock_storage

import (
	"OpenSearchAdvancedProxy/internal/core/models"
	"time"
)

type MockStorage struct {
}

func (m *MockStorage) Name() string {
	return "mock"
}

func (m *MockStorage) Fields() *models.Fields {
	return &models.Fields{
		Fields: map[string]map[models.FieldType]*models.Field{
			"datetime": {
				models.DateType: models.NewField(models.DateType, true, true),
			},
			"message": {
				models.TextType: models.NewField(models.TextType, true, false),
			},
		},
	}
}

func (m *MockStorage) Search(r *models.SearchRequest) (*models.SearchResult, error) {
	datetimeNow := time.Now().UTC().Format(time.RFC3339)
	return &models.SearchResult{
		Took:     1,
		TimedOut: false,
		Shards: &models.Shards{
			Total:      1,
			Successful: 1,
			Skipped:    0,
			Failed:     0,
		},
		Hits: &models.Hits{
			Total: &models.TotalValue{
				Value: 1,
			},
			Hits: []*models.Hit{
				{
					Index: m.Name(),
					Source: map[string]interface{}{
						"datetime": datetimeNow,
						"message":  "hello world",
					},
				},
			},
		},
	}, nil
}

// NewMockStorage creates a new MockStorage struct
func NewMockStorage() *MockStorage {
	return &MockStorage{}
}
