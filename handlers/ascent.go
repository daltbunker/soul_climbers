package handlers

import (
	"fmt"
	"log"
	"net/http"
	"slices"
	"strconv"
	"time"

	"github.com/daltbunker/soul_climbers/db"
	"github.com/daltbunker/soul_climbers/types"
	"github.com/daltbunker/soul_climbers/utils"
	"github.com/go-chi/chi"
)

func AddAscent(w http.ResponseWriter, r *http.Request) {
	climbType := r.URL.Query().Get("climbType")
	climbIdStr := r.URL.Query().Get("climbId")
	newAscent := r.URL.Query().Get("newAscent")

	climbId, err := strconv.Atoi(climbIdStr)
	if err != nil {
		log.Printf("Climb id must be type int: %v", err)
		HandleNotFound(w, r)
		return
	}

	date := sanitize(r.FormValue("date"))
	grade := sanitize(r.FormValue("grade"))
	rating := sanitize(r.FormValue("rating"))
	attempts := sanitize(r.FormValue("attempts"))
	weight := sanitize(r.FormValue("weight"))
	comment := sanitize(r.FormValue("comment"))

	dateError := ""
	gradeError := ""
	ratingError := ""
	attemptError := ""
	weightError := ""
	commentError := ""

	ratings := utils.GetAscentRatings()
	weights := utils.GetAscentWeights() 
	attemptOptions := utils.GetAscentAttempts()
	if climbType != "boulder" {
		attemptOptions = append(attemptOptions, "onsight")
	}

	clientErrors := []string{}
	formError := false

	_, err = time.Parse("2006-01-02", date)
	if date == "" {
		formError = true 
		dateError = "required"
	} else if err != nil {
		clientErrors = append(clientErrors, "invalid date format, expected yyyy-mm-dd") 
	}
	if grade == "" {
		formError = true 
		gradeError = "required"
	} else if len(grade) != 7 || string(grade[0]) != "#" {
		clientErrors = append(clientErrors, "invalid grade") 
	}
	if rating == "" {
		formError = true 
		ratingError = "required"
	} else if !slices.Contains(ratings, rating) {
		clientErrors = append(clientErrors, "invalid rating") 
	}
	if attempts == "" {
		formError = true 
		attemptError = "required"
	} else if !slices.Contains(attemptOptions, attempts) {
		clientErrors = append(clientErrors, "invalid attempts") 
	}
	if weight == "" {
		formError = true 
		weightError = "required"
	} else if !slices.Contains(weights, weight) {
		clientErrors = append(clientErrors, "invalid weight") 
	}
	if len(comment) > 500 {
		clientErrors = append(clientErrors, "comment exceeded max 500 characters") 
	}
	if len(clientErrors) > 0 {
		log.Printf("%v", clientErrors)
		HandleClientError(w, fmt.Errorf("invalid input"))
		return
	}
	if formError {
		ascentForm := types.AscentForm{
			ClimbId: climbId,
			ClimbType: climbType,
			NewAscent: newAscent == "true",
			Date: date,
			DateError: dateError,
			Grade: grade,
			GradeError: gradeError,
			RatingOptions: newFormOptions(ratings, rating),
			RatingError: ratingError,
			WeightOptions: newFormOptions(weights, weight),
			WeightError: weightError,
			AttemptOptions: newFormOptions(attemptOptions, attempts),
			AttemptError: attemptError,
			Comment: comment,
			CommentError: commentError,
		}

		renderComponent(w, "ascent-form", "ascent-form", ascentForm)
		return
	}

	over200Pounds := false
	if (weight == "over 200 pounds") {
		over200Pounds = true
	}

	user, err := GetSessionUser(r)
	if err != nil {
		HandleServerError(w, r, err)
		return
	}

	ascent := types.Ascent{
		ClimbId: climbId,
		Grade: grade,
		Rating: rating,
		AscentDate: date,
		Over200Pounds: over200Pounds,
		Attempts: attempts,
		Comment: comment,
		CreatedBy: user.Username,

	}

	err = db.CreateAscent(r, ascent)
	if err != nil {
		HandleServerError(w, r, err)
		return
	}

	w.Header().Set("HX-Refresh", "true")
}

func DeleteAscent(w http.ResponseWriter, r *http.Request) {
	user, err := GetSessionUser(r)
	if err != nil {
		HandleServerError(w, r, err)
		return
	}

	paramClimbId := chi.URLParam(r, "climbId")
	climbId, err := strconv.Atoi(paramClimbId)
	if err != nil {
		log.Printf("Climb id must be type int: %v", err)
		HandleNotFound(w, r)
		return
	}

	err = db.DeleteAscent(r, user.Username, climbId)
	if err != nil {
		HandleServerError(w, r, err)
		return
	}

	w.Header().Set("HX-Refresh", "true")
}