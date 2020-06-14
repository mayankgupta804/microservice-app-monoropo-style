package entity

import "database/sql"

type Comments struct {
	ID         sql.NullInt64
	IncidentID sql.NullInt64
	Comment    sql.NullString
}
