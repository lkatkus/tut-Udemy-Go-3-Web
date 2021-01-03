package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

type user struct {
	UserName  string
	Password  []byte
	FirstName string
	LastName  string
	Role      string
}

type session struct {
	un           string
	lastActivity time.Time
}

var dbUsers = map[string]user{}
var dbSessions = map[string]session{}
var dbSessionsClean time.Time
var parsedTemplates *template.Template

func init() {
	tmpl, err := template.New("").ParseGlob("./templates/*")
	parsedTemplates = template.Must(tmpl, err)
	dbSessionsClean = time.Now()
}

func main() {
	ticker := time.NewTicker(5 * time.Second)

	go checkSessions(ticker)

	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.HandleFunc("/", userHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/signup", signupHandler)
	http.HandleFunc("/agent", agentHandler)
	http.ListenAndServe(":8080", nil)
}

func checkSessions(ticker *time.Ticker) {
	for {
		select {
		case <-ticker.C:
			fmt.Println("Cleaning sessions")

			for k, s := range dbSessions {
				if time.Now().Sub(s.lastActivity) > (time.Second * 10) {
					delete(dbSessions, k)
				}
			}
			// case <-quit:
			// 	ticker.Stop()
			// 	return
			// }
		}
	}
}

func userHandler(w http.ResponseWriter, req *http.Request) {
	if alreadyLoggedIn(req) {
		c, err := req.Cookie("session")
		if err != nil {
			http.Redirect(w, req, "/login", http.StatusSeeOther)
			return
		}

		sn, ok := dbSessions[c.Value]
		if !ok {
			http.Redirect(w, req, "/login", http.StatusSeeOther)
			return
		}

		u, ok := dbUsers[sn.un]
		if ok {
			dbSessions[c.Value] = session{
				un:           sn.un,
				lastActivity: time.Now(),
			}

			w.Write([]byte(getResponseBody("user.gohtml", u)))
			return
		}
	}

	http.Redirect(w, req, "/login", http.StatusSeeOther)
}

func agentHandler(w http.ResponseWriter, req *http.Request) {
	if alreadyLoggedIn(req) {
		c, err := req.Cookie("session")
		if err != nil {
			http.Redirect(w, req, "/login", http.StatusSeeOther)
			return
		}

		sn, ok := dbSessions[c.Value]
		if !ok {
			http.Redirect(w, req, "/login", http.StatusSeeOther)
			return
		}

		u, ok := dbUsers[sn.un]
		if ok {
			dbSessions[c.Value] = session{
				un:           sn.un,
				lastActivity: time.Now(),
			}

			if u.Role == "agent" {
				w.Write([]byte("Hello secret agent"))
				return
			} else {
				http.Error(w, "This is not the page you are looking for", http.StatusForbidden)
				return
			}
		}
	}

	http.Redirect(w, req, "/login", http.StatusSeeOther)
}

func loginHandler(w http.ResponseWriter, req *http.Request) {
	if alreadyLoggedIn(req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	if req.Method == http.MethodGet {
		w.Write([]byte(getResponseBody("login.gohtml", nil)))
	}

	if req.Method == http.MethodPost {
		un := req.FormValue("username")
		p := req.FormValue("password")

		u, ok := dbUsers[un]
		if !ok {
			http.Error(w, "Wrong username of password", http.StatusForbidden)
			return
		}

		err := bcrypt.CompareHashAndPassword(u.Password, []byte(p))
		if err != nil {
			http.Error(w, "Wrong username of password", http.StatusForbidden)
			return
		}

		sID, err := uuid.NewV4()
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}

		c := &http.Cookie{
			Name:  "session",
			Value: sID.String(),
		}
		http.SetCookie(w, c)
		dbSessions[c.Value] = session{un, time.Now()}
		http.Redirect(w, req, "/", http.StatusSeeOther)
	}
}

func logoutHandler(w http.ResponseWriter, req *http.Request) {
	if !alreadyLoggedIn(req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	c, _ := req.Cookie("session")

	delete(dbSessions, c.Value)

	c = &http.Cookie{
		Name:   "session",
		Value:  "",
		MaxAge: -1,
	}

	http.SetCookie(w, c)
	http.Redirect(w, req, "/", http.StatusSeeOther)
}

func signupHandler(w http.ResponseWriter, req *http.Request) {
	if alreadyLoggedIn(req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	if req.Method == http.MethodGet {
		w.Write([]byte(getResponseBody("signup.gohtml", nil)))
	}

	if req.Method == http.MethodPost {
		un := req.FormValue("username")
		p := req.FormValue("password")
		fn := req.FormValue("firstname")
		ln := req.FormValue("lastname")
		r := "civilian"

		if ln == "bond" {
			r = "agent"
		}

		if _, ok := dbUsers[un]; ok {
			http.Error(w, "Username already taken", http.StatusForbidden)
		}

		sID, err := uuid.NewV4()
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}

		c := &http.Cookie{
			Name:  "session",
			Value: sID.String(),
		}
		http.SetCookie(w, c)
		dbSessions[c.Value] = session{un, time.Now()}

		ep, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.MinCost)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}

		u := user{
			UserName:  un,
			Password:  ep,
			FirstName: fn,
			LastName:  ln,
			Role:      r,
		}
		dbUsers[un] = u

		http.Redirect(w, req, "/", http.StatusSeeOther)
	}

}

func alreadyLoggedIn(req *http.Request) bool {
	c, err := req.Cookie("session")
	if err != nil {
		return false
	}

	s := dbSessions[c.Value]
	_, ok := dbUsers[s.un]

	return ok
}

func getResponseBody(tn string, d interface{}) string {
	var tpl bytes.Buffer

	err := parsedTemplates.ExecuteTemplate(&tpl, tn, d)
	if err != nil {
		log.Fatal(err)
	}

	return tpl.String()
}
