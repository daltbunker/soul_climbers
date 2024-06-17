package handlers

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
)

var pages map[string]*template.Template

func InitPages() {
	pages = make(map[string]*template.Template)
}

func GetPage(w http.ResponseWriter, name string) (*template.Template, error) {
	if pages[name] == nil {
		HandleServerError(w, fmt.Errorf("page '%s' is not defined", name))
	}
	return pages[name], nil
}

func renderPage(t *template.Template, w http.ResponseWriter, data interface{}) {
	err := t.Execute(w, data)
	if err != nil {
		HandleServerError(w, err)
	}
}

func renderComponent(w http.ResponseWriter, page string, name string, data interface{}) {
	t, err := GetPage(w, page)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	err = t.ExecuteTemplate(w, name, data)
	if err != nil {
		HandleServerError(w, err)
	}
}