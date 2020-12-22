package main

import (
	"fmt"
	html "html/template"
	"log"
	"os"
	text "text/template"
)

var parsedTemplates *text.Template
var htmlParsedTemplates *html.Template

func init() {
	os.Mkdir("html", os.ModePerm)

	tmpl, err := text.New("").ParseGlob("./templates/*")
	parsedTemplates = text.Must(tmpl, err)

	otherTmpl, err := html.New("").ParseGlob("./templates/*")
	htmlParsedTemplates = html.Must(otherTmpl, err)
}

type person struct {
	Name, LastName string
}

func (p person) GetHacked() string {
	return fmt.Sprintln("<script>alert(\"you were haxxed\")</script>")
}

func main() {
	p := person{
		Name:     "James",
		LastName: "Bond",
	}

	printToFile("./html/hackedTemplate.html", "index.gohtml", p)
	printSafeToFile("./html/safeTemplate.html", "index.gohtml", p)
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

func printSafeToFile(fn string, tn string, data interface{}) {
	f, err := os.Create(fn)
	if err != nil {
		log.Fatal("Failed to create a file", err)
	}
	defer f.Close()

	err = htmlParsedTemplates.ExecuteTemplate(f, tn, data)
	if err != nil {
		log.Fatal(err)
	}
}
