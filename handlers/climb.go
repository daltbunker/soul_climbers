package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/daltbunker/soul_climbers/db"
	"github.com/daltbunker/soul_climbers/types"
	"github.com/go-chi/chi"
)

func ClimbSearch(w http.ResponseWriter, r *http.Request) {
	// POSTGRES search value min is 0.3 / 1. Can change value for more strict matching
	searchType := sanitize(r.FormValue("search-type"))
	searchValue := sanitize(r.FormValue("search-value"))
	climbType := sanitize(r.FormValue("climb-type"))

	if searchValue == "" {
		HandleClientError(w, fmt.Errorf("search-value cannot be empty: %v", searchType))
		return
	}

	if searchType == "area" {
		areaResults, err := db.SearchAreas(r, searchValue)
		if err != nil {
			HandleServerError(w, r, err)
			return
		}

		subAreasResults, err := db.SearchSubAreas(r, searchValue)
		if err != nil {
			HandleServerError(w, r, err)
			return
		}

		areaResults = append(areaResults, subAreasResults...)

		renderComponent(w, "climb-search", "area-search-results", areaResults)
	} else if searchType == "climb" {
		climbResults, err := db.SearchClimbs(r, climbType, searchValue)
		if err != nil {
			HandleServerError(w, r, err)
			return
		}

		renderComponent(w, "climb-search", "climb-search-results", climbResults)
	} else {
		HandleClientError(w, fmt.Errorf("invalid search-type: %v", searchType))
		return
	}
}

func SearchClimbArea(w http.ResponseWriter, r *http.Request) {
	areaQuery := r.URL.Query().Get("area")
	areaResults, err := db.SearchAreas(r, areaQuery)
	if err != nil {
		HandleServerError(w, r, err)
		return
	}

	results := []types.Result{}
	for i := 0; i < len(areaResults); i++ {
		result := types.Result{Name: areaResults[i].Name, Id: int(areaResults[i].AreaId)}
		results = append(results, result)
	}

	searchResults := types.SearchResult{Results: results, InputId: "area"}
	renderComponent(w, "climbForm", "search-results", searchResults)
}

func SearchSubArea(w http.ResponseWriter, r *http.Request) {
	sessionUser, err := GetSessionUser(r)
	if err != nil {
		HandleServerError(w, r, err)
		return
	}

	climbDraft, err := db.GetClimbDraft(r, sessionUser.Username)
	if err != nil && err != sql.ErrNoRows {
		HandleServerError(w, r, err)
		return
	}

	subAreaQuery := r.URL.Query().Get("sub-area")
	subAreaResults, err := db.SearchSubAreas(r, subAreaQuery)

	if err != nil {
		HandleServerError(w, r, err)
		return
	}

	results := []types.Result{}
	for i := 0; i < len(subAreaResults); i++ {
		if subAreaResults[i].Name == climbDraft.Area {
			result := types.Result{Name: subAreaResults[i].SubArea, Id: int(subAreaResults[i].AreaId)}
			results = append(results, result)
		}
	}

	searchResults := types.SearchResult{Results: results, InputId: "sub-area"}
	renderComponent(w, "climbForm", "search-results", searchResults)
}

func HandleDeleteSubArea(w http.ResponseWriter, r *http.Request) {
	subArea := r.URL.Query().Get("sub-area")

	sessionUser, err := GetSessionUser(r)
	if err != nil {
		HandleServerError(w, r, err)
		return
	}

	climbDraft, err := db.GetClimbDraft(r, sessionUser.Username)
	if err != nil && err != sql.ErrNoRows {
		HandleServerError(w, r, err)
		return
	}

	if strings.Index(climbDraft.SubAreas, subArea) == -1 {
		HandleClientError(w, fmt.Errorf("sub-area not found: %v", subArea))
		return
	}

	oldSubAreas := strings.Split(climbDraft.SubAreas, ",")
	newSubAreas := []string{}
	for _, s := range oldSubAreas {
		if s != subArea {
			newSubAreas = append(newSubAreas, s)
		}
	}
	climbDraft.SubAreas = strings.Join(newSubAreas, ",")

	err = db.CreateClimbDraft(r, climbDraft, sessionUser.Username)
	if err != nil {
		HandleServerError(w, r, err)
		return
	}

	renderComponent(w, "climbForm", "sub-area-element", newSubAreas)
}

