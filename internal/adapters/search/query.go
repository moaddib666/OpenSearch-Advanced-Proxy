package search

import (
	"fmt"
	"github.com/moaddib666/OpenSearch-Advanced-Proxy/internal/core/models"
	"github.com/moaddib666/OpenSearch-Advanced-Proxy/internal/core/ports"
	log "github.com/sirupsen/logrus"
	"strings"
)

type SQLBoolQueryBuilder struct{}

type SQLMatchAllQueryBuilder struct{}

func NewMatchAllQueryBuilder() *SQLMatchAllQueryBuilder {
	return &SQLMatchAllQueryBuilder{}
}

func (S *SQLMatchAllQueryBuilder) BuildQuery() (string, error) {
	return "", nil
}

type SQLMatchPhraseQueryBuilder struct {
	model  models.MatchPhrase
	fields []string
}

func NewMatchPhraseQueryBuilder(model models.MatchPhrase, fields []string) *SQLMatchPhraseQueryBuilder {
	return &SQLMatchPhraseQueryBuilder{
		model:  model,
		fields: fields,
	}
}

func (S *SQLMatchPhraseQueryBuilder) BuildQuery() (string, error) {
	var conditions []string
	for _, field := range S.fields {
		filterValue, ok := S.model[field]
		if !ok {
			continue
		}
		conditions = append(conditions, fmt.Sprintf("%s = '%s'", field, filterValue))
	}
	return strings.Join(conditions, " OR "), nil
}

type SQLMultiMatchQueryBuilder struct {
	model  *models.MultiMatch
	fields []string
}

func NewMultiMatchQueryBuilder(model *models.MultiMatch, fields []string) *SQLMultiMatchQueryBuilder {
	return &SQLMultiMatchQueryBuilder{
		model:  model,
		fields: fields,
	}
}

func (S *SQLMultiMatchQueryBuilder) BuildQuery() (string, error) {
	var conditions []string
	for _, field := range S.fields {
		conditions = append(conditions, fmt.Sprintf("%s LIKE '%%%s%%'", field, S.model.Query))
	}
	return strings.Join(conditions, " OR "), nil
}

type SQLExcludeQueryBuilder struct {
	condition ports.QueryBuilder
}

func NewExcludeQueryBuilder(condition ports.QueryBuilder) *SQLExcludeQueryBuilder {
	return &SQLExcludeQueryBuilder{
		condition: condition,
	}
}

func (S *SQLExcludeQueryBuilder) BuildQuery() (string, error) {
	condition, err := S.condition.BuildQuery()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("NOT (%s)", condition), nil
}

type SQLRangeQueryBuilder struct {
	model         *models.Range
	dateTimeField string
}

func NewRangeQueryBuilder(model *models.Range, dateTimeField string) *SQLRangeQueryBuilder {
	return &SQLRangeQueryBuilder{
		model:         model,
		dateTimeField: dateTimeField,
	}
}
func (S *SQLRangeQueryBuilder) BuildQuery() (string, error) {
	var query string
	if S.model.DateTime != nil {
		if !S.model.DateTime.GTE.IsZero() {
			query += fmt.Sprintf("%s >= %d", S.dateTimeField, S.model.DateTime.GTE.Unix())
		}
		if !S.model.DateTime.LTE.IsZero() {
			if query != "" {
				query += " AND "
			}
			query += fmt.Sprintf("%s <= %d", S.dateTimeField, S.model.DateTime.LTE.Unix())
		}
	}
	return query, nil
}

type SQLBoolConditionQueryBuilder struct {
	conditions []ports.QueryBuilder
}

func (S *SQLBoolConditionQueryBuilder) BuildQuery() (string, error) {
	var conditions []string
	for _, condition := range S.conditions {
		conditionQuery, err := condition.BuildQuery()
		if err != nil {
			return "", err
		}
		if conditionQuery == "" {
			continue
		}
		conditionQuery = "(" + conditionQuery + ")"
		conditions = append(conditions, conditionQuery)
	}
	return strings.Join(conditions, " AND "), nil
}

