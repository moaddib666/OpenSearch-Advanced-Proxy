package log_provider

import (
	"OpenSearchAdvancedProxy/internal/core/models"
	"OpenSearchAdvancedProxy/internal/core/ports"
	"bufio"
	log "github.com/sirupsen/logrus"
	"os"
	"sync"
	"time"
)

type LogFileProvider struct {
	file    string
	scanner *bufio.Scanner
	fh      *os.File

	filterFactory    ports.SearchFilterFactory
	entryConstructor ports.EntryConstructor
	indexer          ports.Indexer
	intervalParser   ports.SearchInternalParser

	//endPos   int64
	startPos int64
	mux      sync.Mutex

	timeRange *models.Range

	currentFilter   ports.SearchFilter
	currentLogEntry ports.LogEntry

	returnedResults int
	limitResults    int
}

// NewLogFileProvider creates a new LogFileProvider struct
func NewLogFileProvider(filePath string, constructor ports.EntryConstructor, indexer ports.Indexer, intervalParser ports.SearchInternalParser, filterFactory ports.SearchFilterFactory) *LogFileProvider {
	return &LogFileProvider{
		file:             filePath,
		entryConstructor: constructor,
		indexer:          indexer,
		mux:              sync.Mutex{},
		intervalParser:   intervalParser,
		filterFactory:    filterFactory,
	}
}

func (f *LogFileProvider) BeginScan(r *models.SearchRequest) {
	// TODO: due to nature of the file it's more efficient to scan from the end;
	var err error
	f.mux.Lock()
	f.open()
	f.timeRange = r.GetRange()
	f.limitResults = r.Size
	if f.indexer != nil {
		err := f.indexer.LoadOrCreateIndex()
		if err != nil {
			log.Warnf("Index was not loaded: %s", err.Error())
			return
		}
		f.startPos, err = f.indexer.SearchStartPos(f.timeRange.DateTime.GTE)
		if err == nil {
			log.Debugf("Start position: %d", f.startPos)
			_, err = f.fh.Seek(f.startPos, 0)
			if err != nil {
				log.Warnf("Error seeking to start position: %s", err.Error())
			}
		}
	}
	f.currentFilter, err = f.filterFactory.FromQuery(r.Query)
	if err != nil {
		log.Errorf("Error creating filter: %s", err.Error())
	}
}

func (f *LogFileProvider) EndScan() {
	f.currentLogEntry = nil
	f.currentFilter = nil
	f.timeRange = nil
	f.returnedResults = 0
	f.limitResults = 0
	f.close()
	f.mux.Unlock()
}

func (f *LogFileProvider) LogEntry() ports.LogEntry {
	return f.currentLogEntry
}

func (f *LogFileProvider) Scan() bool {
	if f.startPos == -1 {
		return false
	}
	if f.returnedResults > f.limitResults {
		return false
	}
	for f.scanner.Scan() {
		entry := f.entryConstructor()
		_ = entry.LoadBytes(f.scanner.Bytes())
		if entry.Timestamp().After(f.timeRange.DateTime.LTE) {
			return false
		}
		if entry.Timestamp().Before(f.timeRange.DateTime.GTE) {
			continue
		}
		if !f.currentFilter.Match(entry) {
			continue
		}
		f.currentLogEntry = entry
		f.returnedResults++
		return true
	}
	return false
}

func (f *LogFileProvider) Text() string {
	return f.currentLogEntry.Raw()
}

func (f *LogFileProvider) Err() error {
	return f.scanner.Err()
}

// open a file and return a scanner
func (f *LogFileProvider) open() {
	f.close()
	log.Debugf("Opening file: %s", f.file)
	file, err := os.Open(f.file)
	if err != nil {
		panic(err)
	}
	f.fh = file
	f.scanner = bufio.NewScanner(file)
}

// close the file
func (f *LogFileProvider) close() {
	log.Debugf("Closing file: %s", f.file)
	_ = f.fh.Close()
	f.fh = nil
	f.scanner = nil
}

func (f *LogFileProvider) AggregateResult(request *models.SearchAggregation) *models.AggregationResult {
	result := models.NewAggregationResult()
	if request.DateHistogram == nil {
		return result
	}
	var interval time.Duration
	err := f.intervalParser.Parse(request.DateHistogram, &interval)
	if err != nil {
		log.Warnf("Error while parsing interval: %s", err)
		return result
	}
	duration := f.timeRange.DateTime.LTE.Sub(f.timeRange.DateTime.GTE)
	intervalCount := int(duration / interval)
	buckets := make([]*models.Bucket, intervalCount)
	for i := range buckets {
		ts := f.timeRange.DateTime.GTE.Add(time.Duration(i) * interval)
		b := models.NewBucket()
		b.FromTime(ts)
		buckets[i] = b
	}

	for f.scanner.Scan() {
		entry := f.entryConstructor()
		err = entry.LoadBytes(f.scanner.Bytes())
		if err != nil {
			log.Warnf("Error while loading entry: %s", err)
			continue
		}

		if entry.Timestamp().Before(f.timeRange.DateTime.GTE) {
			continue
		}

		if entry.Timestamp().After(f.timeRange.DateTime.LTE) {
			break
		}

		if !f.currentFilter.Match(entry) {
			continue
		}

		bucketIndex := int(entry.Timestamp().Sub(f.timeRange.DateTime.GTE) / interval)
		if bucketIndex < 0 || bucketIndex >= intervalCount {
			continue
		}
		buckets[bucketIndex].AddDoc()
	}
	for _, bucket := range buckets {
		if bucket.HasDocs() {
			result.AddBucket(bucket)
		}
	}
	return result
}
