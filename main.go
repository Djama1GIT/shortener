package main

import (
	"net/http"
)

func home_page(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/index.html")
}

func agreement(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/agreement.html")
}

func main() {
	fs := http.FileServer(http.Dir("templates/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", home_page)
	http.HandleFunc("/agreement", agreement)
	http.ListenAndServe(":8081", nil)
}
