package log_storage

import (
	"context"
	"github.com/moaddib666/OpenSearch-Advanced-Proxy/internal/core/models"
	"github.com/moaddib666/OpenSearch-Advanced-Proxy/internal/core/ports"
	"time"
)

type GenericStorage struct {
	name      string
	fields    *models.Fields
	engine    ports.SearchEngine
	searchTTL time.Duration
}

// NewGenericStorage - create a new GenericStorage
func NewGenericStorage(name string, fields *models.Fields, engine ports.SearchEngine) *GenericStorage {
	return &GenericStorage{
		name:      name,
		fields:    fields,
		engine:    engine,
		searchTTL: 60 * time.Second,
	}
}

func (f *GenericStorage) Name() string {
	return f.name
}

func (f *GenericStorage) Fields() *models.Fields {
	return f.fields
}

func (f *GenericStorage) Search(r *models.SearchRequest) (*models.SearchResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), f.searchTTL)
	defer cancel()
	found, err := f.engine.ProcessSearch(ctx, r)
	if err != nil {
		return nil, err
	}
	return found, nil
}
