package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/untemi/carshift/internal/db"
	"github.com/untemi/carshift/internal/template"
)

var dateFormat = "02-01-2006"

func GETcarFinder(w http.ResponseWriter, r *http.Request) {
	template.CarFinder().Render(r.Context(), w)
}

func POSTcarFinder(w http.ResponseWriter, r *http.Request) {
	var page int64
	var startDate time.Time
	var endDate time.Time

	err := r.ParseForm()
	if err != nil {
		reTargetAlert("bad data", w, r)
		return
	}

	district := r.FormValue("district")
	startDateStr := r.FormValue("startdate")
	endDateStr := r.FormValue("enddate")
	query := r.FormValue("carname")
	pageStr := r.FormValue("page")

	if district == "" {
		reTargetAlert("Select district", w, r)
		return
	}

	// Parse Dates
	if startDateStr != "" && endDateStr != "" {
		startDate, err = time.Parse(dateFormat, startDateStr)
		endDate, err = time.Parse(dateFormat, endDateStr)
		if err != nil {
			reTargetAlert("bad date format", w, r)
			return
		}

		// Chack dates
		if endDate.Unix() < startDate.Unix() {
			reTargetAlert("start date is greated than end date", w, r)
			return
		}
	}

	if pageStr != "" {
		page, _ = strconv.ParseInt(pageStr, 10, 64)
	}

	cars, err := db.FetchCars(r.Context(), district, startDate, endDate, query, 5, page)
	if err != nil {
		log.Printf("DB: Error fetching cars: %v", err)
		reTargetAlert("internal error", w, r)
		return
	}

	hxVars := fmt.Sprintf(
		`{"district": "%s","startdate": "%s","enddate": "%s","carname": "%s","page": %d}`,
		district,
		startDateStr,
		endDateStr,
		query,
		page+1,
	)
	template.CarFinderResults(cars, hxVars).Render(r.Context(), w)
}
