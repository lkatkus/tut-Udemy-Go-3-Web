package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func main() {
	router := httprouter.New()
	router.NotFound = handler404{}

	router.GET("/", index)
	router.GET("/hello/:name", hello)
	router.GET("/dog/:breed", dogs)

	log.Fatal(http.ListenAndServe(":8080", router))
}

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

func hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Println(ps)
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}

func dogs(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("breed"))
}

type handler404 struct{}

func (h handler404) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	io.WriteString(res, "--- 404 ---")
}
