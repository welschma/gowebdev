package views

import (
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
)

type Template struct {
	htmlTpl *template.Template
}

func (t Template) Execute(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	err := t.htmlTpl.Execute(w, data)

	if err != nil {
		log.Printf("error executing template: %v", err)
		http.Error(w, "There was an error executing the template", http.StatusInternalServerError)
		return
	}
}

func Parse(filepath string) (Template, error) {
	tpl, err := template.ParseFiles(filepath)

	if err != nil {
		return Template{}, fmt.Errorf("error parsing template: %w", err)
	}

	return Template{htmlTpl: tpl}, nil
}

func ParseFS(fs fs.FS, patterns ...string) (Template, error) {
	tpl, err := template.ParseFS(fs, patterns...)

	if err != nil {
		return Template{}, fmt.Errorf("error parsing template: %w", err)
	}

	return Template{htmlTpl: tpl}, nil
}

func Must(t Template, err error) Template {
	if err != nil {
		panic(err)
	}
	return t
}
