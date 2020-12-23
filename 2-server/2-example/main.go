package main

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net"
	"strings"
)

var parsedTemplates *template.Template

func init() {
	tmpl, err := template.New("").ParseGlob("./templates/*")
	parsedTemplates = template.Must(tmpl, err)
}

func main() {
	li, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer li.Close()

	for {
		conn, err := li.Accept()
		if err != nil {
			log.Fatalln(err.Error())
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	handleRequest(conn)
}

func handleRequest(conn net.Conn) {
	i := 0
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		ln := scanner.Text()
		fmt.Println(ln)

		if i == 0 {
			m := strings.Fields(ln)[1]

			handleResponse(conn, m)
		}
		if ln == "" {
			break
		}
		i++
	}
}

func handleResponse(conn net.Conn, url string) {
	body := getResponseBody(url)

	fmt.Fprint(conn, "HTTP/1.1 200 OK\r\n")
	fmt.Fprintf(conn, "Content-Length: %d\r\n", len(body))
	fmt.Fprint(conn, "Content-Type: text/html\r\n")
	fmt.Fprint(conn, "\r\n")
	fmt.Fprint(conn, body)
}

func getResponseBody(url string) string {
	var tpl bytes.Buffer

	err := parsedTemplates.ExecuteTemplate(&tpl, "index.gohtml", url)
	if err != nil {
		log.Fatal(err)
	}

	return tpl.String()
}
