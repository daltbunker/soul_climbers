package handlers

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/daltbunker/soul_climbers/types"
	"github.com/gorilla/sessions"
	"github.com/microcosm-cc/bluemonday"
)

var store *sessions.CookieStore
var pages map[string]*template.Template

func InitSession(key string) {
	store = sessions.NewCookieStore([]byte(key))
}

func GetSession(r *http.Request) (*sessions.Session, error) {
	return store.Get(r, "session")
}

func NewSession(r *http.Request, w http.ResponseWriter, user types.User) error {
	session, err := store.Get(r, "session")
	if err != nil {
		return err
	}
	session.Values["authenticated"] = true
	session.Values["email"] = user.Email
	session.Values["username"] = user.Username
	session.Values["role"] = user.Role
	session.Options.MaxAge = 30
	err = session.Save(r, w)
	if err != nil {
		return err
	}
	return nil
}

func NewDevSession(w http.ResponseWriter, r *http.Request) {
	NewSession(r, w, types.User{Username: "Wayne", Email: "wayne_ker@aol.com", Role: "admin"})	
}

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
