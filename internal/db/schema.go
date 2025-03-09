package db

import "time"

type User struct {
	ID       int    `gorm:"primaryKey,autoIncrement"`
	Username string `gorm:"index,unique"`

	Firstname string
	Lastname  string

	Passhash string

	Phone string
	Email string

	CreatedAt time.Time
	UpdatedAt time.Time
}

type Car struct {
	ID int `gorm:"primaryKey,autoIncrement"`

	Name  string
	Price float64

	UserID int
	User   User

	DistrictID int
	District   District

	StartAt time.Time
	EndAt   time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}

type District struct {
	ID   int    `gorm:"primaryKey"`
	Name string `gorm:"index,unique"`
}
