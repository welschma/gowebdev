package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func addHandler(w http.ResponseWriter, r *http.Request) {
	a, err := strconv.Atoi(r.URL.Query().Get("a"))

	if err != nil {
		fmt.Fprintf(w, "Invalid value for a")
		return
	}

	b, err := strconv.Atoi(r.URL.Query().Get("b"))

	if err != nil {
		fmt.Fprintf(w, "Invalid value for b")
		return
	}

	fmt.Fprintf(w, "Add endpoint: %d + %d = %d", a, b, a+b)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "<h1>Welcome to my awesome site!</h1>")
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "<h1>Contact Page </h1><p>To get in touch, please send an email "+
		"to <a href=\"mailto:maxwelsch93@gmail.com\">maxwelsch93@gmail.com</a>.</p>")
}

func pageNotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	http.Error(w, "Page not found", http.StatusNotFound)
}

type Router struct {
}

func pathHandler (w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		homeHandler(w, r)
	case "/contact":
		contactHandler(w, r)
	case "/add":
		addHandler(w, r)
	default:
		pageNotFoundHandler(w, r)
	}
}

func (router Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pathHandler(w, r)
}

func main() {
	var router Router
	fmt.Println("Starting the server on :3000...")
	http.ListenAndServe(":3000", router)
}
