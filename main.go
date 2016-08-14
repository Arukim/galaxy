package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/arukim/galaxy/game"
)

var addr = flag.String("addr", ":8080", "http service address")

func homeHandler(c http.ResponseWriter, r *http.Request) {
	var indexTemplate = template.Must(template.ParseFiles("templates/index.html"))
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

	log.Print("server started")
	g := game.NewGame(1, 100, 25*time.Millisecond)
	g.Start()

	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
