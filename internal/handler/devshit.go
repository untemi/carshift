package handler

import (
	"encoding/json"
	"math/rand"
	"net/http"

	"github.com/untemi/carshift/internal/db"
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

	car := db.Car{
		Name:       misc.RandString(10, 18),
		Price:      float64(rand.Int31n(200) + 200),
		UserID:     1,
		DistrictID: rand.Intn(7) + 1,
		StartAt:    startDate,
		EndAt:      endDate,
	}

	db.AddCar(&car)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(car)
}
