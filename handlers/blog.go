package handlers

import (
	"net/http"

	"github.com/daltbunker/soul_climbers/db"
	"github.com/daltbunker/soul_climbers/types"
)

func GetBlogs(w http.ResponseWriter, r *http.Request) []types.Blog {
	blogs, err := db.GetAllBlogs(r)
	if err != nil {
		HandleServerError(w, err)
	}
	return blogs
}

func HandleNewBlog(w http.ResponseWriter, r *http.Request) {

	//TODO: this method should publish a draft	
	// fmt.Fprintf(w, "<h2>%s</h2> author %s <p>%s</p>", savedBlog.Title, savedBlog.CreatedBy, savedBlog.Body)
}

