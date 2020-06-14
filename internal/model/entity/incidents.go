package entity

import "database/sql"

type Incidents struct {
	ID       sql.NullInt64
	Status   sql.NullString
	Message  sql.NullString
	Ack      sql.NullString
	Comments []sql.NullString
}
