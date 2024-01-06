package log_provider

import (
	"OpenSearchAdvancedProxy/internal/adapters/search/search_internval"
	"OpenSearchAdvancedProxy/internal/core/models"
	"OpenSearchAdvancedProxy/internal/core/ports"
	"database/sql"
	"fmt"
	_ "github.com/ClickHouse/clickhouse-go"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

type SQLDatabaseProvider struct {
	queryBuilder     ports.QueryBuilderFactory
	db               *sql.DB
	table            string
	rows             *sql.Rows
	columns          []string
	err              error
	mux              sync.Mutex
	filterQuery      ports.QueryBuilder
	entryConstructor ports.EntryConstructor
	intervalParser   ports.SearchInternalParser
}

func NewSQLDatabaseProvider(queryBuilder ports.QueryBuilderFactory, db *sql.DB, table string, entryConstructor ports.EntryConstructor, iparser ports.SearchInternalParser) *SQLDatabaseProvider {
	return &SQLDatabaseProvider{
		queryBuilder:     queryBuilder,
		db:               db,
		table:            table,
		entryConstructor: entryConstructor,
		intervalParser:   iparser,
	}
}

func NewClickhouseProvider(dsn, table string, queryBuilder ports.QueryBuilderFactory, entryConstructor ports.EntryConstructor) *SQLDatabaseProvider {
	db, err := sql.Open("clickhouse", dsn)
	if err != nil {
		log.Fatalf("Error while opening database: %s", err)
	}
	return NewSQLDatabaseProvider(queryBuilder, db, table, entryConstructor, search_internval.NewClickHouseSearchIntervalParser())
}

func (s *SQLDatabaseProvider) Text() string {
	var text string
	err := s.rows.Scan(&text)
	if err != nil {
		log.Errorf("Error while scanning row: %s", err)
		return ""
	}
	return text
}

func (s *SQLDatabaseProvider) Err() error {
	if s.rows != nil {
		return s.rows.Err()
	}
	return s.err
}

func (s *SQLDatabaseProvider) BeginScan(r *models.SearchRequest) {
	var err error
	s.mux.Lock()
	log.Debugf("SQLDatabaseProvider scan started")
	s.filterQuery, err = s.queryBuilder.FromQuery(r.Query)
	if err != nil {
		log.Errorf("Error while building query: %s", err)
	}
	sqlString, err := s.filterQuery.BuildQuery()
	// FIXME: Extarnal base query
	baseQuery := fmt.Sprintf("SELECT * FROM %s WHERE %s", s.table, sqlString)
	limitedBaseQuery := fmt.Sprintf("%s LIMIT %d", baseQuery, r.Size)
	log.Debugf("Query for Hits: %s", limitedBaseQuery)
	s.rows, err = s.db.Query(limitedBaseQuery)
	if err != nil {
		log.Errorf("Error while querying database: %s", err)
	}
	s.columns, err = s.rows.Columns()
	if err != nil {
		log.Errorf("Error while getting columns: %s", err)
	}

}

func (s *SQLDatabaseProvider) Scan() bool {
	return s.rows.Next()
}

func (s *SQLDatabaseProvider) getDataFromSql() map[string]interface{} {
	columnsSlice := make([]interface{}, len(s.columns))
	columnsPtrs := make([]interface{}, len(s.columns))
	for i := range columnsSlice {
		columnsPtrs[i] = &columnsSlice[i]
	}

	// Scan the result into the column pointers...
	if err := s.rows.Scan(columnsPtrs...); err != nil {
		log.Errorf("Error while scanning row: %s", err)
		return nil
	}

	item := make(map[string]interface{})
	for i, colName := range s.columns {
		val := columnsPtrs[i].(*interface{})
		item[colName] = *val
	}
	return item
}

func (s *SQLDatabaseProvider) LogEntry() ports.LogEntry {
	item := s.getDataFromSql()
	if item == nil {
		log.Warnf("Error while getting data from sql")
		return nil
	}
	entry := s.entryConstructor()
	err := entry.LoadMap(item)
	if err != nil {
		log.Warnf("Error while loading map: %s", err)
		return nil
	}
	return entry
}

func (s *SQLDatabaseProvider) EndScan() {
	log.Debugf("SQLDatabaseProvider scan ended")
	s.mux.Unlock()
	if s.rows != nil {
		s.err = s.rows.Close()
	} else {
		s.err = nil
	}
	s.rows = nil
}

func (s *SQLDatabaseProvider) AggregateResult(request *models.SearchAggregation) *models.AggregationResult {
	// TODO: Create sepparate abstraction for this as it could changed from DB to DB
	result := &models.AggregationResult{
		Buckets: make([]*models.Bucket, 0),
	}
	if request.DateHistogram == nil {
		return result
	}
	// TODO: make supoort for other databases not on clickhouse
	//interval := request.DateHistogram.Interval
	// TODO: add resolve interval
	var interval string
	err := s.intervalParser.Parse(request.DateHistogram, &interval)
	if err != nil {
		log.Warnf("Error while parsing interval: %s", err)
		return result
	}
	field := request.DateHistogram.Field
	var filterBy string
	if s.filterQuery != nil {
		sqlString, err := s.filterQuery.BuildQuery()
		if err != nil {
			log.Errorf("Error while building query: %s", err)
		}
		filterBy = fmt.Sprintf("WHERE %s", sqlString)
	}
	baseQuery := fmt.Sprintf(fmt.Sprintf(`
	SELECT
	   toStartOfInterval(%s, INTERVAL %s) AS interval,
	   count() AS count
	FROM %s 
	%s
	GROUP BY interval
	ORDER BY interval
	`, field, interval, s.table, filterBy))

	log.Debugf("Query Metadata: %s", baseQuery)
	rows, err := s.db.Query(baseQuery)
	if err != nil {
		log.Errorf("Error while querying database: %s", err)
		return result
	}
	defer rows.Close()
	for rows.Next() {
		var ts time.Time
		bucket := &models.Bucket{}
		err := rows.Scan(&ts, &bucket.DocCount)
		if err != nil {
			log.Warnf("Error while scanning row: %s", err)
		}
		bucket.FromTime(ts)
		result.AddBucket(bucket)
	}
	return result
}
