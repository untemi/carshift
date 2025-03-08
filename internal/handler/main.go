package handler

import (
	"net/http"

	"github.com/untemi/carshift/internal/template"
)

func reTargetAlert(message string, w http.ResponseWriter, r *http.Request) {
	w.Header().Add("HX-Retarget", "#hxtoast")
	w.Header().Add("HX-Reswap", "beforeend")
	template.AlertError(message).Render(r.Context(), w)
}
