package main

import (
	"log"
	"os"
	"text/template"
)

var parsedTemplates *template.Template

func init() {
	parsedTemplates = template.Must(template.ParseGlob("./templates/*"))
}

func main() {
	templateFromFile()
	templateFromGlob()
	fromInitialTemplates()
}

func templateFromFile() {
	tpl, err := template.ParseFiles("./templates/index.gohtml")
	if err != nil {
		log.Fatal("Failed to read template", err)
	}

	// Adding additional templates to a template element
	tpl, err = tpl.ParseFiles("./templates/other.gohtml")
	if err != nil {
		log.Fatal("Failed to read template", err)
	}

	f, err := os.Create("index.html")
	if err != nil {
		log.Fatal("Failed to create a file", err)
	}
	defer f.Close()

	err = tpl.ExecuteTemplate(f, "index.gohtml", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func templateFromGlob() {
	tpl, err := template.ParseGlob("./templates/*.gohtml")
	if err != nil {
		log.Fatal("Failed to read template", err)
	}

	// Adding additional templates to a template element
	tpl, err = tpl.ParseFiles("./templates/other.gohtml")
	if err != nil {
		log.Fatal("Failed to read template", err)
	}

	f, err := os.Create("glob.html")
	if err != nil {
		log.Fatal("Failed to create a file", err)
	}
	defer f.Close()

	err = tpl.ExecuteTemplate(f, "other.gohtml", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func fromInitialTemplates() {
	f, err := os.Create("fromInitial.html")
	if err != nil {
		log.Fatal("Failed to create a file", err)
	}
	defer f.Close()

	err = parsedTemplates.ExecuteTemplate(f, "other.gohtml", nil)
	if err != nil {
		log.Fatal(err)
	}
}
