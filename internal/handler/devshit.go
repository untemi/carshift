package handler

import (
	"encoding/json"
	"math/rand"
	"net/http"

	"github.com/untemi/carshift/internal/db"
	"github.com/untemi/carshift/internal/db/sqlc"
	"github.com/untemi/carshift/internal/misc"
)

func DevAddRandCar(w http.ResponseWriter, r *http.Request) {
	startDate := misc.RanDate()
	endDate := misc.RanDate()
	for {
		if endDate.Unix() > startDate.Unix() {
			break
		}

		endDate = misc.RanDate()
	}

	car := sqlc.Car{
		Name:       misc.RandString(10, 18),
		Price:      float64(rand.Int31n(200) + 200),
		OwnerID:    1,
		DistrictID: rand.Int63n(7) + 1,
		StartAt:    misc.TimeToNull(startDate),
		EndAt:      misc.TimeToNull(endDate),
	}

	db.AddCar(r.Context(), &car)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(car)
}
