package handlers

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/daltbunker/soul_climbers/types"
	"github.com/microcosm-cc/bluemonday"
)

var pages map[string]*template.Template

func InitPages() {
	pages = make(map[string]*template.Template)
}

func GetPage(w http.ResponseWriter, name string) (*template.Template, error) {
	if pages[name] == nil {
		return &template.Template{}, fmt.Errorf("page '%s' is not defined", name)
	}
	return pages[name], nil
}

func renderPage(t *template.Template, w http.ResponseWriter, r *http.Request, data interface{}) {
	user := types.User{} 
	session, err := GetSession(r)
	if err != nil {
		HandleServerError(w, r, err)
		return
	}
	if auth, ok := session.Values["authenticated"].(bool); ok && auth {
		user.Username = session.Values["email"].(string)
		user.Email = session.Values["username"].(string)
		user.Role = session.Values["role"].(string)
	}

	err = t.Execute(w, types.Base{MainContent: data, User: user})
	if err != nil {
		HandleServerError(w, r, err)
	}
}

func renderComponent(w http.ResponseWriter, r *http.Request, page string, name string, data interface{}) {
	t, err := GetPage(w, page)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	err = t.ExecuteTemplate(w, name, data)
	if err != nil {
		HandleServerError(w, r, err)
	}
}

func sanitize(s string) string {
    var p = bluemonday.UGCPolicy()
    return p.Sanitize(s)
}
