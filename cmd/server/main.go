package main

import (
	"net/http"

	"github.com/andrei-daniloff/monitoring/internal/handler"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/update/{type}/{name}/{value}", handler.Update)

	http.ListenAndServe(":8080", mux)
}
