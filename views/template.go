package views

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"

	"github.com/gorilla/csrf"
)

type Template struct {
	htmlTpl *template.Template
}

func (t Template) Execute(w http.ResponseWriter, r *http.Request, data any) {

    tpl, err := t.htmlTpl.Clone()

    if err !=nil {
        log.Printf("cloning template: %v", err)
        http.Error(w, "There was an error rendering the page.", http.StatusInternalServerError)
        return
    }

    tpl = tpl.Funcs(
        template.FuncMap{
            "csrfField": func() (template.HTML, error) {
                return csrf.TemplateField(r), nil
            },
        },
        )

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

    var buf bytes.Buffer

	err = tpl.Execute(&buf, data)

	if err != nil {
		log.Printf("error executing template: %v", err)
		http.Error(w, "There was an error executing the template", http.StatusInternalServerError)
		return
	}

    io.Copy(w, &buf)
}

func ParseFS(fs fs.FS, patterns ...string) (Template, error) {
    tpl := template.New(patterns[0])

    tpl = tpl.Funcs(
        template.FuncMap{
            "csrfField": func() (template.HTML, error) {
                return "", fmt.Errorf("csrfField not implemented")
            },
        },
        )

	tpl, err := tpl.ParseFS(fs, patterns...)

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
