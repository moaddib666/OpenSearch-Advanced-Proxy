package search

import (
	"OpenSearchAdvancedProxy/internal/core/models"
	"OpenSearchAdvancedProxy/internal/core/ports"
	"fmt"
	log "github.com/sirupsen/logrus"
	"sort"
	"sync"
	"time"
)

// TimeIntervalAggregator is an implementation of Aggregator that aggregates data into time intervals.
type TimeIntervalAggregator struct {
	FixedInterval string
	Field         string
	getHitTime    ports.HitTimeParser
}

// NewTimeIntervalAggregator creates a new TimeIntervalAggregator
func NewTimeIntervalAggregator(fixedInterval string, field string, getHitTime ports.HitTimeParser) *TimeIntervalAggregator {
	return &TimeIntervalAggregator{
		FixedInterval: fixedInterval,
		Field:         field,
		getHitTime:    getHitTime,
	}
}

// AutoDetectInterval if interval is not set, auto detect interval
func (t *TimeIntervalAggregator) AutoDetectInterval(hits []*models.Hit, rule SortRule) (time.Duration, error) {

	if len(hits) < 2 {
		return 0, fmt.Errorf("not enough data to auto-detect interval")
	}

	index1 := 0
	index2 := len(hits) - 1
	if rule == Desc {
		index1 = len(hits) - 1
		index2 = 0
	}
	firstHitTime, err := t.getHitTime(hits[index1], t.Field)
	if err != nil {
		log.Warnf("Invalid hit time: %v", err)
		return 0, err
	}
	lastHitTime, err := t.getHitTime(hits[index2], t.Field)
	if err != nil {
		log.Warnf("Invalid hit time: %v", err)
		return 0, err
	}

	// Calculate interval
	totalDuration := lastHitTime.Sub(firstHitTime)
	var interval time.Duration
	// From OpenSearch docs: https://opensearch.org/docs/latest/field-types/supported-field-types/date/
	// also https://www.elastic.co/guide/en/elasticsearch/reference/current/search-aggregations-bucket-datehistogram-aggregation.html
	switch {
	case totalDuration <= 2*time.Hour:
		interval = 1 * time.Minute
	case totalDuration <= 24*time.Hour:
		interval = 10 * time.Minute
	case totalDuration <= 7*24*time.Hour:
		interval = 30 * time.Minute
	case totalDuration <= 31*24*time.Hour:
		interval = 1 * time.Hour
	case totalDuration <= 3*31*24*time.Hour:
		interval = 12 * time.Hour
	case totalDuration <= 6*31*24*time.Hour:
		interval = 1 * 24 * time.Hour
	default:
		interval = 7 * 24 * time.Hour
	}

	return interval, nil
}

type SortRule string

const Asc SortRule = "asc"
const Desc SortRule = "desc"

// DetectSortRule detects the sort rule from the hits
func (t *TimeIntervalAggregator) DetectSortRule(hits []*models.Hit) SortRule {
	// if second hit is after first hit, then asc
	if len(hits) < 2 {
		return Asc
	}
	firstHitTime, err := t.getHitTime(hits[0], t.Field)
	if err != nil {
		log.Warnf("Invalid hit time: %v", err)
		return Asc
	}
	secondHitTime, err := t.getHitTime(hits[1], t.Field)
	if err != nil {
		log.Warnf("Invalid hit time: %v", err)
		return Asc
	}
	if secondHitTime.After(firstHitTime) {
		return Asc
	}
	return Desc
}

func (t *TimeIntervalAggregator) Aggregate(hits []*models.Hit) *models.AggregationResult {
	// Parse interval
	sortRule := t.DetectSortRule(hits) // Detect sort rule (asc or desc)
	intervalDuration, err := time.ParseDuration(t.FixedInterval)
	if err != nil {
		log.Debugf("Invalid interval: %s", t.FixedInterval)
		intervalDuration, err = t.AutoDetectInterval(hits, sortRule)
		if err != nil {
			log.Warnf("Invalid interval: %s", t.FixedInterval)
			return &models.AggregationResult{}
		}
		t.FixedInterval = intervalDuration.String()
	}
	log.Debugf("Starting aggregation with interval: %s and field: %s", t.FixedInterval, t.Field)

	log.Debugf("Sort rule: %s", sortRule)

	var buckets []*models.Bucket
	var currentBucket *models.Bucket

	// Function to create a new bucket
	createNewBucket := func(t time.Time) {
		key := t.UnixMilli()
		keyAsString := t.Format(time.RFC3339)
		currentBucket = &models.Bucket{
			KeyAsString: keyAsString,
			Key:         key,
			DocCount:    0,
		}
		buckets = append(buckets, currentBucket)
	}

	// Iterate over hits
	for _, hit := range hits {
		// Check if we need a new bucket
		hitTime, err := t.getHitTime(hit, t.Field)
		if err != nil {
			log.Warnf("Invalid hit time: %v", err)
			continue
		}
		//if currentBucket == nil || hitTime.After(time.UnixMilli(currentBucket.Key).Add(intervalDuration)) {
		//	createNewBucket(hitTime)
		//}
		if currentBucket == nil {
			createNewBucket(hitTime)
		} else {
			backetStartTime := time.UnixMilli(currentBucket.Key)
			if sortRule == "asc" {
				if hitTime.After(backetStartTime.Add(intervalDuration)) {
					createNewBucket(hitTime)
				}
			} else if sortRule == "desc" {
				if hitTime.Before(backetStartTime.Add(-intervalDuration)) {
					createNewBucket(hitTime)
				}
			} else {
				log.Warnf("Invalid sort rule: %s", sortRule)
				break
			}
		}

		// Increment document count
		currentBucket.DocCount++
	}

	return &models.AggregationResult{Buckets: buckets}
}

