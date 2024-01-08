package search_interval

import (
	"github.com/moaddib666/OpenSearch-Advanced-Proxy/internal/core/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimeDurationIntervalParser(t *testing.T) {
	parser := NewTimeDurationIntervalParser()

	testCases := []struct {
		name           string
		dateHistogram  models.DateHistogram
		expectedResult time.Duration
	}{
		{
			name: "Test with empty interval",
			dateHistogram: models.DateHistogram{
				Interval: "",
			},
			expectedResult: 24 * time.Hour, // default interval
		},
		{
			name: "Test with valid daily interval",
			dateHistogram: models.DateHistogram{
				Interval: "1d",
			},
			expectedResult: 24 * time.Hour,
		},
		{
			name: "Test with valid hourly interval",
			dateHistogram: models.DateHistogram{
				Interval: "1h",
			},
			expectedResult: 1 * time.Hour,
		},
		{
			name: "Test with valid minute interval",
			dateHistogram: models.DateHistogram{
				Interval: "30m",
			},
			expectedResult: 30 * time.Minute,
		},
		{
			name: "Test with invalid interval",
			dateHistogram: models.DateHistogram{
				Interval: "10x",
			},
			expectedResult: 24 * time.Hour, // default interval due to invalid input
		},
		{
			name: "Test with valid daily calendar interval",
			dateHistogram: models.DateHistogram{
				CalendarInterval: "1d",
			},
			expectedResult: 24 * time.Hour,
		},
		{
			name: "Test with valid weekly calendar interval",
			dateHistogram: models.DateHistogram{
				CalendarInterval: "1w",
			},
			expectedResult: 7 * 24 * time.Hour,
		},
		{
			name: "Test with valid monthly calendar interval",
			dateHistogram: models.DateHistogram{
				CalendarInterval: "1M",
			},
			expectedResult: 30 * 24 * time.Hour,
		},
		{
			name: "Test with valid yearly calendar interval",
			dateHistogram: models.DateHistogram{
				CalendarInterval: "1y",
			},
			expectedResult: 365 * 24 * time.Hour,
		},
		{
			name: "Test with invalid calendar interval",
			dateHistogram: models.DateHistogram{
				CalendarInterval: "10x",
			},
			expectedResult: 24 * time.Hour, // default interval due to invalid input
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var result time.Duration
			err := parser.Parse(&tc.dateHistogram, &result)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedResult, result)
		})
	}
}