func NewBoolConditionQueryBuilder(conditions []ports.QueryBuilder) *SQLBoolConditionQueryBuilder {
	return &SQLBoolConditionQueryBuilder{
		conditions: conditions,
	}
}

type SQLLimitQueryBuilder struct {
	limit int
}

func (S *SQLLimitQueryBuilder) BuildQuery() (string, error) {
	return fmt.Sprintf("LIMIT %d", S.limit), nil
}

type SQLQueryBuilderFactory struct {
	baseQueryString  string
	searchableFields []string
	dateTimeField    string
}

func NewSQLQueryBuilderFactory(searchableFields []string, dateTimeField string) *SQLQueryBuilderFactory {
	return &SQLQueryBuilderFactory{
		dateTimeField:    dateTimeField,
		searchableFields: searchableFields,
		baseQueryString:  string("SELECT * FROM logs WHERE "), // FIXME
	}
}

func (S *SQLQueryBuilderFactory) CreateQueryBuilder(filter *models.Filter) (ports.QueryBuilder, error) {
	if filter == nil {
		log.Debugf("Filter is nil, returning match-all filter")
		return NewMatchAllQueryBuilder(), nil
	}
	if filter.Bool != nil {
		bf, err := S.CreateBoolConditionBuilder(filter.Bool)
		if err != nil {
			return nil, err
		}
		return bf, nil
	}
	if filter.MultiMatch != nil {
		return NewMultiMatchQueryBuilder(filter.MultiMatch, S.searchableFields), nil
	}
	if filter.MatchPhrase != nil {
		return NewMatchPhraseQueryBuilder(filter.MatchPhrase, S.searchableFields), nil
	}
	if filter.MatchAll != nil {
		return NewMatchAllQueryBuilder(), nil
	}
	if filter.Range != nil {
		return NewRangeQueryBuilder(filter.Range, S.dateTimeField), nil
	}
	log.Debugf("Unsupported filter type: %+v", filter)
	return NewMatchAllQueryBuilder(), nil
}

func (S *SQLQueryBuilderFactory) CreateBoolConditionBuilder(filter *models.BoolFilter) (ports.QueryBuilder, error) {
	currentFilter := make([]ports.QueryBuilder, len(filter.Filter))
	for i, nestedFilter := range filter.Filter {
		nested, err := S.CreateQueryBuilder(nestedFilter)
		if err != nil {
			return nil, err
		}
		currentFilter[i] = nested
	}
	if filter.MustNot != nil {
		nested, err := S.CreateQueryBuilder(filter.MustNot)
		if err != nil {
			return nil, err
		}
		currentFilter = append(currentFilter, NewExcludeQueryBuilder(nested))
	}
	return NewBoolConditionQueryBuilder(currentFilter), nil
}

func (S *SQLQueryBuilderFactory) FromQuery(query *models.Query) (ports.QueryBuilder, error) {
	querySet := make([]ports.QueryBuilder, 0)
	if query == nil {
		log.Debugf("Query is nil, returning match-all filter")
		return NewMatchAllQueryBuilder(), nil
	}
	if query.Bool != nil {
		for _, filter := range query.Bool.Filter {
			f, err := S.CreateQueryBuilder(filter)
			if err != nil {
				return nil, err
			}
			querySet = append(querySet, f)
		}
		for _, filter := range query.Bool.MustNot {
			f, err := S.CreateQueryBuilder(filter)
			if err != nil {
				return nil, err
			}
			querySet = append(querySet, NewExcludeQueryBuilder(f))
		}
		for _, filter := range query.Bool.Must {
			f, err := S.CreateQueryBuilder(filter)
			if err != nil {
				return nil, err
			}
			querySet = append(querySet, f)
		}
		return NewBoolConditionQueryBuilder(querySet), nil

	}
	return nil, fmt.Errorf("unsupported query type")
}
