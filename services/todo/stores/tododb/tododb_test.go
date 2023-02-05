package tododb

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	terr "github.com/stumacwastaken/todo/errors"
	"github.com/stumacwastaken/todo/todoitem"
)

func newId(id string) *string {
	return &id
}
func newSummary(summary string) *string {
	return &summary
}
func newTime(ti time.Time) *time.Time {
	return &ti
}
func newBool(b bool) *bool {
	return &b
}

var testTime = newTime(time.Date(2023, time.January, 12, 12, 12, 12, 12, time.Local))

func TestUpdate(t *testing.T) {

}
func TestGetAll(t *testing.T) {
	var rows = sqlmock.NewRows([]string{"id", "summary", "date_created", "date_updated", "completed", "deleted"})
	type test struct {
		name      string
		expect    []todoitem.TodoItem
		expectErr error
		mockRows  *sqlmock.Rows
		mockErr   error
	}
	tests := []test{
		{
			name:      "unknown",
			expect:    []todoitem.TodoItem{},
			expectErr: terr.ErrorWithCode("internal error", "Could not query for todos", 500),
			mockRows:  nil,
			mockErr:   errors.New("a random sql test error"),
		},
		{
			name:      "no rows",
			expect:    []todoitem.TodoItem{},
			expectErr: nil,
			mockRows:  nil,
			mockErr:   sql.ErrNoRows,
		},
		{
			name: "happy path",
			expect: []todoitem.TodoItem{
				{
					Id:        newId("1111"),
					Created:   testTime,
					Updated:   testTime,
					Deleted:   newBool(false),
					Completed: newBool(false),
					Summary:   newSummary("test summary"),
				},
			},
			expectErr: nil,
			mockRows:  rows.AddRow("1111", "test summary", testTime, testTime, false, false),
			mockErr:   nil,
		},
	}
	for _, tt := range tests {
		tf := func(t *testing.T) {
			mockDB, mock, err := sqlmock.New()
			if err != nil {
				t.Fatal(err)
			}
			defer mockDB.Close()
			db := sqlx.NewDb(mockDB, "sqlmock")
			store := NewStore(db)
			// testTime := newTime(time.Date(2023, time.January, 12, 12, 12, 12, 12, time.Local))
			mock.ExpectBegin()
			query := mock.ExpectQuery(`SELECT \* FROM todo_item WHERE deleted=false ORDER BY date_created DESC`).WithArgs()

			if tt.mockErr != nil {
				query.WillReturnError(tt.mockErr)
			} else {
				query.WillReturnRows(tt.mockRows)
			}

			val, err := store.GetAll(context.Background())
			mock.ExpectCommit()
			assert.Equal(t, tt.expect, val)
			assert.Equal(t, tt.expectErr, err)
		}
		t.Run(tt.name, tf)
	}
}

func TestGetById(t *testing.T) {
	var rows = sqlmock.NewRows([]string{"id", "summary", "date_created", "date_updated", "completed", "deleted"})
	type test struct {
		name      string
		expect    todoitem.TodoItem
		expectErr error
		mockRows  *sqlmock.Rows
		mockErr   error
	}
	tests := []test{
		{
			name: "happy path",
			expect: todoitem.TodoItem{
				Id:        newId("1111"),
				Created:   testTime,
				Updated:   testTime,
				Deleted:   newBool(false),
				Completed: newBool(false),
				Summary:   newSummary("test summary"),
			},
			expectErr: nil,
			mockRows:  rows.AddRow("1111", "test summary", testTime, testTime, false, false),
			mockErr:   nil,
		},
		{
			name:      "no rows found",
			expect:    todoitem.TodoItem{},
			expectErr: terr.ErrorWithCode("not found", "Item with id 1111 not found", 404),
			mockErr:   sql.ErrNoRows,
		},
		{
			name:      "unknown error",
			expect:    todoitem.TodoItem{},
			expectErr: terr.UnknownError(),
			mockErr:   errors.New("some random mysql error"),
		},
	}
	for _, tt := range tests {
		tf := func(t *testing.T) {
			mockDB, mock, err := sqlmock.New()
			if err != nil {
				t.Fatal(err)
			}
			defer mockDB.Close()
			db := sqlx.NewDb(mockDB, "sqlmock")
			store := NewStore(db)
			// testTime := newTime(time.Date(2023, time.January, 12, 12, 12, 12, 12, time.Local))
			query := mock.ExpectQuery(`SELECT \* FROM todo_item where id=\?`).WithArgs("1111")
			if tt.mockErr != nil {
				query.WillReturnError(tt.mockErr)
			} else {
				query.WillReturnRows(tt.mockRows)
			}
			val, err := store.GetById(context.Background(), "1111")
			assert.Equal(t, tt.expect, val)
			assert.Equal(t, tt.expectErr, err)
		}
		t.Run(tt.name, tf)
	}
}
