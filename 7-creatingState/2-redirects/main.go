package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	http.HandleFunc("/foo", func(res http.ResponseWriter, req *http.Request) {
		fmt.Println("REQUEST - /foo")

		res.Header().Set("Location", "/bar")
		res.WriteHeader(http.StatusSeeOther)
	})

	http.HandleFunc("/bar", func(res http.ResponseWriter, req *http.Request) {
		http.Redirect(res, req, "/foobar", http.StatusSeeOther)
	})

	http.HandleFunc("/foobar", func(res http.ResponseWriter, req *http.Request) {
		io.WriteString(res, "Welcome to foobar")
	})

	http.ListenAndServe(":8080", nil)
}
