package main

import (
	"io"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/file", fileHandler)   // Use function instead of Handler type
	http.HandleFunc("/other", otherHandler) // Use function instead of Handler type
	http.HandleFunc("/image", imageHandler) // Use function instead of Handler type

	http.ListenAndServe("localhost:8080", nil) // Use default ServerMux, when passing nil
}

func fileHandler(res http.ResponseWriter, req *http.Request) {
	f, err := os.Open("file.txt")
	defer f.Close()

	if err != nil {
		http.Error(res, "File not found", http.StatusNotFound)
		return
	}

	io.Copy(res, f)
}

func otherHandler(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, "other.txt")
}

func imageHandler(res http.ResponseWriter, req *http.Request) {
	f, err := os.Open("image.png")
	defer f.Close()

	if err != nil {
		http.Error(res, "File not found", 404)
		return
	}

	fi, err := f.Stat()
	if err != nil {
		http.Error(res, "File not found", 404)
	}

	http.ServeContent(res, req, f.Name(), fi.ModTime(), f)
}
