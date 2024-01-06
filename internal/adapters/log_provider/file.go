package log_provider

import (
	"OpenSearchAdvancedProxy/internal/core/models"
	"OpenSearchAdvancedProxy/internal/core/ports"
	"bufio"
	log "github.com/sirupsen/logrus"
	"os"
	"sync"
)

type LogFileProvider struct {
	file             string
	scanner          *bufio.Scanner
	fh               *os.File
	entryConstructor ports.EntryConstructor
	indexer          ports.Indexer
	endPos           int64
	startPos         int64
	mux              sync.Mutex
}

func (f *LogFileProvider) BeginScan(r *models.SearchRequest) {
	f.mux.Lock()
	f.open()
	rg := r.GetRange()
	if f.indexer != nil {
		err := f.indexer.LoadOrCreateIndex()
		if err != nil {
			log.Warnf("Index was not loaded: %s", err.Error())
			return
		}
		f.startPos, err = f.indexer.SearchStartPos(rg.DateTime.GTE)
		if err == nil {
			log.Debugf("Start position: %d", f.startPos)
			_, err = f.fh.Seek(f.startPos, 0)
			if err != nil {
				log.Warnf("Error seeking to start position: %s", err.Error())
			}
		}
	}
}

func (f *LogFileProvider) EndScan() {
	f.close()
	f.mux.Unlock()
}

func (f *LogFileProvider) LogEntry() ports.LogEntry {
	entry := f.entryConstructor()
	_ = entry.LoadString(f.scanner.Text())
	return entry
}

func (f *LogFileProvider) Scan() bool {
	if f.startPos == -1 {
		return false
	}
	return f.scanner.Scan()
}

func (f *LogFileProvider) Text() string {
	return f.scanner.Text()
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
	return nil
}

// NewLogFileProvider creates a new LogFileProvider struct
func NewLogFileProvider(filePath string, constructor ports.EntryConstructor, indexer ports.Indexer) *LogFileProvider {
	return &LogFileProvider{
		file:             filePath,
		entryConstructor: constructor,
		indexer:          indexer,
		mux:              sync.Mutex{},
	}
}
