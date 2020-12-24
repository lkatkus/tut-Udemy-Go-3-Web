package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"text/template"
)

type myMux struct{}

func (m myMux) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	fmt.Println("Request:", req.Method, req.URL.Path)

	if req.Method == http.MethodGet {
		res.Header().Set("My-Key", "Well this is workign")
		// res.Header().Set("Content-Type", "text/html; charset=utf-8")

		switch req.URL.Path {
		case "/":
			res.Write([]byte(getResponseBody("index.gohtml", nil)))
		case "/first":
			res.Write([]byte(getResponseBody("first.gohtml", nil)))
		case "/second":
			fmt.Fprintln(res, getResponseBody("second.gohtml", nil))
		default:
			fmt.Fprintln(res, `<h1>404</h1>`)
		}
	}

	if req.Method == http.MethodPost {
		err := req.ParseForm()
		if err != nil {
			log.Fatalln(err)
		}

		res.Write([]byte(getResponseBody("name.gohtml", req.Form)))
	}

}

var parsedTemplates *template.Template

func init() {
	tmpl, err := template.New("").ParseGlob("./templates/*")
	parsedTemplates = template.Must(tmpl, err)
}

func main() {
	mux := myMux{}

	http.ListenAndServe("localhost:8080", mux)
}

func getResponseBody(tn string, d interface{}) string {
	var tpl bytes.Buffer

	err := parsedTemplates.ExecuteTemplate(&tpl, tn, d)
	if err != nil {
		log.Fatal(err)
	}

	return tpl.String()
}
