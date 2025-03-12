// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: query.sql

package sqlc

import (
	"context"
	"database/sql"
)

const createCar = `-- name: CreateCar :one
INSERT INTO cars (
  name, price, start_at, end_at, owner_id, district_id 
) VALUES (
  ?, ?, ?, ?, ?, ?
)
RETURNING id
`

type CreateCarParams struct {
	Name       string
	Price      float64
	StartAt    sql.NullTime
	EndAt      sql.NullTime
	OwnerID    int64
	DistrictID int64
}

// Cars     --
func (q *Queries) CreateCar(ctx context.Context, arg CreateCarParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, createCar,
		arg.Name,
		arg.Price,
		arg.StartAt,
		arg.EndAt,
		arg.OwnerID,
		arg.DistrictID,
	)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const createUser = `-- name: CreateUser :one
INSERT INTO users (
  username, firstname, lastname, passhash, phone, email 
) VALUES (
  ?, ?, ?, ?, ?, ?
)
RETURNING id
`

type CreateUserParams struct {
	Username  string
	Firstname string
	Lastname  string
	Passhash  string
	Phone     string
	Email     string
}

// Users     --
func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.Username,
		arg.Firstname,
		arg.Lastname,
		arg.Passhash,
		arg.Phone,
		arg.Email,
	)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const findUserById = `-- name: FindUserById :one
SELECT id, username, firstname, lastname, passhash, phone, email FROM users
  WHERE id = ?
LIMIT 1
`

func (q *Queries) FindUserById(ctx context.Context, id int64) (User, error) {
	row := q.db.QueryRowContext(ctx, findUserById, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Firstname,
		&i.Lastname,
		&i.Passhash,
		&i.Phone,
		&i.Email,
	)
	return i, err
}

const findUserByUsername = `-- name: FindUserByUsername :one
SELECT id, username, firstname, lastname, passhash, phone, email FROM users
  WHERE username = ?
LIMIT 1
`

func (q *Queries) FindUserByUsername(ctx context.Context, username string) (User, error) {
	row := q.db.QueryRowContext(ctx, findUserByUsername, username)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Firstname,
		&i.Lastname,
		&i.Passhash,
		&i.Phone,
		&i.Email,
	)
	return i, err
}

const isUsernameUsed = `-- name: IsUsernameUsed :one
SELECT COUNT(id) FROM users
  WHERE username = ?
LIMIT 1
`

func (q *Queries) IsUsernameUsed(ctx context.Context, username string) (int64, error) {
	row := q.db.QueryRowContext(ctx, isUsernameUsed, username)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const queryCars = `-- name: QueryCars :many
SELECT c.id, c.name, c.price, c.start_at, c.end_at, c.owner_id, c.district_id FROM cars c
  INNER JOIN districts d ON c.district_id = d.id
WHERE c.name LIKE ?1
  AND d.name = ?6
  AND (
    (start_at <= ?2 AND end_at >= ?3) 
    OR (?2 IS NULL AND ?3 IS NULL)
  )
LIMIT ?4 OFFSET ?5
`

type QueryCarsParams struct {
	Name         string
	StartAt      sql.NullTime
	EndAt        sql.NullTime
	Limit        int64
	Offset       int64
	DistrictName string
}

func (q *Queries) QueryCars(ctx context.Context, arg QueryCarsParams) ([]Car, error) {
	rows, err := q.db.QueryContext(ctx, queryCars,
		arg.Name,
		arg.StartAt,
		arg.EndAt,
		arg.Limit,
		arg.Offset,
		arg.DistrictName,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Car
	for rows.Next() {
		var i Car
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Price,
			&i.StartAt,
			&i.EndAt,
			&i.OwnerID,
			&i.DistrictID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const queryUsers = `-- name: QueryUsers :many
SELECT id, username, firstname, lastname, passhash, phone, email FROM users
  WHERE username LIKE ?
LIMIT ? OFFSET ?
`

type QueryUsersParams struct {
	Username string
	Limit    int64
	Offset   int64
}

func (q *Queries) QueryUsers(ctx context.Context, arg QueryUsersParams) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, queryUsers, arg.Username, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Username,
			&i.Firstname,
			&i.Lastname,
			&i.Passhash,
			&i.Phone,
			&i.Email,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
