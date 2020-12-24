package main

import (
	"io"
	"net/http"
)

func main() {
	var indexHandler indexHandler
	var dogsHandler dogsHandler

	// mux := http.NewServeMux()
	// mux.Handle("/", indexHandler)
	// mux.Handle("/dog", dogsHandler)

	// http.ListenAndServe("localhost:8080", mux)

	http.Handle("/", indexHandler)
	http.Handle("/dogs", dogsHandler)     // Use Handler type
	http.HandleFunc("/cats", catsHandler) // Use function instead of Handler type

	http.ListenAndServe("localhost:8080", nil) // Use default ServerMux, when passing nil
}

type indexHandler struct{}

func (h indexHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	io.WriteString(res, "--- HOME ROUTE ---")
}

type dogsHandler struct{}

func (h dogsHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	io.WriteString(res, "--- DOGS ROUTE ---")
}

func catsHandler(res http.ResponseWriter, req *http.Request) {
	io.WriteString(res, "--- CATS ROUTE ---")
}
