package misc

import (
	"database/sql"
	"time"
)

func TimeToNull(t time.Time) sql.NullTime {
	if !t.IsZero() {
		return sql.NullTime{Time: t, Valid: true}
	} else {
		return sql.NullTime{Valid: false}
	}
}
