package handlers

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/daltbunker/soul_climbers/db"
	"github.com/daltbunker/soul_climbers/types"
	"github.com/google/uuid"
)

var baseTemplate = "templates/base.html" 

func HandleGetHome(w http.ResponseWriter, r *http.Request) {
	if pages["home"] == nil {
		var err error
		pages["home"], err = template.ParseFiles(baseTemplate, "templates/pages/home.html")
		if err != nil {
			HandleServerError(w, err)
			return
		}
	}
	renderPage(pages["home"], w, nil)
}

func HandleGetLogin(w http.ResponseWriter, r *http.Request) {
	if pages["login"] == nil {
		var err error
		pages["login"], err = template.ParseFiles(baseTemplate, "templates/pages/login.html", "templates/components/login.html")
		if err != nil {
			HandleServerError(w, err)
			return
		}
	}
	renderPage(pages["login"], w, types.LoginForm{})
}

func HandleGetSignup(w http.ResponseWriter, r *http.Request) {
	if pages["signup"] == nil {
		var err error
		pages["signup"], err = template.ParseFiles(baseTemplate, "templates/pages/signup.html", "templates/components/signup.html")
		if err != nil {
			HandleServerError(w, err)
		}
	}
	renderPage(pages["signup"], w, types.SignupForm{})
}

func HandleGetResetEmail(w http.ResponseWriter, r *http.Request) {
	if pages["resetEmail"] == nil {
		var err error
		pages["resetEmail"], err = template.ParseFiles(baseTemplate, "templates/pages/reset-email.html")
		if err != nil {
			HandleServerError(w, err)
		}
	}
	renderPage(pages["resetEmail"], w, nil)
}

func HandleGetResetPassword(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	parsedToken, err := uuid.Parse(token)
	if err != nil {
		log.Printf("Failed to get parse token %v", err)
		fmt.Fprint(w, "Invalid token")
		return
	}
	requestToken, err := db.GetResetTokenByToken(r, parsedToken)
	if err != nil {
		log.Printf("Failed to get reset token %v", err)
		fmt.Fprint(w, "Invalid token")
		return
	}

	if time.Now().After(requestToken.Expiration) {
		log.Printf("Failed to get reset token %v", err)
		fmt.Fprint(w, "Token is expired")
		return
	}

	setCookie(w, "Reset-Token", token)

	if pages["resetPassword"] == nil {
		var err error
		pages["resetPassword"], err = template.ParseFiles(baseTemplate, "templates/pages/reset-password.html")
		if err != nil {
			HandleServerError(w, err)
		}
	}
	renderPage(pages["resetPassword"], w, nil)
}

func HandleGetAdmin(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Admin Page</h1>")
}

func setCookie(w http.ResponseWriter, name string, value string) {
	cookie := &http.Cookie{
		Name:  name,
		Value: value,
		Path:  "/",
	}
	http.SetCookie(w, cookie)
}
