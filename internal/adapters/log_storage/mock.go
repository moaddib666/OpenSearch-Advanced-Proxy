package log_storage

import (
	"OpenSearchAdvancedProxy/internal/adapters/log_provider"
	"OpenSearchAdvancedProxy/internal/adapters/search"
	"OpenSearchAdvancedProxy/internal/core/models"
	"OpenSearchAdvancedProxy/internal/core/ports"
	"encoding/json"
	log "github.com/sirupsen/logrus"
)

type MockStorage struct {
	engine ports.SearchEngine
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
	log.Debugf("Searching storage: `%s`", m.Name())
	jsonRequest, _ := json.Marshal(r)
	log.Debugf("Search request: %s", string(jsonRequest))

	found, err := m.engine.ProcessSearch(r)
	if err != nil {
		return nil, err
	}
	count := len(found)
	hits := make([]*models.Hit, count)
	for i, entry := range found {
		hits[i] = &models.Hit{
			Index:  m.Name(),
			Source: entry.Map(),
		}
	}

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
				Value: count,
			},
			Hits: hits,
		},
	}, nil
}

// NewMockStorage creates a new MockStorage struct
func NewMockStorage() *MockStorage {
	provider := log_provider.NewLogFileProvider("tmp/test.log", log_provider.JsonLogEntryConstructor)
	return &MockStorage{
		engine: search.NewLogSearchEngine(provider),
	}
}
