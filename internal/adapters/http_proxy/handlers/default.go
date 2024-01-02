package handlers

import (
	"bytes"
	"crypto/tls"
	"io"
	"net/http"
	"strings"
)

func DefaultHandler(dest string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//log.Debugf("Proxy request: %s %s %s\n", r.Method, r.URL.Path, r.Proto)
		requestBody, _ := io.ReadAll(r.Body)
		destReq, err := http.NewRequest(r.Method, dest+r.URL.Path, bytes.NewBuffer(requestBody))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		destReq.Header = r.Header
		// Disable TLS verification
		// TODO: make it configurable
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		// Send the request
		client := &http.Client{
			Transport: tr,
		}
		resp, err := client.Do(destReq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
		responseBody, _ := io.ReadAll(resp.Body)
		// ignore _nodes requests
		if r.URL.Path != "/_nodes" {
			//log.Debugf("Request Body: %s", string(requestBody))
			//log.Debugf("Response Body: %s", string(responseBody))
		}
		// Write the response back to the original client
		for key, value := range resp.Header {
			w.Header().Set(key, strings.Join(value, ", "))
		}
		w.WriteHeader(resp.StatusCode)
		_, _ = w.Write(responseBody)
	}
}
