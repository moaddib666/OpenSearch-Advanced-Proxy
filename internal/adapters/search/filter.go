package search

import (
	"OpenSearchAdvancedProxy/internal/core/models"
	"OpenSearchAdvancedProxy/internal/core/ports"
	"fmt"
	log "github.com/sirupsen/logrus"
	"strings"
)

// MatchAllFilter is a filter that matches all
type MatchAllFilter struct {
}

func (m *MatchAllFilter) Match(entry ports.LogEntry) bool {
	return true
}

// NewMatchAllFilter creates a new match-all filter
func NewMatchAllFilter() *MatchAllFilter {
	return &MatchAllFilter{}
}

// BoolFilter is a filter that matches a boolean value
type BoolFilter struct {
	filters []ports.SearchFilter
}

// NewBoolFilter creates a new boolean filter
func NewBoolFilter(filters []ports.SearchFilter) *BoolFilter {
	return &BoolFilter{
		filters: filters,
	}
}

func (b *BoolFilter) Match(entry ports.LogEntry) bool {
	for _, filter := range b.filters {
		if !filter.Match(entry) {
			return false
		}
	}
	return true
}

// ExcludeFilter is a filter that excludes a value
type ExcludeFilter struct {
	filter ports.SearchFilter
}

func (e *ExcludeFilter) Match(entry ports.LogEntry) bool {
	return !e.filter.Match(entry)
}

func NewExcludeFilter(filter ports.SearchFilter) *ExcludeFilter {
	return &ExcludeFilter{
		filter: filter,
	}
}

// MultiMatchFilter is a filter that matches a multi-match query
type MultiMatchFilter struct {
	model *models.MultiMatch
}

func (m *MultiMatchFilter) Match(entry ports.LogEntry) bool {
	return strings.Contains(entry.Raw(), m.model.Query)
}

// NewMultiMatchFilter creates a new multi-match filter
func NewMultiMatchFilter(m *models.MultiMatch) *MultiMatchFilter {
	return &MultiMatchFilter{
		model: m,
	}
}

// FilterFactory is a factory for filters
type FilterFactory struct {
}

type MatchPhraseFilter struct {
	model models.MatchPhrase
}

func (m *MatchPhraseFilter) Match(entry ports.LogEntry) bool {
	fieldsMap := entry.Map()
	for k, v := range m.model {
		if entryValue, ok := fieldsMap[k]; !ok {
			return false
		} else if entryValue != v {
			return false
		}
	}
	return true
}
func NewMatchPhraseFilter(m models.MatchPhrase) *MatchPhraseFilter {
	return &MatchPhraseFilter{
		model: m,
	}
}

type RangeFilter struct {
	model *models.Range
}

func (r *RangeFilter) Match(entry ports.LogEntry) bool {
	if entry.Timestamp().Before(r.model.DateTime.GTE) {
		return false
	}
	if entry.Timestamp().After(r.model.DateTime.LTE) {
		return false
	}
	return true
}

// NewRangeFilter creates a new range filter
func NewRangeFilter(m *models.Range) *RangeFilter {
	return &RangeFilter{
		model: m,
	}
}

// CreateFilter creates a new filter
func (f *FilterFactory) CreateFilter(filter *models.Filter) (ports.SearchFilter, error) {
	if filter == nil {
		log.Debugf("Filter is nil, returning match-all filter")
		return NewMatchAllFilter(), nil
	}
	if filter.Bool != nil {
		bf, err := f.CreateBoolFilter(filter.Bool)
		if err != nil {
			return nil, err
		}
		return bf, nil
	}
	if filter.MultiMatch != nil {
		return NewMultiMatchFilter(filter.MultiMatch), nil
	}
	if filter.MatchPhrase != nil {
		return NewMatchPhraseFilter(filter.MatchPhrase), nil
	}
	if filter.MatchAll != nil {
		return NewMatchAllFilter(), nil
	}
	if filter.Range != nil {
		return NewRangeFilter(filter.Range), nil
	}
	log.Debugf("Unsupported filter type: %+v", filter)
	return NewMatchAllFilter(), nil
	//return nil, fmt.Errorf("unsupported filter type")
}

// CreateBoolFilter creates a new boolean filter
func (f *FilterFactory) CreateBoolFilter(q *models.BoolFilter) (ports.SearchFilter, error) {
	currentFilter := make([]ports.SearchFilter, len(q.Filter))
	for i, nestedFilter := range q.Filter {
		nested, err := f.CreateFilter(nestedFilter)
		if err != nil {
			return nil, err
		}
		currentFilter[i] = nested
	}
	if q.MustNot != nil {
		nested, err := f.CreateFilter(q.MustNot)
		if err != nil {
			return nil, err
		}
		currentFilter = append(currentFilter, NewExcludeFilter(nested))
	}
	return NewBoolFilter(currentFilter), nil
}

// FromQuery creates a new filter from a query
func (f *FilterFactory) FromQuery(query *models.Query) (ports.SearchFilter, error) {
	filterSet := make([]ports.SearchFilter, 0)
	if query == nil {
		log.Debugf("Query is nil, returning match-all filter")
		return NewMatchAllFilter(), nil
	}
	if query.Bool != nil {
		for _, filter := range query.Bool.Filter {
			f, err := f.CreateFilter(filter)
			if err != nil {
				return nil, err
			}
			filterSet = append(filterSet, f)
		}
		return NewBoolFilter(filterSet), nil
	}

	return nil, fmt.Errorf("unsupported query type")
}

// NewFilterFactory creates a new filter factory
func NewFilterFactory() *FilterFactory {
	return &FilterFactory{}
}
