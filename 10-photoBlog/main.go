package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"photoBlog/internal/db"

	_ "github.com/go-sql-driver/mysql"
)

var parsedTemplates *template.Template

func init() {
	tmpl, err := template.New("").ParseGlob("./templates/*")
	parsedTemplates = template.Must(tmpl, err)
}

func main() {
	db.Connect()
	defer db.Close()

	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/signup", signupHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)

	fmt.Println("Server running...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func indexHandler(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	c, err := req.Cookie("session")
	if err != nil {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}

	err = db.CheckSession(c.Value)
	if err != nil {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}

	w.Write([]byte(getResponseBody("index.gohtml", nil)))
}

func signupHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		w.Write([]byte(getResponseBody("signup.gohtml", nil)))
	}

	if req.Method == http.MethodPost {
		un := req.FormValue("username")
		p := req.FormValue("password")

		_, err := db.CreateUser(db.User{
			UserName: un,
			Password: p,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		http.Redirect(w, req, "/login", http.StatusSeeOther)
	}
}

func loginHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		w.Write([]byte(getResponseBody("login.gohtml", nil)))
	}

	if req.Method == http.MethodPost {
		un := req.FormValue("username")
		p := req.FormValue("password")

		sID, err := db.LoginUser(db.User{
			UserName: un,
			Password: p,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		c := &http.Cookie{
			Name:  "session",
			Value: sID,
		}
		http.SetCookie(w, c)
		http.Redirect(w, req, "/", http.StatusSeeOther)
	}
}

func logoutHandler(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("logoutHandler"))
}

func getResponseBody(tn string, d interface{}) string {
	var tpl bytes.Buffer

	err := parsedTemplates.ExecuteTemplate(&tpl, tn, d)
	if err != nil {
		log.Fatal(err)
	}

	return tpl.String()
}
