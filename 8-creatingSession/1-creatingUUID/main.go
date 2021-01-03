package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"text/template"

	uuid "github.com/satori/go.uuid"
)

type user struct {
	UserName  string
	FirstName string
	LastName  string
}

var dbUsers = map[string]user{}
var dbSessions = map[string]string{}
var parsedTemplates *template.Template

func init() {
	tmpl, err := template.New("").ParseGlob("./templates/*")
	parsedTemplates = template.Must(tmpl, err)
}

func main() {
	http.HandleFunc("/foo", fooHandler)
	http.HandleFunc("/bar", barHandler)
	http.HandleFunc("/foobar", foobarHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func fooHandler(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte(getResponseBody("index.gohtml", nil)))
}

func barHandler(w http.ResponseWriter, req *http.Request) {
	c, err := req.Cookie("session")
	if err != nil {
		sID, err := uuid.NewV4()
		if err != nil {
			log.Fatal("Unabled to create uuid", err)
		}

		c = &http.Cookie{
			Name:  "session",
			Value: sID.String(),
		}

		http.SetCookie(w, c)
	}

	var u user

	if un, ok := dbSessions[c.Value]; ok {
		u = dbUsers[un]
	}

	if req.Method == http.MethodPost {
		un := req.FormValue("username")
		f := req.FormValue("firstname")
		l := req.FormValue("lastname")

		u = user{
			UserName:  un,
			FirstName: f,
			LastName:  l,
		}

		dbSessions[c.Value] = un
		dbUsers[un] = u

		http.Redirect(w, req, "/foobar", http.StatusSeeOther)
	}
}

func foobarHandler(w http.ResponseWriter, req *http.Request) {
	c, err := req.Cookie("session")
	if err != nil {
		fmt.Println("ERROR req.Cookie")

		http.Redirect(w, req, "/foo", http.StatusSeeOther)
		return
	}

	un, ok := dbSessions[c.Value]
	if !ok {
		fmt.Println("ERROR dbSessions[c.Value]")

		http.Redirect(w, req, "/foo", http.StatusSeeOther)
		return
	}

	u := dbUsers[un]
	w.Write([]byte(getResponseBody("user.gohtml", u)))
}

func getResponseBody(tn string, d interface{}) string {
	var tpl bytes.Buffer

	err := parsedTemplates.ExecuteTemplate(&tpl, tn, d)
	if err != nil {
		log.Fatal(err)
	}

	return tpl.String()
}
