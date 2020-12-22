package main

import (
	"fmt"
	"log"
	"os"
	"text/template"
	"time"
)

var parsedTemplates *template.Template
var fm = template.FuncMap{
	"fdateMDY":     fdateMDY,
	"appendString": appendString,
}

func fdateMDY(t time.Time) string {
	return t.Format("2006-01-02")
}

func appendString(s string) string {
	return fmt.Sprint(s, "-appended-")
}

func init() {
	os.Mkdir("html", os.ModePerm)

	tmpl, err := template.New("").Funcs(fm).ParseGlob("./templates/*")
	parsedTemplates = template.Must(tmpl, err)
}

func main() {
	printToFile("./html/date.html", "index.gohtml", time.Now())
}

func printToFile(fn string, tn string, data interface{}) {
	f, err := os.Create(fn)
	if err != nil {
		log.Fatal("Failed to create a file", err)
	}
	defer f.Close()

	err = parsedTemplates.ExecuteTemplate(f, tn, data)
	if err != nil {
		log.Fatal(err)
	}
}