func HandleAddSubArea(w http.ResponseWriter, r *http.Request) {
	subArea := r.FormValue("sub-area")

	if subArea == "" {
		HandleClientError(w, fmt.Errorf("sub-area can not be empty"))
		return
	}

	sessionUser, err := GetSessionUser(r)
	if err != nil {
		HandleServerError(w, r, err)
		return
	}

	climbDraft, err := db.GetClimbDraft(r, sessionUser.Username)
	if err != nil && err != sql.ErrNoRows {
		HandleServerError(w, r, err)
		return
	}

	if strings.Contains(climbDraft.SubAreas, subArea) {
		HandleClientError(w, fmt.Errorf("sub-area '%v' is already added", subArea))
		return
	}

	if climbDraft.SubAreas != "" {
		subArea = "," + subArea
	}

	climbDraft.SubAreas += subArea
	err = db.CreateClimbDraft(r, climbDraft, sessionUser.Username)
	if err != nil {
		HandleServerError(w, r, err)
		return
	}
	renderComponent(w, "climbForm", "sub-area-element", strings.Split(climbDraft.SubAreas, ","))
}

func HandleClimbForm(w http.ResponseWriter, r *http.Request) {
	sessionUser, err := GetSessionUser(r)
	if err != nil {
		HandleServerError(w, r, err)
		return
	}

	formPartParam := chi.URLParam(r, "part")
	formPart, err := strconv.Atoi(formPartParam)
	if err != nil {
		log.Println(err)
	}

	countries := []string{"United States", "France"}
	routeTypes := []string{"boulder", "sport", "trad"}

	climbDraft, err := db.GetClimbDraft(r, sessionUser.Username)
	if err != nil && err != sql.ErrNoRows {
		HandleServerError(w, r, err)
		return
	}

	switch formPart {
	case 1:
		fdClimbName := sanitize(r.FormValue("climb-name"))
		fdRouteType := sanitize(r.FormValue("route-type"))
		fdCountry := sanitize(r.FormValue("country"))

		formData := types.ClimbForm{}
		hasError := false
		if len(fdClimbName) < 1 || len(fdClimbName) > 30 {
			formData.NameError = "*required"
			hasError = true
		}
		if fdRouteType == "" {
			formData.RouteTypeError = "*required"
			hasError = true
		}
		if fdCountry == "" {
			formData.CountryError = "*required"
			hasError = true
		}
		if hasError {
			formData.Part = 1
			formData.Name = fdClimbName
			formData.CountryOptions = newFormOptions(countries, fdCountry)
			formData.RouteTypeOptions = newFormOptions(routeTypes, fdRouteType)
			renderComponent(w, "climbForm", "climb-form-1", formData)
			return
		} else {
			climbDraft.Name = fdClimbName
			climbDraft.Country = fdCountry
			climbDraft.RouteType = fdRouteType

			err = db.CreateClimbDraft(r, climbDraft, sessionUser.Username)
			if err != nil {
				HandleServerError(w, r, err)
				return
			}
		}
	case 2:
		fdArea := sanitize(r.FormValue("area"))
		fdAreaId := sanitize(r.FormValue("area-id"))
		areaId, err := strconv.Atoi(fdAreaId)
		if err == nil {
			climbDraft.AreaId = int32(areaId)
		}
		formData := types.ClimbForm{}
		hasError := false
		if fdArea == "" {
			formData.AreaError = "*required"
			hasError = true
		}
		if hasError {
			formData.Part = 2
			formData.Name = climbDraft.Name
			formData.RouteType = climbDraft.RouteType
			formData.Country = climbDraft.Country
			renderComponent(w, "climbForm", "climb-form-2", formData)
			return
		} else {
			climbDraft.Area = fdArea
			err = db.CreateClimbDraft(r, climbDraft, sessionUser.Username)
			if err != nil {
				HandleServerError(w, r, err)
				return
			}
		}
	case 3:
		// Nothing to do, add/delete sub-area handlers are used for part 3
	case 4:
		err = db.CreateClimbWithArea(r, climbDraft, sessionUser.Username)
		if err != nil {
			HandleServerError(w, r, err)
			return
		}
		err = db.DeleteClimbDraft(r, sessionUser.Username)
		if err != nil {
			HandleServerError(w, r, err)
			return
		}
		w.Header().Set("HX-Redirect", "/climbform/1")
		return
	default:
		HandleClientError(w, fmt.Errorf("invalid form part '%v'", formPart))
		return
	}

	w.Header().Set("HX-Redirect", "/climbform/"+strconv.Itoa(formPart+1))
}

func newFormOptions(options []string, preSelected string) []types.FormOption {
	formOptions := []types.FormOption{}
	for _, option := range options {
		selected := false
		if option == preSelected {
			selected = true
		}
		formOption := types.FormOption{Name: option, Selected: selected}
		formOptions = append(formOptions, formOption)
	}
	return formOptions
}
