package main

import (
	"io"
	"net/http"
)

func main() {
	http.HandleFunc("/", dogHandler)
	http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("./assets"))))

	http.ListenAndServe(":8080", nil)
}

func dogHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.WriteString(res, `<div><h1>Dogo</h1><img src="assets/dog.jpg"></div>`)
}
