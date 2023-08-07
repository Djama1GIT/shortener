package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strings"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client
var HOST = "http://ama1.ru"

func init() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

func shortener(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

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
			return
		}
		defer r.Body.Close()
		// Process the request body and generate the shortened URL
		shortURL := generateShortURL(string(strings.Trim(string(body), " ")))

		// Return the shortened URL as the response
		fmt.Fprintf(w, shortURL)
	}
}

func home_page(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.URL.Path != "/" {
		ctx := context.Background()
		shortedURL := r.URL.Path[1:]

		longURL, err := redisClient.Get(ctx, shortedURL).Result()
		if err != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		if !strings.HasPrefix(longURL, "http://") && !strings.HasPrefix(longURL, "https://") {
			longURL = "http://" + longURL
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
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", home_page)
	http.HandleFunc("/agreement/", agreement)
	http.HandleFunc("/shortener/", shortener)

	go func() {
		if err := http.ListenAndServe("212.109.218.42:80", nil); err != nil {
			log.Fatal("HTTP server error: ", err)
		}
	}()

	log.Println("Server started on port 80")

	<-make(chan struct{})
}

func generateShortURL(longURL string) string {
	ctx := context.Background()
	if validateURL(longURL) != "" {
		for {
			shortedURL := shortURL()
			_, err := redisClient.Get(ctx, shortedURL).Result()
			if err != nil {
				err := redisClient.Set(ctx, shortedURL, longURL, 0).Err()
				if err != nil {
					log.Println("Error saving URL to Redis:", err)
					return ""
				}
				return HOST + "/" + shortedURL
			}
		}
	}
	return ""
}

func validateURL(URL string) string {
	pattern := `^((http|https)://)?[a-zA-Z0-9\-\.]+\.[a-zA-Z]{2,}(\/\S*)?$`
	match, _ := regexp.MatchString(pattern, URL)
	if match {
		return URL
	} else {
		return ""
	}
}

func shortURL() string {
	alphabet := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	shortedURL := make([]rune, 5)

	for i := range shortedURL {
		shortedURL[i] = alphabet[rand.Intn(len(alphabet))]
	}

	return string(shortedURL)
}
