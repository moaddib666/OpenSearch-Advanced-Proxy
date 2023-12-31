package log_provider

import (
	"OpenSearchAdvancedProxy/internal/core/ports"
	"bufio"
	log "github.com/sirupsen/logrus"
	"os"
)

type LogFileProvider struct {
	file             string
	scanner          *bufio.Scanner
	fh               *os.File
	entryConstructor ports.EntryConstructor
}

func (f *LogFileProvider) BeginScan() {
	f.open()
}

func (f *LogFileProvider) EndScan() {
	f.close()
}

func (f *LogFileProvider) LogEntry() ports.LogEntry {
	entry := f.entryConstructor()
	_ = entry.Load(f.scanner.Text())
	return entry
}

func (f *LogFileProvider) Scan() bool {
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

// NewLogFileProvider creates a new LogFileProvider struct
func NewLogFileProvider(filePath string, constructor ports.EntryConstructor) *LogFileProvider {
	return &LogFileProvider{
		file:             filePath,
		entryConstructor: constructor,
	}
}
