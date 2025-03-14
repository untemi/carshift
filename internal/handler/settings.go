package handler

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/google/uuid"

	"github.com/untemi/carshift/internal/db"
	"github.com/untemi/carshift/internal/db/sqlc"
	"github.com/untemi/carshift/internal/misc"
	"github.com/untemi/carshift/internal/template"
)

const maxPfpSize = 5 << 20

var (
	tabs = []template.Tab{
		{Name: "Profile", Content: template.SettingsProfile(), URL: "/settings/0"},
		{Name: "Account", Content: template.SettingsAccount(), URL: "/settings/1"},
	}
	validPfpTypes = map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
	}
)

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

	// DPS parses
	err := r.ParseForm()
	if err != nil {
		log.Printf("SERVER: Error parsing form : %v", err)
		template.AlertError("bad data").Render(r.Context(), w)
		return
	}

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
		return
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

	err = db.UpdateUser(r.Context(), &u)
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

	// DPS parses
	err := r.ParseForm()
	if err != nil {
		log.Printf("SERVER: Error parsing form : %v", err)
		template.AlertError("bad data").Render(r.Context(), w)
		return
	}

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

	err = db.UpdateUser(r.Context(), &u)
	if err != nil {
		log.Printf("DB: Error updating user: %v", err)
		template.AlertError("internal error").Render(r.Context(), w)
		return
	}

	template.AlertSuccess("saved").Render(r.Context(), w)
}

func POSTsettingsUpdatePFP(w http.ResponseWriter, r *http.Request) {
	// limit body size
	r.Body = http.MaxBytesReader(w, r.Body, maxPfpSize)
	err := r.ParseMultipartForm(maxPfpSize)
	if err != nil {
		template.AlertError("5MB is the max").Render(r.Context(), w)
		return
	}

	// fetching user from ctx
	u, ok := r.Context().Value("userdata").(sqlc.User)
	if !ok {
		template.AlertError("internal error").Render(r.Context(), w)
		return
	}

	// taking the file reader
	fileReader, fileHeader, err := r.FormFile("profile")
	if err != nil {
		template.AlertError("bad data").Render(r.Context(), w)
		return
	}
	defer fileReader.Close()

	// reading the first 512 bytes for verifying MIME
	buf := make([]byte, 512)
	_, err = fileReader.Read(buf)
	if err != nil {
		template.AlertError("bad data").Render(r.Context(), w)
		return
	}
	fileType := http.DetectContentType(buf)

	// reader reset
	fileReader.Seek(0, io.SeekStart)

	// verifying MIME type
	if !validPfpTypes[fileType] {
		template.AlertError("not valid MIME type").Render(r.Context(), w)
		return
	}

	// reading the whole file
	fileBytes, err := io.ReadAll(fileReader)
	if err != nil {
		template.AlertError("bad data").Render(r.Context(), w)
		return
	}

	// making the file name
	fileSplit := strings.Split(fileHeader.Filename, ".")
	fileExtention := fileSplit[len(fileSplit)-1]
	filename := fmt.Sprintf("pfp_%s.%s", uuid.New().String(), fileExtention)

	// writing the file to disk
	err = os.WriteFile(fmt.Sprintf("pictures/pfp/%s", filename), fileBytes, 0644)
	if err != nil {
		log.Printf("SERVER: Error saving file : %v", err)
		template.AlertError("internal error").Render(r.Context(), w)
		return
	}

	// removing the old file
	if u.PfpName != "" && !strings.Contains(u.PfpName, "/") {
		err = os.Remove(fmt.Sprintf("pictures/pfp/%s", u.PfpName))
		if err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				log.Printf("SERVER: Error removing file : %v", err)
				template.AlertError("internal error").Render(r.Context(), w)
				return
			}
		}
	}

	// saving the name to db
	u.PfpName = filename
	err = db.UpdateUser(r.Context(), &u)
	if err != nil {
		log.Printf("DB: Error updating user: %v", err)
		template.AlertError("internal error").Render(r.Context(), w)
		return
	}

	template.AlertSuccess("saved").Render(r.Context(), w)
}
