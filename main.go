package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/welschma/gowebdev/controllers"
	"github.com/welschma/gowebdev/views"
)

func addHandler(w http.ResponseWriter, r *http.Request) {
	a, err := strconv.Atoi(chi.URLParam(r, "a"))

	if err != nil {
		fmt.Fprintf(w, "Invalid value for a")
		return
	}

	b, err := strconv.Atoi(chi.URLParam(r, "b"))

	if err != nil {
		fmt.Fprintf(w, "Invalid value for b")
		return
	}

	log.Println("URL Parameters: a =", a, ", b =", b)
	fmt.Fprintf(w, "Add endpoint: %d + %d = %d", a, b, a+b)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tplPath := filepath.Join("templates", "home.gohtml")
	executeTemplate(w, tplPath)
}

func executeTemplate(w http.ResponseWriter, filepath string) {
	tpl, err := views.Parse(filepath)

	if err != nil {
		log.Printf("error parsing template: %v", err)
		http.Error(w, "There was an error executing the template", http.StatusInternalServerError)
	}

	tpl.Execute(w, nil)
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	tplPath := filepath.Join("templates", "contact.gohtml")
	executeTemplate(w, tplPath)
}

func pageNotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	http.Error(w, "Page not found", http.StatusNotFound)
}

func faqHandler(w http.ResponseWriter, r *http.Request) {
	executeTemplate(w, "templates/faq.gohtml")
}

func main() {
	r := chi.NewRouter()

	tpl, err := views.Parse(filepath.Join("templates", "home.gohtml"))
	if err != nil {
		panic(err)
	}
	r.Get("/", controllers.StaticHandler(tpl))

	tpl, err = views.Parse(filepath.Join("templates", "contact.gohtml"))
	if err != nil {
		panic(err)
	}
	r.Get("/contact", controllers.StaticHandler(tpl))

	tpl, err = views.Parse(filepath.Join("templates", "faq.gohtml"))
	if err != nil {
		panic(err)
	}
	r.Get("/faq", controllers.StaticHandler(tpl))

	fmt.Println("Starting the server on :3000...")
	http.ListenAndServe(":3000", r)
}
