package db

import (
	"context"
	"time"

	"github.com/untemi/carshift/internal/db/sqlc"
	"github.com/untemi/carshift/internal/misc"
)

func AddCar(ctx context.Context, c *sqlc.Car) (int64, error) {
	return runner.CreateCar(
		ctx,
		sqlc.CreateCarParams{
			Name:       c.Name,
			Price:      c.Price,
			StartAt:    c.StartAt,
			EndAt:      c.EndAt,
			OwnerID:    c.OwnerID,
			DistrictID: c.DistrictID,
		},
	)
}

func FetchCars(
	ctx context.Context,
	district string,
	startDate time.Time,
	endDate time.Time,
	query string,
	limit int64,
	page int64,
) (*[]sqlc.Car, error) {
	cars, err := runner.QueryCars(
		ctx,
		sqlc.QueryCarsParams{
			DistrictName: district,
			Name:         "%" + query + "%",
			StartAt:      misc.TimeToNull(startDate),
			EndAt:        misc.TimeToNull(endDate),
			Limit:        limit,
			Offset:       page,
		})
	return &cars, err
}
