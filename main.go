package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func init() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

func shortener(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	} else {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}
		if string(body) == "" {
			http.Error(w, "Invalid link", http.StatusTeapot)
		}
		defer r.Body.Close()
		// Process the request body and generate the shortened URL
		shortURL := generateShortURL(string(body))

		// Return the shortened URL as the response
		fmt.Fprintf(w, shortURL)
	}
}

func home_page(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		ctx := context.Background()
		shortedURL := r.URL.Path[1:]

		longURL, err := redisClient.Get(ctx, shortedURL).Result()
		if err != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			fmt.Println(shortedURL, longURL)
			return
		}

		http.Redirect(w, r, longURL, http.StatusSeeOther)
		return
	}

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

	go func() {
		if err := http.ListenAndServe(":8083", nil); err != nil {
			panic(err)
		}
	}()

	<-make(chan struct{})
}

func generateShortURL(longURL string) string {
	ctx := context.Background()

	shortedURL := shortenerURL(longURL)
	err := redisClient.Set(ctx, shortedURL, longURL, 0).Err()
	if err != nil {
		return ""
	}

	return "http://127.0.0.1:8083/" + shortedURL
}

func shortenerURL(longURL string) string {
	alphabet := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	shortedURL := make([]rune, 5)

	for i := range shortedURL {
		shortedURL[i] = alphabet[rand.Intn(len(alphabet))]
	}

	return string(shortedURL)
}
