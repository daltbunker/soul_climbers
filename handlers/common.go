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

func GetSessionUser(r *http.Request) (types.User, error) {
	session, err := store.Get(r, "session")
	if err != nil {
		return types.User{}, err
	}
	// User is not logged in
	if session.IsNew {
		return types.User{}, nil
	}
	user := types.User{
		Email: session.Values["email"].(string),
		Username: session.Values["username"].(string),
		Role: session.Values["role"].(string),
		SoulScore: session.Values["soul_score"].(int32),
	}

	return user, nil
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
	session.Values["soul_score"] = user.SoulScore
	session.Options.MaxAge = 30
	err = session.Save(r, w)
	if err != nil {
		return err
	}
	return nil
}

// TODO: Change to handle strings and validate session value
func UpdateSession(r *http.Request, w http.ResponseWriter, key string, value int32) error {
	session, err := store.Get(r, "session")
	if err != nil {
		return err
	}

	session.Values[key] = value 
	err = session.Save(r, w)
	if err != nil {
		return err
	}
	return nil
}

func NewDevSession(w http.ResponseWriter, r *http.Request) {
	NewSession(r, w, types.User{Username: "garth", Email: "garth@aol.com", Role: "admin", SoulScore: 0})	
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
	user, err := GetSessionUser(r)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, types.Base{MainContent: data, User: user})
	if err != nil {
		log.Print(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func renderComponent(w http.ResponseWriter, page string, name string, data interface{}) {
	t, err := GetPage(w, page)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = t.ExecuteTemplate(w, name, data)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func sanitize(s string) string {
    var p = bluemonday.UGCPolicy()
    return p.Sanitize(s)
}
