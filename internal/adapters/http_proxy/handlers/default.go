package handlers

import (
	"bytes"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strings"
)

func DefaultHandler(dest string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Debugf("Proxy request: %s %s %s\n", r.Method, r.URL.Path, r.Proto)
		requestBody, _ := io.ReadAll(r.Body)
		destReq, err := http.NewRequest(r.Method, "http://"+dest+r.URL.Path, bytes.NewBuffer(requestBody))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		destReq.Header = r.Header

		// Send the request
		client := &http.Client{}
		resp, err := client.Do(destReq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
		responseBody, _ := io.ReadAll(resp.Body)

		// Write the response back to the original client
		for key, value := range resp.Header {
			w.Header().Set(key, strings.Join(value, ", "))
		}
		w.WriteHeader(resp.StatusCode)
		_, _ = w.Write(responseBody)
	}
}
