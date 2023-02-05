package tododb

import "time"

type dbTodoItem struct {
	Id          string    `db:"id"`
	Summary     string    `db:"summary"`
	DateCreated time.Time `db:"date_created"`
	DateUpdated time.Time `db:"date_updated"`
	Deleted     bool      `db:"deleted"`
	Completed   bool      `db:"completed"`
}
