package handlers

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
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
			HandleServerError(w, r, err)
			return
		}
	}

	blogs, err := db.GetAllBlogs(r)
	if err != nil {
		HandleServerError(w, r, err)
		return
	}
	d := types.Home{Blogs: blogs}

	renderPage(pages["home"], w, r, d)
}

func HandleGetLogin(w http.ResponseWriter, r *http.Request) {
	if pages["login"] == nil {
		var err error
		pages["login"], err = template.ParseFiles(baseTemplate, "templates/pages/login.html", "templates/components/login.html")
		if err != nil {
			HandleServerError(w, r, err)
			return
		}
	}
	renderPage(pages["login"], w, r, types.LoginForm{})
}

func HandleGetSignup(w http.ResponseWriter, r *http.Request) {
	if pages["signup"] == nil {
		var err error
		pages["signup"], err = template.ParseFiles(baseTemplate, "templates/pages/signup.html", "templates/components/signup.html")
		if err != nil {
			HandleServerError(w, r, err)
			return
		}
	}
	renderPage(pages["signup"], w, r, types.SignupForm{})
}

func HandleGetAccount(w http.ResponseWriter, r *http.Request) {
	if pages["account"] == nil {
		var err error
		pages["account"], err = template.ParseFiles(baseTemplate, "templates/pages/account.html")
		if err != nil {
			HandleServerError(w, r, err)
			return
		}
	}
	user := types.User{} 
	session, err := GetSession(r)
	if err != nil {
		HandleServerError(w, r, err)
		return
	}
	if auth, ok := session.Values["authenticated"].(bool); ok && auth {
		user.Username = session.Values["username"].(string)
		user.Email = session.Values["email"].(string)
		user.Role = session.Values["role"].(string)
	}
	err = pages["account"].Execute(w, types.Base{MainContent: user, User: user})
	if err != nil {
		HandleServerError(w, r, err)
	}
}

func HandleGetResetEmail(w http.ResponseWriter, r *http.Request) {
	if pages["resetEmail"] == nil {
		var err error
		pages["resetEmail"], err = template.ParseFiles(baseTemplate, "templates/pages/reset-email.html")
		if err != nil {
			HandleServerError(w, r, err)
			return
		}
	}
	renderPage(pages["resetEmail"], w, r, nil)
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
			HandleServerError(w, r, err)
			return
		}
	}
	renderPage(pages["resetPassword"], w, r, nil)
}

func HandleGetBlog(w http.ResponseWriter, r *http.Request) {
	if pages["blog"] == nil {
		var err error
		pages["blog"], err = template.ParseFiles(baseTemplate, "templates/pages/blog.html")
		if err != nil {
			HandleServerError(w, r, err)
			return
		}
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Printf("Blog id must be type int: %v", err)
		HandleClientError(w, err)
		return
	}

	blog, err := db.GetBlogById(r, int32(id))
	if err != nil {
		HandleServerError(w, r, err)
		return
	}

	renderPage(pages["blog"], w, r, blog)
}

func HandleGetBlogForm(w http.ResponseWriter, r *http.Request) {
	if pages["blogForm"] == nil {
		var err error
		pages["blogForm"], err = template.ParseFiles(baseTemplate, "templates/pages/blog-form.html")
		if err != nil {
			HandleServerError(w, r, err)
			return
		}
	}

	paramId := chi.URLParam(r, "id")
	if paramId == ""{
		blogForm := types.BlogForm {
			FormAction: "/admin/blog/preview",
			FormMethod: "post",
		}
		renderPage(pages["blogForm"], w, r, blogForm)
		return
	}

	id, err := strconv.Atoi(paramId)	
	if err != nil {
		log.Printf("Blog id must be type int: %v", err)
		HandleClientError(w, err)
		return
	}

	blog, err := db.GetBlogById(r, int32(id))
	if err != nil {
		HandleServerError(w, r, err)
		return
	}	
	
	blogForm := types.BlogForm {
		Id: blog.Id,
		Body: blog.Body,
		Title: blog.Title,
		Excerpt: blog.Excerpt,
		ImgName: blog.ImgName, 
		FormAction: "/admin/blog/preview/" + paramId,
		FormMethod: "post",
	}

	renderPage(pages["blogForm"], w, r, blogForm)
}

func HandleGetBlogPreview(w http.ResponseWriter, r *http.Request) {

	if pages["blogPreview"] == nil {
		var err error
		pages["blogPreview"], err = template.ParseFiles(baseTemplate, "templates/pages/blog-preview.html")
		if err != nil {
			HandleServerError(w, r, err)
			return
		}
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Printf("Blog id must be type int: %v", err)
		HandleClientError(w, err)
		return
	}

	blog, err := db.GetBlogById(r, int32(id))
	if err != nil {
		HandleServerError(w, r, err)
		return
	}

	renderPage(pages["blogPreview"], w, r, blog)	
}

