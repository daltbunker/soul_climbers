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

	dbBlogImg, err := db.GetBlogImg(r, int32(id))
	if err != nil {
		HandleServerError(w, r, err)
		return
	}
	if dbBlogImg.ImgName != paramImgName {
		HandleClientError(w, fmt.Errorf("Image '%v' not found", paramImgName))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(dbBlogImg.Img)
}

func DeleteBlogImg(w http.ResponseWriter, r *http.Request)  {
	paramId := chi.URLParam(r, "id")
	id, err := strconv.Atoi(paramId)	

	_, err = db.DeleteBlogImg(r, int32(id))
	if err != nil {
		HandleServerError(w, r, err)
		return
	}

	fmt.Fprint(w, "<label for=\"thumbnail\">Thumbnail</label><input type=\"file\" name=\"thumbnail\" id=\"thumbnail\">")
}