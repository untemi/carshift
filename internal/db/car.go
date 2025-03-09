package db

import (
	"time"
)

func AddCar(c *Car) error {
	return db.Create(&c).Error
}

func FetchCars(
	district string,
	startDate time.Time,
	endDate time.Time,
	query string,
	limit int,
	page int,
) (*[]Car, error) {
	var cars *[]Car
	tx := db.Limit(limit).Offset(page * limit)
	tx.Joins("District").Where("District.name = ?", district)

	if !startDate.IsZero() && !endDate.IsZero() {
		tx.Where("start_at <= ? AND end_at >= ?", startDate, endDate)
	}

	if query != "" {
		query = "%" + query + "%"
		tx.Where("cars.name LIKE ?", query)
	}

	tx.Find(&cars)
	return cars, tx.Error
}
