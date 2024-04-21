package handlers

import (
	"fmt"
	"log"
	"net/http"
	"net/mail"
	// "regexp"
	"text/template"

	"github.com/daltbunker/soul_climbers/db"
	"github.com/daltbunker/soul_climbers/types"
	"golang.org/x/crypto/bcrypt"
)

func HandleNewUser(w http.ResponseWriter, r *http.Request) {

	user := types.User{}
	user.Email = r.FormValue("email")
	user.Password = r.FormValue("password")
	user.Username = r.FormValue("username")

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user.Password = string(hash)

	dbUser, err := db.NewUser(r, user)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "<h2>Username: %s</h2><div>Email: %s</div>", dbUser.Username, dbUser.Email)
}

func HandleLoginUser(w http.ResponseWriter, r *http.Request) {
	type LoginForm struct {
		Email         string
		Password      string
		EmailError    string
		PasswordError string
	}
	loginValidation := LoginForm{}

	user := types.User{}
	user.Email = r.FormValue("email")
	user.Password = r.FormValue("password")

	dbUser, err := db.GetUserByEmail(r, user.Email)
	if err != nil {
		log.Printf("failed getting user from DB: %v", err)
		loginValidation.Email = user.Email
		loginValidation.EmailError = "email not found"
		renderComponent(w, "login.html", loginValidation)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
	if err != nil {
		log.Printf("failed checking password: %v", err)
		loginValidation.Email = user.Email
		loginValidation.PasswordError = "password is incorrect"
		renderComponent(w, "login.html", loginValidation)
		return
	}

	fmt.Fprintf(w, "<h2>Welcome, %s</h2>", dbUser.Username)
}

func HandleEmailValidate(w http.ResponseWriter, r *http.Request) {
	formInput := types.FormInput{}
	formInput.Name = "email"

	email := r.FormValue("email")
	formInput.Value = email

	_, err := mail.ParseAddress(email)
	if err != nil {
		formInput.Message = "please enter a valid email address"
		renderComponent(w, "input.html", formInput)
		return
	}

	_, err = db.GetUserByEmail(r, email)
	if err == nil {
		formInput.Message = "this email address is already taken"
		renderComponent(w, "input.html", formInput)
		return
	}

	renderComponent(w, "input.html", formInput)
}

func HandleUsernameValidate(w http.ResponseWriter, r *http.Request) {
	formInput := types.FormInput{}
	formInput.Name = "username"

	username := r.FormValue("username")
	formInput.Value = username

	_, err := db.GetUserByUsername(r, username)
	if err == nil {
		formInput.Message = "this username is already taken"
		renderComponent(w, "input.html", formInput)
		return
	}

	renderComponent(w, "input.html", formInput)
}

func HandlePasswordValidate(w http.ResponseWriter, r *http.Request) {

	// password := r.FormValue("password")
	// matched, err := regexp.Match("^(?=.*[A-Za-z])(?=.*\d)(?=.*[@$!%*#?&])[A-Za-z\d@$!%*#?&]{8,}$", []byte(password))

}

// todo: make this global / move to util
func renderComponent(w http.ResponseWriter, htmlFile string, data interface{}) {
	t, err := template.ParseFiles("templates/components/" + htmlFile)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = t.Execute(w, data)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
