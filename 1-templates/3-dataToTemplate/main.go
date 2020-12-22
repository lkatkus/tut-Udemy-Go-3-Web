package main

import (
	"log"
	"os"
	"strings"
	"text/template"
)

var parsedTemplates *template.Template
var fm = template.FuncMap{
	"upperCase":  strings.ToUpper,
	"firstThree": firstThree,
}

func firstThree(s string) string {
	s = strings.TrimSpace(s)
	s = s[:3]
	return s
}

func init() {
	os.Mkdir("html", os.ModePerm)
	tmpl, err := template.New("").Funcs(fm).ParseGlob("./templates/*")
	parsedTemplates = template.Must(tmpl, err)
}

type person struct {
	Name     string
	lastName string
}

func main() {
	dataInt()
	dataSlice()
	dataMap()
	dataStruct()
	dataSliceStruct()
	dataSliceStructSlice()
	functions()
}

func dataInt() {
	printToFile("./html/int.html", "int.gohtml", 42)
}

func dataSlice() {
	sages := map[string]string{
		"Internet": "Musk",
		"Jedi":     "Emperor",
		"Destiny":  "Xur",
	}

	printToFile("./html/slice.html", "sages.gohtml", sages)
}

func dataMap() {
	sages := []string{"Ghandi", "Aristotel", "Plato"}

	printToFile("./html/map.html", "sages.gohtml", sages)
}

func dataStruct() {
	type person struct {
		Name     string
		LastName string
	}

	p := person{
		Name:     "James",
		LastName: "Bond",
	}

	printToFile("./html/struct.html", "struct.gohtml", p)
}

func dataSliceStruct() {
	type person struct {
		Name     string
		LastName string
	}

	people := []person{
		person{
			Name:     "James",
			LastName: "Bond",
		},
		person{
			Name:     "Miss",
			LastName: "Moneypenny",
		},
		person{
			Name:     "Dr.",
			LastName: "Evil",
		},
	}

	printToFile("./html/sliceStruct.html", "sliceStruct.gohtml", people)
}

func dataSliceStructSlice() {
	type person struct {
		Name      string
		LastName  string
		Favorites []string
	}

	people := []person{
		person{
			Name:      "James",
			LastName:  "Bond",
			Favorites: []string{"Moneypenny", "Martini", "Guns"},
		},
		person{
			Name:      "Miss",
			LastName:  "Moneypenny",
			Favorites: []string{"James", "Bombs"},
		},
		person{
			Name:      "Dr.",
			LastName:  "Evil",
			Favorites: []string{"Nothing"},
		},
	}

	printToFile("./html/sliceStructSlice.html", "sliceStructSlice.gohtml", people)
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

func functions() {
	type person struct {
		Name     string
		LastName string
	}

	p := person{
		Name:     "James",
		LastName: "Bond",
	}

	printToFile("./html/functions.html", "functions.gohtml", p)
}
