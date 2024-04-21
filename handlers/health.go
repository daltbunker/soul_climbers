package handlers

import "net/http"

func HandleGetHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}
