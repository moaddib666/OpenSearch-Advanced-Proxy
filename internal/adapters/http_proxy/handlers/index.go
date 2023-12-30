package handlers

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type index struct {
	Indices []struct {
		Name       string   `json:"name"`
		Attributes []string `json:"attributes"`
	} `json:"indices"`
	Aliases     []interface{} `json:"aliases"`
	DataStreams []interface{} `json:"data_streams"`
}

// IndexHandler is a custom handler for the /_resolve/index endpoint
// Response example: {"indices":[{"name":".kibana_1","aliases":[".kibana"],"attributes":["open"]},{"name":"opensearch_dashboards_sample_data_ecommerce","attributes":["open"]}],"aliases":[{"name":".kibana","indices":[".kibana_1"]}],"data_streams":[]}
func IndexHandler(name string) http.HandlerFunc {
	i := index{
		Indices: []struct {
			Name       string   `json:"name"`
			Attributes []string `json:"attributes"`
		}{
			{
				Name: name,
			},
		},
		Aliases:     []interface{}{},
		DataStreams: []interface{}{},
	}
	raw, err := json.Marshal(i)
	if err != nil {
		log.Fatalf("error marshalling index: %v", err)
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