// TimeIntervalAggregatorFactory is an implementation of SearchHitAggregatorFactory that creates TimeIntervalAggregator
type TimeIntervalAggregatorFactory struct {
}

func (t *TimeIntervalAggregatorFactory) CreateAggregator(settings *models.SearchAggregation) ports.SearchHitAggregator {
	return NewTimeIntervalAggregator(settings.DateHistogram.Interval, settings.DateHistogram.Field, parseHitTime)
}

// NewHitAggregatorFactory creates a new TimeIntervalAggregatorFactory
func NewHitAggregatorFactory() *TimeIntervalAggregatorFactory {
	return &TimeIntervalAggregatorFactory{}
}

// FIXME - Rename this to search processor
type Aggregator struct {
	request       *models.SearchRequest
	result        *models.SearchResult
	searchResults []*models.SearchResult
	mux           *sync.Mutex
	getHitTime    ports.HitTimeParser
	aggr          ports.SearchHitAggregatorFactory
}

// NewAggregator creates a new aggregator
func NewAggregator(request *models.SearchRequest, aggr ports.SearchHitAggregatorFactory) *Aggregator {
	return &Aggregator{
		searchResults: make([]*models.SearchResult, 0),
		mux:           &sync.Mutex{},
		request:       request,
		result: &models.SearchResult{
			Took:         0,
			TimedOut:     false,
			Shards:       nil,
			Hits:         nil,
			Aggregations: make(map[string]*models.AggregationResult),
		},
		getHitTime: parseHitTime,
		aggr:       aggr,
	}
}

func (a *Aggregator) AddResult(result *models.SearchResult) {
	a.mux.Lock()
	a.searchResults = append(a.searchResults, result)
	a.mux.Unlock()
}

func parseHitTime(hit *models.Hit, fieldName string) (time.Time, error) {
	// check if iter already in hit.Sort
	//if hit.Sort == nil {
	//	hit.Sort = make([]int, sortIndex+1)
	//}
	//if len(hit.Sort) > sortIndex {
	//	return time.Unix(int64(hit.Sort[sortIndex]), 0), nil
	//}

	value, ok := hit.Source[fieldName]
	if !ok {
		return time.Time{}, nil
	}
	switch value.(type) {
	case string:
		result, err := time.Parse(time.RFC3339, value.(string))
		//if err == nil {
		//	// insert to hit.Sort
		//	hit.Sort[sortIndex] = int(result.Unix())
		//}
		return result, err
	case time.Time:
		return value.(time.Time), nil
	default:
		return time.Time{}, nil
	}
}

// Sort sorts the hits by request parameters
func (a *Aggregator) Sort() {
	// TODO add sort abstraction
	if a.request.Sort == nil || len(a.request.Sort) == 0 {
		return
	}

	for _, sortRule := range a.request.Sort {
		for fieldName, sortOrder := range sortRule {
			// Sorting logic
			sort.SliceStable(a.result.Hits.Hits, func(i, j int) bool {
				hitI := a.result.Hits.Hits[i]
				hitJ := a.result.Hits.Hits[j]

				timeI, errI := parseHitTime(hitI, fieldName)
				timeJ, errJ := parseHitTime(hitJ, fieldName)

				if errI != nil || errJ != nil {
					return false
				}
				hitI.Sort = append(hitI.Sort, int(timeI.Unix()))
				hitJ.Sort = append(hitJ.Sort, int(timeI.Unix()))

				if sortOrder.Order == "desc" {
					return timeI.After(timeJ)
				}
				return timeI.Before(timeJ)
			})
		}
	}
}

// Aggregate fill in the aggregations
func (a *Aggregator) Aggregate() {
	log.Warnf("Aggregation is not implemented yet %+v", a.request.Aggregations)
	for name, settings := range a.request.Aggregations {
		log.Debugf("Aggregating %s, %+v", name, settings.DateHistogram)
		aggr := a.aggr.CreateAggregator(settings)
		a.result.Aggregations[name] = aggr.Aggregate(a.result.Hits.Hits)
	}
}

func (a *Aggregator) GetResult() *models.SearchResult {
	a.mux.Lock()
	defer a.mux.Unlock()
	took := 0
	shards := &models.Shards{
		Total:      0,
		Successful: 0,
		Skipped:    0,
	}

	hits := &models.Hits{
		Total: &models.TotalValue{
			Value: 0,
		},
		Hits: make([]*models.Hit, 0),
	}

	for _, result := range a.searchResults {
		took += result.Took
		shards.Total += result.Shards.Total
		shards.Failed += result.Shards.Failed
		shards.Skipped += result.Shards.Skipped
		shards.Successful += result.Shards.Successful
		hits.Total.Value += result.Hits.Total.Value
		hits.Hits = append(hits.Hits, result.Hits.Hits...) // Not ordered yet
	}

	a.result.Took = took
	a.result.Shards = shards
	a.result.Hits = hits

	a.Sort()
	a.Aggregate()

	if a.request.Size > 0 && len(a.result.Hits.Hits) > a.request.Size {
		log.Debugf("Limiting result to %d", a.request.Size)
		a.result.Hits.Hits = a.result.Hits.Hits[:a.request.Size]
	}
	return a.result
}

type AggregatorFactory struct {
}

func (a *AggregatorFactory) CreateAggregator(request *models.SearchRequest) ports.SearchAggregator {
	return NewAggregator(request, NewHitAggregatorFactory())
}

func NewAggregatorFactory() *AggregatorFactory {
	return &AggregatorFactory{}
}
