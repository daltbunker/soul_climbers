package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/daltbunker/soul_climbers/db"
	"github.com/daltbunker/soul_climbers/types"
)

func HandleSubmitPlacementTest(w http.ResponseWriter, r *http.Request) {
	sessionUser, err := GetSessionUser(r)
	if err != nil {
		HandleServerError(w, r, err)
		return
	}
	username := sessionUser.Username

	if sessionUser.SoulScore > 0 {
		HandleClientError(w, fmt.Errorf("user %v has already taken the test", username))
		return
	}

	testQuestions, err := db.GetPlacementTestQuestions(r)
	if err != nil {
		HandleServerError(w, r, err)
		return
	}

	points := 0
	possiblePoints := 0
	for _, val := range testQuestions {
		possiblePoints += val.PossiblePoints
		if val.InputType == "select" {
			selectPoints, err := getSelectPoints(val, r)
			if err != nil && err.Error()[0:7] == "missing" {
				HandleClientError(w, err)
				return
			}
			if err != nil {
				HandleServerError(w, r, err)
				return
			}
			points += selectPoints
		} else if val.InputType == "checkbox" {
			checkboxPoints, err := getCheckboxPoints(val, r)
			if err != nil {
				HandleServerError(w, r, err)
				return
			}
			points += checkboxPoints 
		}
	}

	err = db.InsertPlacementTest(r, username, int32(points))
	if err != nil {
		HandleServerError(w, r, err)
	}

	pointsPercentage := int(float32(points) / float32(possiblePoints) * 100)
	soulScore, err := db.SetUserSoulScore(r, username, int32(pointsPercentage))
	if err != nil {
		HandleServerError(w, r, err)
	}

	err = UpdateSession(r, w, "soul_score", soulScore)
	if err != nil {
		HandleServerError(w, r, err)
	}

	renderComponent(w, "placementTest", "test-result", pointsPercentage)
}


func getSelectPoints(question types.Question, r *http.Request) (int, error) {
	formValue := r.FormValue("q-" + question.Id)
	if formValue == "" {
		return 0, fmt.Errorf("missing form value: q-%v", question.Id)
	}

	pointsIndex, err := strconv.Atoi(formValue)
	if err != nil {
		return 0, err
	}

	points, err := strconv.Atoi(question.AnswerPoints[pointsIndex])
	if err != nil {
		return 0, err
	}

	return points, nil
}

func getCheckboxPoints(question types.Question, r *http.Request) (int, error) {
	checkboxPoints := 0
	for i := range question.Answers {
		formValue := r.FormValue("q-" + question.Id + "-" + strconv.Itoa(i))

		pointsIndex, err := strconv.Atoi(formValue)
		if err != nil {
			continue
		}

		points, err := strconv.Atoi(question.AnswerPoints[pointsIndex])
		if err != nil {
			return 0, err
		}

		checkboxPoints += points
	}

	return checkboxPoints, nil
}