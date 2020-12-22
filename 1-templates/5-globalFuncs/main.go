package main

import (
	"log"
	"os"
	"text/template"
)

var parsedTemplates *template.Template

func init() {
	os.Mkdir("html", os.ModePerm)

	tmpl, err := template.New("").ParseGlob("./templates/*")
	parsedTemplates = template.Must(tmpl, err)
}

func main() {
	type data struct {
		Slice     []string
		Ten, Five int
	}

	d := data{
		Slice: []string{"First", "Second", "Third"},
		Ten:   10,
		Five:  5,
	}

	printToFile("./html/globalFuncs.html", "index.gohtml", d)
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
