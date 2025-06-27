package helpers

import (
	"database/sql"
)

func CloseRows(rows *sql.Rows) {
	if rows == nil {
		return
	}
	err := rows.Close()
	if err != nil {
		return
	}
}
