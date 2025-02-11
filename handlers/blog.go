package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/daltbunker/soul_climbers/db"
	"github.com/go-chi/chi"
)

func GetBlogImg(w http.ResponseWriter, r *http.Request) {
	paramId := chi.URLParam(r, "id")
	paramImgName := chi.URLParam(r, "imgName")
	id, err := strconv.Atoi(paramId)	
	if err != nil {
		HandleClientError(w, fmt.Errorf("invalid param 'id': %v", paramId))
		return
	}

	dbBlogImg, err := db.GetBlogImg(r, int32(id))
	if err != nil {
		HandleServerError(w, r, err)
		return
	}
	if dbBlogImg.ImgName != paramImgName {
		HandleClientError(w, fmt.Errorf("image '%v' not found", paramImgName))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(dbBlogImg.Img)
}

func DeleteBlogImg(w http.ResponseWriter, r *http.Request)  {
	paramId := chi.URLParam(r, "id")
	id, err := strconv.Atoi(paramId)	
	if err != nil {
		HandleClientError(w, fmt.Errorf("invalid param 'id': %v", paramId))
		return
	}

	_, err = db.DeleteBlogImg(r, int32(id))
	if err != nil {
		HandleServerError(w, r, err)
		return
	}

	fmt.Fprint(w, "<div class=\"input-container\"><label for=\"thumbnail\">Thumbnail:</label><input type=\"file\" name=\"thumbnail\" id=\"thumbnail\"></div>")
}

func DeleteBlog(w http.ResponseWriter, r *http.Request)  {
	paramId := chi.URLParam(r, "id")
	id, err := strconv.Atoi(paramId)	
	if err != nil {
		HandleClientError(w, fmt.Errorf("invalid param 'id': %v", paramId))
		return
	}

	err = db.DeleteBlog(r, int32(id))
	if err != nil {
		HandleServerError(w, r, err)
		return
	}

	w.Header().Set("HX-Redirect", "/admin")
}