package handler

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/untemi/carshift/internal/db"
	"github.com/untemi/carshift/internal/db/sqlc"
	"github.com/untemi/carshift/internal/misc"
	"github.com/untemi/carshift/internal/template"
)

var tabs = []template.Tab{
	{Name: "Account", Content: template.SettingsAccount(), URL: "/settings/0"},
	{Name: "Profile", Content: template.SettingsProfile(), URL: "/settings/1"},
}

const maxProfileSize = 2 << 20

func GETsettings(w http.ResponseWriter, r *http.Request) {
	template.Settings().Render(r.Context(), w)
}

func GETsettingsTabs(w http.ResponseWriter, r *http.Request) {
	// HACK middleware instead
	isHTMX := r.Header.Get("HX-Request")
	if isHTMX != "true" {
		http.Error(w, "Not allowed", http.StatusMethodNotAllowed)
		return
	}

	sel, err := strconv.Atoi(r.PathValue("tab"))
	if err != nil || len(tabs) < sel+1 {
		http.NotFound(w, r)
		return
	}

	template.Tabbed(tabs, sel, "#settings-tabs").Render(r.Context(), w)
}

func POSTsettingsAccount(w http.ResponseWriter, r *http.Request) {
	u, ok := r.Context().Value("userdata").(sqlc.User)
	if !ok {
		template.AlertError("internal error").Render(r.Context(), w)
		return
	}
	ou := u

	r.ParseForm()
	username := strings.TrimSpace(r.FormValue("username"))
	email := strings.TrimSpace(r.FormValue("email"))
	phone := strings.TrimSpace(r.FormValue("phone"))

	// Verifying username
	if !misc.ValidateUsername(username, true) {
		template.AlertError("invalid username").Render(r.Context(), w)
		return
	}

	// Verifying phone
	if phone != "" && !misc.ValidatePhone(phone) {
		template.AlertError("invalid phone").Render(r.Context(), w)
		return
	}

	// Verifying email
	if !misc.ValidateEmail(email) {
		template.AlertError("invalid email").Render(r.Context(), w)
	}

	// Cheching if username is used
	if u.Username != username {
		e, err := db.IsUsernameUsed(r.Context(), username)
		if err != nil {
			log.Printf("DB: Error checking user existence: %v", err)
			template.AlertError("internal error").Render(r.Context(), w)
			return
		}

		if e {
			template.AlertError("username is already taken").Render(r.Context(), w)
			return
		}
	}

	u.Username = username
	u.Email = email
	u.Phone = phone

	if ou == u {
		template.AlertWarning("no changes").Render(r.Context(), w)
		return
	}

	err := db.UpdateUser(r.Context(), &u)
	if err != nil {
		log.Printf("DB: Error updating user: %v", err)
		template.AlertError("internal error").Render(r.Context(), w)
		return
	}

	template.AlertSuccess("saved").Render(r.Context(), w)
}

func POSTsettingsProfile(w http.ResponseWriter, r *http.Request) {
	u, ok := r.Context().Value("userdata").(sqlc.User)
	if !ok {
		template.AlertSuccess("internal error").Render(r.Context(), w)
		return
	}
	ou := u

	r.ParseForm()
	firstname := strings.TrimSpace(r.FormValue("firstname"))
	lastname := strings.TrimSpace(r.FormValue("lastname"))

	if !misc.ValidateName(firstname, false) {
		template.AlertError("invalid firstname").Render(r.Context(), w)
		return
	}

	if !misc.ValidateName(lastname, true) {
		template.AlertError("invalid lastname").Render(r.Context(), w)
		return
	}

	u.Firstname = misc.FormaterName(firstname)
	u.Lastname = misc.FormaterName(lastname)

	if ou == u {
		template.AlertWarning("no changes").Render(r.Context(), w)
		return
	}

	err := db.UpdateUser(r.Context(), &u)
	if err != nil {
		log.Printf("DB: Error updating user: %v", err)
		template.AlertError("internal error").Render(r.Context(), w)
		return
	}

	template.AlertSuccess("saved").Render(r.Context(), w)
}

func POSTsettingsUpdatePFP(w http.ResponseWriter, r *http.Request) {
	template.AlertSuccess("alright we started").Render(r.Context(), w)
}
