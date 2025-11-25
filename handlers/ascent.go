package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/daltbunker/soul_climbers/db"
	"github.com/daltbunker/soul_climbers/types"
)

func AddAscent(w http.ResponseWriter, r *http.Request) {

	date := sanitize(r.FormValue("date"))
	grade := sanitize(r.FormValue("grade"))
	rating := sanitize(r.FormValue("rating"))
	attempts := sanitize(r.FormValue("attempts"))
	weight := sanitize(r.FormValue("weight"))
	comment := sanitize(r.FormValue("comment"))
	//TODO: add to form
	// ascentType := sanitize(r.FormValue("ascentType"))

	dateError := ""
	gradeError := ""
	ratingError := ""
	attemptError := ""
	weightError := ""
	commentError := ""
	formError := false

	//TODO: validate inputs

	if date == "" {
		formError = true 
		dateError = "required"
	}
	if grade == "" {
		formError = true 
		gradeError = "required"
	}
	if rating == "" {
		formError = true 
		ratingError = "required"
	}
	if attempts == "" {
		formError = true 
		attemptError = "required"
	}
	if weight == "" {
		formError = true 
		weightError = "required"
	}
	if len(comment) > 500 {
		formError = true 
		commentError = "exceeded max 500 characters"
	}

	if formError {

		//TODO: make constant / reusable
		ratings := []string{"*", "**", "***", "****"}
		weights := []string{"over 200 pounds", "under 200 pounds"}
		attemptOptions := []string{"flash", "soft second go", "onsight", "more than 2 tries"}
		
		ascentForm := types.AscentForm{
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

	climbIdStr := r.URL.Query().Get("climbId")
	climbId, err := strconv.Atoi(climbIdStr)
	if err != nil {
		log.Printf("Climb id must be type int: %v", err)
		HandleNotFound(w, r)
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

	w.Header().Set("HX-Redirect", "/climb/"+climbIdStr)
}