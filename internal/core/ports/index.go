package ports

import "time"

type Indexer interface {
	GetIndexName() string
	CreateIndex() error
	LoadIndex() error
	LoadOrCreateIndex() error
	ReIndex() error
	SearchStartPos(ts time.Time) (int64, error)
	SearchEndPos(ts time.Time) (int64, error)
}
