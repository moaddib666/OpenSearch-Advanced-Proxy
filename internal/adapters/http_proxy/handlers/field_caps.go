package handlers

import (
	"OpenSearchAdvancedProxy/internal/core/models"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// FieldCapsHandler is a custom handler for the /_field_caps endpoint
func FieldCapsHandler(fields *models.Fields) http.HandlerFunc {
	raw, err := json.Marshal(fields)
	if err != nil {
		log.Fatalf("error marshalling fields: %v", err)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		log.Debugf("Manually handling request: %s %s %s\n", r.Method, r.URL.Path, r.Proto)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write(raw)
		if err != nil {
			log.Errorf("error writing response: %v", err)
		}
	}
}
