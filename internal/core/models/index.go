package models

import "time"

type Index struct {
	Entries []*IndexEntry `json:"entries"`
	Step    time.Duration `json:"step"`
}

type IndexEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Position  int64     `json:"position"`
}
