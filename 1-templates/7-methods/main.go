package main

import (
	"fmt"
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

type person struct {
	Name, LastName string
}

func (p person) GetFullName() string {
	return fmt.Sprintln("My full name is", p.Name, p.LastName)
}

func (p person) GetAge(age int) string {
	return fmt.Sprintln(p.Name, p.LastName, "is", age, "old")
}

func main() {
	p := person{
		Name:     "James",
		LastName: "Bond",
	}

	printToFile("./html/methods.html", "index.gohtml", p)
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
