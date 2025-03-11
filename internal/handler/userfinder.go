package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/untemi/carshift/internal/db"
	"github.com/untemi/carshift/internal/misc"
	"github.com/untemi/carshift/internal/template"
)

func GETuserFinder(w http.ResponseWriter, r *http.Request) {
	template.UserFinder().Render(r.Context(), w)
}

func POSTuserFinder(w http.ResponseWriter, r *http.Request) {
	var page int64
	err := r.ParseForm()
	if err != nil {
		reTargetAlert("bad data", w, r)
		return
	}

	query := r.FormValue("username")
	pageStr := r.FormValue("page")

	if query == "" {
		reTargetAlert("missing data", w, r)
		return
	}

	if !misc.ValidateUsername(query, false) {
		reTargetAlert("invalid username", w, r)
		return
	}

	if pageStr != "" {
		page, _ = strconv.ParseInt(pageStr, 10, 64)
	}

	users, err := db.FetchUsers(r.Context(), query, 10, page)
	if err != nil {
		log.Printf("DB: Error fetching users: %v", err)
		reTargetAlert("internal error", w, r)
		return
	}

	hxVars := fmt.Sprintf(`{"username": "%s","page": %d}`, query, page+1)
	template.UserFinderResults(users, hxVars).Render(r.Context(), w)
}
