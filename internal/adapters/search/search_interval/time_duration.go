package search_interval

import (
	"fmt"
	"github.com/moaddib666/OpenSearch-Advanced-Proxy/internal/core/models"
	"strconv"
	"time"
)

type TimeDurationIntervalParser struct {
	defaultInterval time.Duration
}

func (c *TimeDurationIntervalParser) getInterval(intervalString string) time.Duration {
	lastChar := intervalString[len(intervalString)-1:]
	numberPart := intervalString[:len(intervalString)-1]
	parseIntNumberPart, err := strconv.ParseInt(numberPart, 10, 64)
	if err != nil {
		return c.defaultInterval
	}
	duration := time.Duration(parseIntNumberPart)

	switch lastChar {
	// Parse Elasticsearch interval starting from the smallest unit to the largest
	case "s":
		return duration * time.Second
	case "m":
		return duration * time.Minute
	case "h":
		return duration * time.Hour
	case "d":
		return duration * 24 * time.Hour
	case "w":
		return duration * 7 * 24 * time.Hour
	case "M":
		return duration * 30 * 24 * time.Hour
	case "y":
		return duration * 365 * 24 * time.Hour
	default:
		return c.defaultInterval

	}
}
func (c *TimeDurationIntervalParser) GetInterval(src *models.DateHistogram) time.Duration {
	if src.Interval != "" {
		return c.getInterval(src.Interval)
	}
	if src.CalendarInterval != "" {
		return c.getInterval(src.CalendarInterval)
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