func HandleNewBlogPreview(w http.ResponseWriter, r *http.Request) {
	newBlog := types.Blog{}
	newBlog.Title = sanitize(r.FormValue("title"))
	newBlog.Body = sanitize(r.FormValue("body"))
	newBlog.Excerpt= sanitize(r.FormValue("excerpt"))
	newBlog.IsPublished = false
	newBlog.CreatedBy = "jonnyX"

	blog, err := db.NewBlog(r, newBlog)
	if err != nil {
		HandleClientError(w, err)
		return
	}

	thumbnail, header, err := r.FormFile("thumbnail") 
	if err != nil {
		HandleClientError(w, err)
		return
	}
	defer thumbnail.Close()

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, thumbnail); err != nil {
		HandleServerError(w, r, err)
		return
	}

	newBlogImg := types.BlogImg{
		ImgName: header.Filename,
		Img: buf.Bytes(),
		BlogId: blog.Id,
	} 

	_, err = db.NewBlogImg(r, newBlogImg)
	if err != nil {
		HandleServerError(w, r, err)
		return
	}

	if pages["blogPreview"] == nil {
		var err error
		pages["blogPreview"], err = template.ParseFiles(baseTemplate, "templates/pages/blog-preview.html")
		if err != nil {
			HandleServerError(w, r, err)
			return
		}
	}
	
	http.Redirect(w, r, "/admin/blog/preview/" + strconv.Itoa(int(blog.Id)), http.StatusSeeOther)
}

func HandleUpdateBlogPreview(w http.ResponseWriter, r *http.Request) {

	paramId := chi.URLParam(r, "id")
	id, err := strconv.Atoi(paramId)	

	updatedBlog := types.Blog{}
	updatedBlog.Id = int32(id)
	updatedBlog.Title = sanitize(r.FormValue("title"))
	updatedBlog.Body = sanitize(r.FormValue("body"))
	updatedBlog.Excerpt= sanitize(r.FormValue("excerpt"))
	updatedBlog.IsPublished = false

	blog, err := db.UpdateBlog(r, updatedBlog)
	if err != nil {
		HandleClientError(w, err)
		return
	}

	thumbnail, header, err := r.FormFile("thumbnail") 
	if err != nil && err != http.ErrMissingFile {
		HandleClientError(w, err)
		return
	}
	if err == nil {
		defer thumbnail.Close()

		buf := bytes.NewBuffer(nil)
		if _, err := io.Copy(buf, thumbnail); err != nil {
			HandleServerError(w, r, err)
			return
		}

		newBlogImg := types.BlogImg{
			ImgName: header.Filename,
			Img: buf.Bytes(),
			BlogId: blog.Id,
		} 

		_, err = db.NewBlogImg(r, newBlogImg)
		if err != nil {
			HandleServerError(w, r, err)
			return
		}
	}
	
	if pages["blogPreview"] == nil {
		var err error
		pages["blogPreview"], err = template.ParseFiles(baseTemplate, "templates/pages/blog-preview.html")
		if err != nil {
			HandleServerError(w, r, err)
			return
		}
	}

	http.Redirect(w, r, "/admin/blog/preview/" + strconv.Itoa(int(blog.Id)), http.StatusSeeOther)

}

func HandlePublishBlog(w http.ResponseWriter, r *http.Request) {
	paramId := chi.URLParam(r, "id")
	id, err := strconv.Atoi(paramId)	

	blog, err := db.GetBlogById(r, int32(id))
	if err != nil {
		HandleServerError(w, r, err)
		return
	}

    var missingFields []string
    if blog.ImgName == "" {
        missingFields = append(missingFields, "image")
    }
    if blog.Title == "" {
        missingFields = append(missingFields, "title")
    }
    if blog.Excerpt == "" {
        missingFields = append(missingFields, "excerpt")
    }
    if blog.Body == "" {
        missingFields = append(missingFields, "body")
    }

    if len(missingFields) > 0 {
        fmt.Fprintf(w, "<span class=\"sub-text warning\">*missing fields: %v</span>", strings.Join(missingFields, ", "))
        return
    }

	updatedBlog := blog
	updatedBlog.IsPublished = true

	blog, err = db.UpdateBlog(r, updatedBlog)
	if err != nil {
		HandleClientError(w, err)
		return
	}
	
	w.Header().Set("HX-Redirect", "/admin")
}

func HandleGetAdmin(w http.ResponseWriter, r *http.Request) {
	if pages["admin"] == nil {
		var err error
		pages["admin"], err = template.ParseFiles(baseTemplate, "templates/pages/admin.html", "templates/components/admin-blog-card.html")
		if err != nil {
			HandleServerError(w, r, err)
			return
		}
	}

    blogs, err := db.GetBlogsByCreator(r, 1) //TODO: don't hardcode userId 
    if err != nil {
        HandleServerError(w, r, err)
        return
    }

	d := types.Admin{Blogs: blogs}

	renderPage(pages["admin"], w, r, d)
}

func setCookie(w http.ResponseWriter, name string, value string) {
	cookie := &http.Cookie{
		Name:  name,
		Value: value,
		Path:  "/",
	}
	http.SetCookie(w, cookie)
}
