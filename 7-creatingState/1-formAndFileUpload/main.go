package main

import (
	"bytes"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var parsedTemplates *template.Template

func init() {
	tmpl, err := template.New("").ParseGlob("./templates/*")
	parsedTemplates = template.Must(tmpl, err)
}

func main() {
	http.Handle("/favicon.ico", http.NotFoundHandler())

	http.HandleFunc("/foo", func(w http.ResponseWriter, req *http.Request) {
		v := req.FormValue("q")
		io.WriteString(w, "Do search on: "+v)
	})

	http.HandleFunc("/bar", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte(getResponseBody("index.gohtml", nil)))

		f, h, err := req.FormFile("file")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer f.Close()

		bs, err := ioutil.ReadAll(f)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		dst, err := os.Create(filepath.Join("./user/", h.Filename))
		if err != nil {
			log.Fatal("Error creating a file", err)
		}
		defer dst.Close()

		_, err = dst.Write(bs)
		if err != nil {
			log.Fatal("Error creating a file", err)
		}

		n := req.FormValue("name")
		ln := req.FormValue("lastName")
		sb := req.FormValue("subscribe")

		if n != "" && ln != "" {
			io.WriteString(w, "Your name: "+n)
			io.WriteString(w, " and lastname: "+ln)
			io.WriteString(w, " and your subscription status: "+sb)
		} else {
			io.WriteString(w, "Please add both name and lastname")
		}
	})

	http.ListenAndServe(":8080", nil)
}

func getResponseBody(tn string, d interface{}) string {
	var tpl bytes.Buffer

	err := parsedTemplates.ExecuteTemplate(&tpl, tn, d)
	if err != nil {
		log.Fatal(err)
	}

	return tpl.String()
}
