package internal

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	gm "github.com/go-chi/chi/v5/middleware"

	h "github.com/untemi/carshift/internal/handler"
	m "github.com/untemi/carshift/internal/middleware"
	"github.com/untemi/carshift/internal/view"
)

type ServerConfig struct {
	Address string
}

func Start(c ServerConfig) {
	r := chi.NewRouter()

	// Global middlewares
	r.Use(gm.Logger)
	r.Use(gm.Recoverer)
	r.Use(h.SM.LoadAndSave)

	// Static and general stuff
	r.Group(func(r chi.Router) {
		r.Get("/logout", h.EndSession)
		r.Get("/favicon.ico", view.ServeFavicon)
	})

	// HTMX
	r.Group(func(r chi.Router) {
		r.Get("/htmx/alert", h.HtmxAlert)
	})

	// Public
	r.Group(func(r chi.Router) {
		r.Use(m.FetchLogin)

		r.Get("/", h.GEThome)
		r.Get("/carfinder", h.GETcarFinder)
		r.Get("/userfinder", h.GETuserFinder)
		r.Get("/profile/{username}", h.GETprofile)

		r.Post("/carfinder", h.POSTcarFinder)
		r.Post("/userfinder", h.POSTuserFinder)
	})

	// Users-only
	r.Group(func(r chi.Router) {
		r.Use(m.FetchLogin)
		r.Use(m.UserOnly)

		r.Get("/me", h.GETprofileSelf)
		r.Get("/settings", h.GETsettings)
		r.Get("/settings/{tab}", h.GETsettingsTabs)

		r.Post("/settings/profile", h.POSTsettingsProfile)
		r.Post("/settings/account", h.POSTsettingsAccount)
		r.Post("/settings/pfp", h.POSTsettingsUpdatePFP)

		// Dev Shit
		r.Get("/dev/randcar", h.DevAddRandCar)
	})

	// Guest-only
	r.Group(func(r chi.Router) {
		r.Use(m.FetchLogin)
		r.Use(m.GuestOnly)

		r.Get("/login", h.GETlogin)
		r.Get("/register", h.GETregister)

		r.Post("/login", h.POSTlogin)
		r.Post("/register", h.POSTregister)
	})

	// Files serving
	view.FileServer(r, "/static", "static")
	view.FileServer(r, "/pictures", "pictures")

	server := http.Server{
		Addr:    c.Address,
		Handler: r,
	}

	log.Println("SERVER: running on", server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("SERVER: Error fetching user %v", err)
	}
}
