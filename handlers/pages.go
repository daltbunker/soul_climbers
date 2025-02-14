package handlers

import (
	"bytes"
	"database/sql"
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
	user, err := GetSessionUser(r)
	if err != nil {
		HandleServerError(w, r, err)
		return
	}
	renderPage(pages["account"], w, r, user)
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
		HandleNotFound(w, r)
		return
	}

	blog, err := db.GetBlogById(r, int32(id))
	if err == sql.ErrNoRows {
		log.Printf("No data found for blog with id: %v", id)
		HandleNotFound(w, r)
		return
	}
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
			RequestURL: "/admin/blog/preview",
		}
		renderPage(pages["blogForm"], w, r, blogForm)
		return
	}

	id, err := strconv.Atoi(paramId)	
	if err != nil {
		log.Printf("Blog id must be type int: %v", err)
		HandleNotFound(w, r)
		return
	}

	blog, err := db.GetBlogById(r, int32(id))
	if err == sql.ErrNoRows {
		HandleNotFound(w, r)
		return
	}
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
		RequestURL: "/admin/blog/preview/" + paramId,
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
		HandleNotFound(w, r)
		return
	}

	blog, err := db.GetBlogById(r, int32(id))
	if err == sql.ErrNoRows {
		HandleNotFound(w, r)
		return
	}
	if err != nil {
		HandleServerError(w, r, err)
		return
	}

	renderPage(pages["blogPreview"], w, r, blog)	
}

func HandleNewBlogPreview(w http.ResponseWriter, r *http.Request) {
	sessionUser, err := GetSessionUser(r)
	if err != nil {
		HandleServerError(w, r, err)
		return
	}

	newBlog := types.Blog{}
	newBlog.Title = sanitize(r.FormValue("title"))
	newBlog.Body = sanitize(r.FormValue("body"))
	newBlog.Excerpt= sanitize(r.FormValue("excerpt"))
	newBlog.IsPublished = false
	newBlog.CreatedBy = sessionUser.Username

	hasUniqueTitle, err := hasUniqueTitle(r, newBlog.Title, -1)
	if err != nil {
		HandleServerError(w, r, err)
		return
	}
	if !hasUniqueTitle {
		HandleClientError(w, fmt.Errorf("title '%v' is already taken", newBlog.Title))
		return
	}

	// File(and all other parameters) are optional for preview
	hasFile := true 
	thumbnail, header, err := r.FormFile("thumbnail") 
	if err == http.ErrMissingFile {
		hasFile = false 
	} else if err != nil {
		HandleClientError(w, err)
		return
	}

	blog, err := db.NewBlog(r, newBlog)
	if err != nil {
		HandleServerError(w, r, err)
		return
	}

	if hasFile {
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
	

	w.Header().Set("HX-Redirect", "/admin/blog/preview/" + strconv.Itoa(int(blog.Id)))
}

func HandleUpdateBlogPreview(w http.ResponseWriter, r *http.Request) {

	paramId := chi.URLParam(r, "id")
	id, err := strconv.Atoi(paramId)	
	if err != nil {
		HandleClientError(w, err)
		return
	}

	updatedBlog := types.Blog{}
	updatedBlog.Id = int32(id)
	updatedBlog.Title = sanitize(r.FormValue("title"))
	updatedBlog.Body = sanitize(r.FormValue("body"))
	updatedBlog.Excerpt= sanitize(r.FormValue("excerpt"))
	updatedBlog.IsPublished = false

	hasUniqueTitle, err := hasUniqueTitle(r, updatedBlog.Title, updatedBlog.Id)
	if err != nil {
		HandleServerError(w, r, err)
		return
	}
	if !hasUniqueTitle {
		HandleClientError(w, fmt.Errorf("title '%v' is already taken", updatedBlog.Title))
		return
	}

	hasFile := true 
	thumbnail, header, err := r.FormFile("thumbnail") 
	if err == http.ErrMissingFile {
		hasFile = false 
	} else if err != nil {
		HandleClientError(w, err)
		return
	}

	blog, err := db.UpdateBlog(r, updatedBlog)
	if err != nil {
		HandleServerError(w, r, err)
		return
	}

	if hasFile {
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

	w.Header().Set("HX-Redirect", "/admin/blog/preview/" + strconv.Itoa(int(blog.Id)))
}

func HandlePublishBlog(w http.ResponseWriter, r *http.Request) {
	paramId := chi.URLParam(r, "id")
	id, err := strconv.Atoi(paramId)	
	if err != nil {
		HandleNotFound(w, r)
		return
	}

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
        fmt.Fprintf(w, "<span id=\"publish-info\" class=\"sub-text warning\">*missing fields: %v</span>", strings.Join(missingFields, ", "))
        return
    }

	updatedBlog := blog
	updatedBlog.IsPublished = true

	blog, err = db.UpdateBlog(r, updatedBlog)
	if err != nil {
		HandleNotFound(w, r)
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

	sessionUser, err := GetSessionUser(r)
	if err != nil {
		HandleServerError(w, r, err)
		return
	}

    blogs, err := db.GetBlogsByCreator(r, sessionUser.Username)
    if err != nil && err != sql.ErrNoRows {
        HandleServerError(w, r, err)
        return
    }

	d := types.Admin{Blogs: blogs}

	renderPage(pages["admin"], w, r, d)
}

func HandleGetPlacementTest(w http.ResponseWriter, r *http.Request) {
	sessionUser, err := GetSessionUser(r)
	if err != nil {
		HandleServerError(w, r, err)
		return
	}
	if sessionUser.SoulScore > 0 {
		http.Redirect(w, r, "/account", http.StatusSeeOther)
		return
	}

	if pages["placementTest"] == nil {
		var err error
		pages["placementTest"], err = template.ParseFiles(baseTemplate, 
			"templates/pages/placement-test.html", "templates/components/select-input.html", "templates/components/checkbox-input.html",
			"templates/components/test-result.html")
		if err != nil {
			HandleServerError(w, r, err)
			return
		}
	}

	questions, err := db.GetPlacementTestQuestions(r)
	if err != nil {
		HandleServerError(w, r, err)
	}

	questionInputs := types.QuestionInputs{}	
	for _, q := range questions {
		questionInput := types.QuestionInput{}
		answers := make([]types.Answer, len(q.Answers))
		for i := 0; i < len(answers); i++ {
			answers[i] = types.Answer{Text: q.Answers[i], Value: i}
		}
		questionInput.Id = q.Id
		questionInput.Label = q.Text
		questionInput.Answers = answers 
		if q.InputType == "select" {
			questionInputs.SelectQuestions = append(questionInputs.SelectQuestions, questionInput)
		} else if q.InputType == "checkbox" {
			questionInputs.CheckboxQuestions = append(questionInputs.CheckboxQuestions, questionInput)
		}
	}

	renderPage(pages["placementTest"], w, r, questionInputs)
}

func setCookie(w http.ResponseWriter, name string, value string) {
	cookie := &http.Cookie{
		Name:  name,
		Value: value,
		Path:  "/",
	}
	http.SetCookie(w, cookie)
}

func hasUniqueTitle(r *http.Request, title string, id int32) (bool, error) {
	// blog Id is 0 if no results
	blog, err := db.GetBlogByTitle(r, title)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}
	if blog.Id > 0 && blog.Id != id {
		return false, nil
	}
	return true, nil
}
