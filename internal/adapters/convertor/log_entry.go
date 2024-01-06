package convertor

import (
	"OpenSearchAdvancedProxy/internal/core/models"
	"OpenSearchAdvancedProxy/internal/core/ports"
)

type DefaultLogEntryConvertor struct {
	Index string
}

func NewDefaultLogEntryConvertor(index string) *DefaultLogEntryConvertor {
	return &DefaultLogEntryConvertor{
		Index: index,
	}
}

func (d *DefaultLogEntryConvertor) Convert(entry ports.LogEntry) (*models.Hit, error) {
	return models.NewHit(d.Index, entry.ID(), entry.Map(), entry.Timestamp()), nil
}

func (d *DefaultLogEntryConvertor) ConvertBatch(entries []ports.LogEntry) (*models.Hits, error) {
	count := len(entries)
	hitList := make([]*models.Hit, count)
	for index, entry := range entries {
		hitList[index], _ = d.Convert(entry)
	}
	hits := &models.Hits{
		Total: &models.TotalValue{
			Value:    count,
			Relation: "",
		},
		Hits: hitList,
	}
	return hits, nil
}
