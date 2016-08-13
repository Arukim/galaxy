package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8080", "http service address")
var indexTemplate = template.Must(template.ParseFiles("templates/index.html"))

func homeHandler(c http.ResponseWriter, r *http.Request) {
	data := struct {
		Host string
	}{r.Host}
	indexTemplate.Execute(c, data)
}

func main() {
	http.HandleFunc("/", homeHandler)

	http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})

	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
