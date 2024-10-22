package handlers

import (
	"log"
	"net/http"
	"text/template"
)

func HandleServerError(w http.ResponseWriter, err error) {
	if pages["serverError"] == nil {
		var err error
		pages["serverError"], err = template.ParseFiles(baseTemplate, "templates/pages/server-error.html")
		if err != nil {
			log.Print(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	// FIXME: returns error message and html, html doesn't render
	w.WriteHeader(http.StatusInternalServerError)
	log.Printf("internal server error: %v", err)
	renderPage(pages["serverError"], w, nil)
}

func HandleUnautorized(w http.ResponseWriter, err error) {
	if pages["unauthorized"] == nil {
		var err error
		pages["unauthorized"], err = template.ParseFiles(baseTemplate, "templates/pages/unauthorized.html")
		if err != nil {
			log.Print(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusForbidden)
	log.Printf("user unauthorized")
	renderPage(pages["unauthorized"], w, nil)
}

func HandleClientError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	log.Printf("bad request: %v", err)
	http.Error(w, err.Error(), http.StatusBadRequest)
}
