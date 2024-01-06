package search_internval

import (
	"OpenSearchAdvancedProxy/internal/core/models"
	"fmt"
	"time"
)

type TimeDurationIntervalParser struct {
	defaultInterval time.Duration
}

func (c *TimeDurationIntervalParser) GetInterval(src *models.DateHistogram) time.Duration {
	if src.Interval != "" {
		// Bug in we go up to `ms` `ns` etc.
		lastChar := src.Interval[len(src.Interval)-1:]
		numberPart := src.Interval[:len(src.Interval)-1]
		parseIntNumberPart, err := time.ParseDuration(numberPart)
		if err != nil {
			return c.defaultInterval
		}
		switch lastChar {
		case "d":
			return parseIntNumberPart * 24 * time.Hour
		case "h":
			return parseIntNumberPart * time.Hour
		case "m":
			return parseIntNumberPart * time.Minute
		case "s":
			return parseIntNumberPart * time.Second
		default:
			return c.defaultInterval
		}
	}
	if src.CalendarInterval != "" {

		lastChar := src.CalendarInterval[len(src.CalendarInterval)-1:]
		numberPart := src.CalendarInterval[:len(src.CalendarInterval)-1]
		parseNumberPart, err := time.ParseDuration(numberPart)
		if err != nil {
			return c.defaultInterval
		}

		// Convert Elasticsearch interval to Clickhouse format
		switch lastChar {
		case "d":
			return parseNumberPart * 24 * time.Hour
		case "w":
			return parseNumberPart * 7 * 24 * time.Hour
		case "M":
			return parseNumberPart * 30 * 24 * time.Hour
		case "y":
			return parseNumberPart * 365 * 24 * time.Hour
		default:
			return c.defaultInterval
		}
	}
	return c.defaultInterval
}

func (c *TimeDurationIntervalParser) Parse(src *models.DateHistogram, dest interface{}) error {
	switch dest.(type) {
	case *time.Duration:
		*(dest.(*time.Duration)) = c.GetInterval(src)
	default:
		return fmt.Errorf("not supported type: %T", dest)
	}

	return nil
}

func NewTimeDurationIntervalParser() *TimeDurationIntervalParser {
	return &TimeDurationIntervalParser{
		defaultInterval: 24 * time.Hour,
	}
}
