package handlers

import (
	"OpenSearchAdvancedProxy/internal/core/models"
	"OpenSearchAdvancedProxy/internal/core/ports"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// SearchHandler is a custom handler for the /custom_remote_infra_index/_search endpoint
func SearchHandler(storage ports.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Debugf("Manually handling request: %s %s %s\n", r.Method, r.URL.Path, r.Proto)
		result, err := storage.Search(&models.SearchRequest{})
		//
		if err != nil {
			log.Errorf("error searching storage: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		raw, err := json.Marshal(result)
		if err != nil {
			log.Errorf("error marshalling search result: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Debugf("Search result: %s", string(raw))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(raw)
	}
}