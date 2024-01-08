package log_provider_test

import (
	"OpenSearchAdvancedProxy/internal/adapters/log_provider"
	"OpenSearchAdvancedProxy/internal/adapters/search"
	"OpenSearchAdvancedProxy/internal/adapters/search/search_interval"
	"OpenSearchAdvancedProxy/internal/core/models"
	"OpenSearchAdvancedProxy/internal/core/ports"
	"encoding/json"
	"os"
	"testing"
)

var filePath = "../../../examples/test.log"
var requestPath = "../../../examples/data/request/search.json"
var timestampField = "datetime"

func LoadSearchRequestFromFile(file string) (*models.SearchRequest, error) {
	var searchRequest *models.SearchRequest
	fh, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(fh, &searchRequest)
	if err != nil {
		return nil, err
	}
	return searchRequest, nil
}

// PerfTestFileProvider is a provider for performance testing
func BenchmarkFileProvider(b *testing.B) {
	request, err := LoadSearchRequestFromFile(requestPath)
	if err != nil {
		b.Fatal(err)
	}

	fileConstructor := func() ports.LogEntry {
		return &log_provider.JsonLogEntry{
			TimeStampField: timestampField,
		}
	}

	intervalParser := search_interval.NewTimeDurationIntervalParser()
	filterFactory := search.NewFilterFactory()

	// Resetting the timer is crucial to get accurate benchmark results
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		provider := log_provider.NewLogFileProvider(filePath, fileConstructor, nil, intervalParser, filterFactory)

		provider.BeginScan(request)
		for provider.Scan() {
			_ = provider.LogEntry()
		}
		for _, rqa := range request.Aggregations {
			_ = provider.AggregateResult(rqa)
		}
		provider.EndScan()
	}
}
