package db

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"

	_ "github.com/mattn/go-sqlite3"

	. "github.com/untemi/carshift/internal/db/sqlc"
	"github.com/untemi/carshift/internal/misc"
)

//go:embed schema.sql
var tbs string

var (
	runner          *Queries
	ErrNoIdentifier = fmt.Errorf("no Identifier provided")
)

func Init(ctx context.Context) (misc.DBClose, error) {
	conn, err := sql.Open("sqlite3", "app.db")
	if err != nil {
		return nil, err
	}

	runner = New(conn)
	return conn.Close, nil
}

func Setup(ctx context.Context) error {
	conn, err := sql.Open("sqlite3", "app.db")
	if err != nil {
		return err
	}
	_, err = conn.ExecContext(ctx, tbs)

	districts := []string{
		"Aïn Chock",
		"Aïn Sebaâ",
		"Al Fida",
		"Ben M'Sick",
		"Hay Hassani",
		"Moulay Rachid",
		"Sidi Bernoussi",
	}

	tmpRunner := New(conn)
	for _, v := range districts {
		if err := tmpRunner.CreateDistrict(ctx, v); err != nil {
			return err
		}
	}

	conn.Close()
	return err
}
