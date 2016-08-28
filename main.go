package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"

	"github.com/arukim/galaxy/core"
	"github.com/arukim/galaxy/login"
)

var addr = flag.String("addr", ":8080", "http service address")

var indexTemplate = template.Must(template.ParseFiles("templates/index.html",
	"templates/_serverStats.html"))

func homeHandler(c http.ResponseWriter, r *http.Request) {
	data := struct {
		Stats *core.StatisticsInfo
		Host  string
	}{core.GameStats.ToInfo(), r.Host}
	indexTemplate.Execute(c, data)
}

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static", fs)

	http.HandleFunc("/", homeHandler)

	http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})

	gs := login.NewServer("/galaxy")
	gs.Listen()

	log.Print("server started")
	//g := game.NewGame(1, 100, 25*time.Millisecond)

	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
