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
	page := 0
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
		page, _ = strconv.Atoi(pageStr)
	}

	query = "%" + query + "%"
	users, err := db.FetchUsers(query, 10, page)
	if err != nil {
		log.Printf("DB: Error fetching users: %v", err)
		reTargetAlert("internal error", w, r)
		return
	}

	hxVars := fmt.Sprintf("{\"username\": \"%s\",\"page\": %d}", r.FormValue("username"), page+1)
	template.UserFinderResults(users, hxVars).Render(r.Context(), w)
}
