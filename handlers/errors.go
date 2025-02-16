package handlers

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
)

func HandleServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("internal server error: %v", err)

	if pages["serverError"] == nil {
		var err error
		pages["serverError"], err = template.ParseFS(templates, baseTemplate, "templates/pages/server-error.html")
		if err != nil {
			log.Print(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// If HTMX request, need to render component not entire page
	if r.Header.Get("HX-Request") != "" {
		w.Header().Set("HX-Retarget", "#toast-container")
		fmt.Fprint(w, "<span class=\"error-toast\">Ahh nuttz 🥜, there was an error on our server. If it's not fixed in 5 minutes, just wait longer</span>")
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
	renderPage(pages["serverError"], w, r, nil)
}

func HandleUnautorized(w http.ResponseWriter, r *http.Request, authenticated bool) {
	log.Printf("user unauthorized, authenticated=%v", authenticated)

	if pages["unauthorized"] == nil {
		var err error
		pages["unauthorized"], err = template.ParseFS(templates, baseTemplate, "templates/pages/unauthorized.html")
		if err != nil {
			log.Print(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// If HTMX request, need to render component not entire page
	if r.Header.Get("HX-Request") != "" {
		w.Header().Set("HX-Retarget", ".main-content")
		renderComponent(w, "unauthorized", "content", authenticated)
		return
	}

	w.WriteHeader(http.StatusForbidden)
	renderPage(pages["unauthorized"], w, r, authenticated)
}

func HandleClientError(w http.ResponseWriter, err error) {
	w.Header().Set("HX-Retarget", "#toast-container")
	log.Printf("client error: %v", err)
	fmt.Fprintf(w, "<span class=\"error-toast\">Error(400) %v</span>", err.Error())
}

func HandleNotFound(w http.ResponseWriter, r *http.Request) {

	if pages["not-found"] == nil {
		var err error
		pages["not-found"], err = template.ParseFS(templates, baseTemplate, "templates/pages/not-found.html")
		if err != nil {
			log.Print(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// If HTMX request, need to render component not entire page
	if r.Header.Get("HX-Request") != "" {
		w.Header().Set("HX-Retarget", ".main-content")
		renderComponent(w, "not-found", "content", nil)
		return
	}

	w.WriteHeader(http.StatusNotFound)
	renderPage(pages["not-found"], w, r, nil)
}
