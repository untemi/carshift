package handler

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
	_ "github.com/mattn/go-sqlite3"

	"github.com/untemi/carshift/internal/misc"
	"github.com/untemi/carshift/internal/template"
)

var SM = scs.New()

func Init() (misc.DBClose, error) {
	conn, err := sql.Open("sqlite3", "session.db")
	if err != nil {
		return nil, err
	}

	SM.Lifetime = time.Hour * 365 * 24
	SM.Store = sqlite3store.New(conn)

	log.Println("SCS: up and running")
	return conn.Close, nil
}

func Setup(ctx context.Context) error {
	conn, err := sql.Open("sqlite3", "session.db")
	if err != nil {
		return err
	}

	if _, err := conn.ExecContext(
		ctx,
		`CREATE TABLE sessions (
	    token TEXT PRIMARY KEY,
	    data BLOB NOT NULL,
	    expiry REAL NOT NULL);`); err != nil {
		return err
	}

	if _, err := conn.ExecContext(
		ctx,
		`CREATE INDEX sessions_expiry_idx
      ON sessions(expiry)`); err != nil {
		return err
	}

	conn.Close()
	return nil
}

func IsLogged(ctx context.Context) bool {
	return SM.Exists(ctx, "userId")
}

func EndSession(w http.ResponseWriter, r *http.Request) {
	err := SM.RenewToken(r.Context())
	if err != nil {
		log.Printf("SCS: Error session renew %v", err)
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	SM.Destroy(r.Context())
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func reTargetAlert(message string, w http.ResponseWriter, r *http.Request) {
	w.Header().Add("HX-Retarget", "#hxtoast")
	w.Header().Add("HX-Reswap", "beforeend")
	template.AlertError(message).Render(r.Context(), w)
}
