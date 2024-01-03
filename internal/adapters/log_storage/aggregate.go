package log_storage

import (
	"OpenSearchAdvancedProxy/internal/core/models"
	"OpenSearchAdvancedProxy/internal/core/ports"
	log "github.com/sirupsen/logrus"
)

type AggregateStorage struct {
	name              string
	storages          []ports.Storage
	fields            *models.Fields
	aggregatorFactory ports.SearchAggregatorFactory
}

// NewAggregateStorage - create a new AggregateStorage
func NewAggregateStorage(name string, storages []ports.Storage, fields *models.Fields, aggregatorFactory ports.SearchAggregatorFactory) *AggregateStorage {
	return &AggregateStorage{
		name:              name,
		storages:          storages,
		fields:            fields,
		aggregatorFactory: aggregatorFactory,
	}
}

func (a *AggregateStorage) Name() string {
	return a.name
}

func (a *AggregateStorage) Fields() *models.Fields {
	return a.fields
}

func (a *AggregateStorage) Search(r *models.SearchRequest) (*models.SearchResult, error) {
	if len(a.storages) == 0 {
		return nil, models.ErrNoStorages
	}
	found := make(chan *models.SearchResult, len(a.storages))
	aggregator := a.aggregatorFactory.CreateAggregator(r)
	for _, storage := range a.storages {
		if storage == nil {
			log.Fatalf("Found nil storage in %s", a.name)
		}
		// TODO: In case of performance issues, we can use a pool of goroutines here
		store := storage
		go func() {
			log.Debugf("Searching in storage %s, %s", store.Name(), r)
			result, err := store.Search(r)
			if err != nil {
				return
			}
			found <- result
		}()
	}
	for i := 0; i < len(a.storages); i++ {
		result := <-found
		aggregator.AddResult(result)
	}
	return aggregator.GetResult(), nil
}
