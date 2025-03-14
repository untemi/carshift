package db

import (
	"context"

	. "github.com/untemi/carshift/internal/db/sqlc"
)

func RegisterUser(ctx context.Context, u *User) (int64, error) {
	return runner.CreateUser(
		ctx,
		CreateUserParams{
			Username:  u.Username,
			Firstname: u.Firstname,
			Lastname:  u.Lastname,
			Passhash:  u.Passhash,
		})
}

func IsUsernameUsed(ctx context.Context, username string) (bool, error) {
	e, err := runner.IsUsernameUsed(ctx, username)
	return e == 1, err
}

func FillUser(ctx context.Context, u *User) error {
	var err error
	if u.Username != "" {
		*u, err = runner.FindUserByUsername(ctx, u.Username)
		return err
	} else if u.ID != 0 {
		*u, err = runner.FindUserById(ctx, u.ID)
		return err
	}
	return ErrNoIdentifier
}

func FetchUsers(ctx context.Context, query string, limit int64, page int64) (*[]User, error) {
	users, err := runner.QueryUsers(
		ctx,
		QueryUsersParams{
			Username: "%" + query + "%",
			Limit:    limit,
			Offset:   page,
		},
	)

	return &users, err
}

func UpdateUser(ctx context.Context, u *User) error {
	return runner.UpdateUser(
		ctx,
		UpdateUserParams{
			Username:  u.Username,
			Firstname: u.Firstname,
			Lastname:  u.Lastname,
			Passhash:  u.Passhash,
			PfpName:   u.PfpName,
			Phone:     u.Phone,
			Email:     u.Email,
			ID:        u.ID,
		},
	)
}
