// Package main demonstrates how to use the lambdahttp package to create a REST API
// using standard Go HTTP handlers in AWS Lambda.
package main

import (
	"encoding/json"
	"net/http"

	"github.com/tjamet/lambdahttp"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		originalRequest := lambdahttp.GetOriginalRequest(r)
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(originalRequest); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	lambdahttp.StartLambdaHandler(mux)
}
