package db

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"

	_ "github.com/mattn/go-sqlite3"

	. "github.com/untemi/carshift/internal/db/sqlc"
)

type closeFunc func() error

//go:embed schema.sql
var tbs string

var (
	runner          *Queries
	ErrNoIdentifier = fmt.Errorf("no Identifier provided")
)

func Init(ctx context.Context) (closeFunc, error) {
	conn, err := sql.Open("sqlite3", "app.db")
	if err != nil {
		return nil, err
	}

	// if err := setup(ctx, conn); err != nil {
	// 	return nil, err
	// }

	runner = New(conn)
	return conn.Close, nil
}

func setup(ctx context.Context, conn *sql.DB) error {
	_, err := conn.ExecContext(ctx, tbs)
	return err
}
