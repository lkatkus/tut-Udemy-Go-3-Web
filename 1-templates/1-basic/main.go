package main

import (
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	name := "Agent Smith"
	template := `
		<div>
			<div>Hello</div>
			<div>
				<h1>` + name + `</h1>
				<p>Matrix is a lie</p>
			</div>
		</div>
		`

	f, err := os.Create("index.html")
	if err != nil {
		log.Fatal("Error creating a file", err)
	}
	defer f.Close()

	io.Copy(f, strings.NewReader(template))

	/* or pipe output to a file with - go run main.go > index.html */
}
