package handler

import (
	"log"
	"net/http"

	"github.com/untemi/carshift/internal/db"
	"github.com/untemi/carshift/internal/db/sqlc"
	"github.com/untemi/carshift/internal/template"
)

func GETprofileSelf(w http.ResponseWriter, r *http.Request) {
	u, ok := r.Context().Value("userdata").(sqlc.User)
	if !ok {
		log.Println("SERVER: error fetching prop userdata")
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	template.Profile(u, true).Render(r.Context(), w)
}

func GETprofile(w http.ResponseWriter, r *http.Request) {
	l, ok := r.Context().Value("logged").(bool)
	if !ok {
		log.Println("SERVER: error fetching prop logged")
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	u := sqlc.User{
		Username: r.PathValue("username"),
	}

	e, err := db.IsUsernameUsed(r.Context(), u.Username)
	if err != nil {
		log.Printf("DB: Error checking user existence: %v", err)
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	// Cheching if username exists
	if !e {
		http.NotFound(w, r)
		return
	}

	// Fetching user
	err = db.FillUser(r.Context(), &u)
	if err != nil {
		log.Printf("DB: Error fetching user: %v", err)
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	// Checking if you are checking ;)
	if l {
		meu, ok := r.Context().Value("userdata").(sqlc.User)
		if !ok {
			log.Println("SERVER: error fetching prop logged")
			http.Error(w, "Internal Error", http.StatusInternalServerError)
			return
		}

		if meu.ID == u.ID {
			http.Redirect(w, r, "/me", http.StatusSeeOther)
			return
		}
	}

	template.Profile(u, false).Render(r.Context(), w)
}
