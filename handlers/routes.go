package handlers

import (
	"log"
	"net/http"
	"text/template"
)

func HandleGetHome(w http.ResponseWriter, r *http.Request) {
	renderPage(w, "home.html")
}

func HandleGetLogin(w http.ResponseWriter, r *http.Request) {
	renderPage(w, "login.html")
}

func HandleGetSignup(w http.ResponseWriter, r *http.Request) {
	renderPage(w, "signup.html")
}

func renderPage(w http.ResponseWriter, htmlFile string) {
	t, err := template.ParseFiles("templates/pages/base.html", "templates/pages/" + htmlFile)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = t.Execute(w, nil)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}