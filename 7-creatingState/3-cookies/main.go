package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		c, err := req.Cookie("my-cookie")
		if err == http.ErrNoCookie {
			http.Redirect(w, req, "/foo", http.StatusSeeOther)
			return
		}

		count, _ := strconv.Atoi(c.Value)
		count++
		c.Value = strconv.Itoa(count)

		http.SetCookie(w, c)

		io.WriteString(w, c.Value)
	})

	http.Handle("/favicon.ico", http.NotFoundHandler())

	http.HandleFunc("/foo", func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("REQUEST - /foo")

		http.SetCookie(w, &http.Cookie{
			Name:   "my-cookie",
			Value:  "0",
			MaxAge: 5,
		})
	})

	http.HandleFunc("/bar", func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("REQUEST - /bar")

		c, err := req.Cookie("my-cookie")
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		fmt.Fprintln(w, "Your cookie: ", c.String())
	})

	http.HandleFunc("/foobar", func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("REQUEST - /foobar")

		c, err := req.Cookie("my-cookie")
		if err == http.ErrNoCookie {
			http.Redirect(w, req, "/bar", http.StatusSeeOther)
			return
		}

		c.MaxAge = -1

		http.SetCookie(w, c)
		fmt.Fprintln(w, "Your cookie was cleared")
	})

	http.ListenAndServe(":8080", nil)
}
