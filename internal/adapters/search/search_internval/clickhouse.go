package search_internval

import (
	"OpenSearchAdvancedProxy/internal/core/models"
	"fmt"
)

type ClickHouseSearchIntervalParser struct {
	defaultInterval string
}

func (c *ClickHouseSearchIntervalParser) GetInterval(src *models.DateHistogram) string {
	if src.Interval != "" {
		// Bug in we go up to `ms` `ns` etc.
		lastChar := src.Interval[len(src.Interval)-1:]
		numberPart := src.Interval[:len(src.Interval)-1]
		switch lastChar {
		case "d":
			return fmt.Sprintf("%s DAY", numberPart)
		case "h":
			return fmt.Sprintf("%s HOUR", numberPart)
		case "m":
			return fmt.Sprintf("%s MINUTE", numberPart)
		case "s":
			return fmt.Sprintf("%s SECOND", numberPart)
		default:
			return c.defaultInterval
		}
	}
	if src.CalendarInterval != "" {

		lastChar := src.CalendarInterval[len(src.CalendarInterval)-1:]
		numberPart := src.CalendarInterval[:len(src.CalendarInterval)-1]

		// Convert Elasticsearch interval to Clickhouse format
		switch lastChar {
		case "d":
			return fmt.Sprintf("%s DAY", numberPart)
		case "w":
			return fmt.Sprintf("%s WEEK", numberPart)
		case "M":
			return fmt.Sprintf("%s MONTH", numberPart)
		case "y":
			return fmt.Sprintf("%s YEAR", numberPart)
		default:
			return c.defaultInterval
		}
	}
	return c.defaultInterval
}

func (c *ClickHouseSearchIntervalParser) Parse(src *models.DateHistogram, dest interface{}) error {
	switch dest.(type) {
	case *string:
		*(dest.(*string)) = c.GetInterval(src)
	default:
		return fmt.Errorf("not supported type: %T", dest)
	}

	return nil
}

func NewClickHouseSearchIntervalParser() *ClickHouseSearchIntervalParser {
	return &ClickHouseSearchIntervalParser{
		defaultInterval: "1 HOUR",
	}
}
