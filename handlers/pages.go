package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/daltbunker/soul_climbers/db"
	"github.com/daltbunker/soul_climbers/types"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

var baseTemplate = "templates/base.html" 

func HandleGetHome(w http.ResponseWriter, r *http.Request) {
	if pages["home"] == nil {
		var err error
		pages["home"], err = template.ParseFiles(baseTemplate, "templates/pages/home.html", "templates/components/blog-card.html")
		if err != nil {
			HandleServerError(w, err)
			return
		}
	}

	blogs, err := db.GetAllBlogs(r)
	if err != nil {
		HandleServerError(w, err)
	}
	d := types.Home{Blogs: blogs}

	renderPage(pages["home"], w, d)
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

func HandleGetBlog(w http.ResponseWriter, r *http.Request) {
	if pages["blog"] == nil {
		var err error
		pages["blog"], err = template.ParseFiles(baseTemplate, "templates/pages/blog.html")
		if err != nil {
			HandleServerError(w, err)
			return
		}
	}
	
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Printf("Blog id must be type int: %v", err)
		HandleClientError(w, err)
	}

	blog, err := db.GetBlogById(r, int32(id))
	if err != nil {
		HandleServerError(w, err)
	}

	renderPage(pages["blog"], w, blog)
}

func HandleGetBlogForm(w http.ResponseWriter, r *http.Request) {
	if pages["blogForm"] == nil {
		var err error
		pages["blogForm"], err = template.ParseFiles(baseTemplate, "templates/pages/blog-form.html")
		if err != nil {
			HandleServerError(w, err)
			return
		}
	}

	renderPage(pages["blogForm"], w, nil)
}

func setCookie(w http.ResponseWriter, name string, value string) {
	cookie := &http.Cookie{
		Name:  name,
		Value: value,
		Path:  "/",
	}
	http.SetCookie(w, cookie)
}
