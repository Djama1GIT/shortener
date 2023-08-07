package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func shortener(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	} else {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()
		// Process the request body and generate the shortened URL
		shortURL := generateShortURL(string(body))

		// Return the shortened URL as the response
		fmt.Fprintf(w, shortURL)
	}
}

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
	http.HandleFunc("/agreement/", agreement)
	http.HandleFunc("/shortener/", shortener)
	http.ListenAndServe(":8081", nil)
}

func generateShortURL(longURL string) string {
	// Generate the shortened URL logic goes here
	return longURL
}
