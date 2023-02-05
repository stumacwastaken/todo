package tododb

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/stumacwastaken/todo/log"
	"go.uber.org/zap"

	"github.com/jmoiron/sqlx"
	"github.com/stumacwastaken/todo/errors"
	"github.com/stumacwastaken/todo/todoitem"
)

type Store struct {
	db *sqlx.DB
}

func NewStore(db *sqlx.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) Create(ctx context.Context, item todoitem.TodoItem) (todoitem.TodoItem, error) {
	statement := `INSERT into todo_item (summary) VALUES (?)`
	tx, err := s.db.Beginx()
	if err != nil {
		log.Default().Error("failed to start transaction", zap.Error(err))
		return todoitem.TodoItem{}, errors.ErrorWithCode("internal error", "Could not query for todos", 500)
	}
	res, err := tx.Exec(statement, item.Summary)
	if err != nil {
		log.Default().Warn("error creating new todo item in database", zap.Error(err))
		tx.Rollback()
		return todoitem.TodoItem{}, errors.UnknownError()
	}
	v := new(dbTodoItem)

	if err != nil {
		log.Default().Warn("unknown error inserting row into database", zap.Error(err))

		return todoitem.TodoItem{}, errors.UnknownError()
	}
	id, _ := res.LastInsertId()
	num, _ := res.RowsAffected()
	log.Default().Info("inserted new todo item", zap.Int64("rows-affected", num), zap.Int64("lastId", id))

	//hack to get the last element because mysql isn't postgres
	getStatement := `SELECT * from todo_item ORDER BY date_created desc LIMIT 1`
	row := tx.QueryRowxContext(ctx, getStatement)
	if row.Err() != nil {
		tx.Rollback()
		log.Default().Warn("unknown error inserting row into database", zap.Error(err))
		return todoitem.TodoItem{}, errors.UnknownError()
	}
	if err := row.StructScan(v); err != nil {
		tx.Rollback()
		log.Default().Warn("unknown error inserting row into database", zap.Error(err))
		return todoitem.TodoItem{}, errors.UnknownError()
	}

	tx.Commit()
	return toCoreItem(*v), nil
}

func (s *Store) Update(ctx context.Context, item todoitem.TodoItem) (todoitem.TodoItem, error) {
	statement := `UPDATE todo_item SET summary = ?, date_updated = ?, deleted = ?, completed = ? WHERE id = ?`
	tx, err := s.db.Beginx()
	if err != nil {
		log.Default().Error("failed to start transaction", zap.Error(err))
		return todoitem.TodoItem{}, errors.ErrorWithCode("internal error", "Could not query for todos", 500)
	}

	row := tx.QueryRowx(statement, item.Summary, item.Updated, item.Deleted, item.Completed, item.Id)
	err = row.Err()
	if err != nil {
		tx.Rollback()
		if err == sql.ErrNoRows {

			log.Default().Error("error no rows on an id that's supposed to be there. How did you get here?", zap.Error(err), zap.String("id", *item.Id))
			return todoitem.TodoItem{}, errors.ErrorWithCode("not found", fmt.Sprintf("Item with id %s not found", *item.Id), 404)
		}
		log.Default().Error("error updating row", zap.Error(err), zap.String("id", *item.Id))
		return todoitem.TodoItem{}, errors.UnknownError()
	}
	tx.Commit()

	return item, nil
}

func (s *Store) GetById(ctx context.Context, id string) (todoitem.TodoItem, error) {
	statement := "SELECT * FROM todo_item where id=?"
	row := s.db.QueryRowx(statement, id)
	v := new(dbTodoItem)
	err := row.StructScan(v)
	if err != nil {
		if err == sql.ErrNoRows {
			return todoitem.TodoItem{}, errors.ErrorWithCode("not found", fmt.Sprintf("Item with id %s not found", id), 404)
		} else {
			log.Default().Error("unknown error querying todo by id", zap.Error(err), zap.String("req id", id))
			return todoitem.TodoItem{}, errors.UnknownError()
		}
	}

	return toCoreItem(*v), nil
}

//could be improved to return additional metadata and better query filtering
func (s *Store) GetAll(ctx context.Context) ([]todoitem.TodoItem, error) {
	q := fmt.Sprintf(`SELECT * FROM todo_item WHERE deleted=false ORDER BY date_created DESC`) //keep fmt here for now....cause queries and filters
	tx, err := s.db.Beginx()
	if err != nil {
		log.Default().Error("failed to start transaction", zap.Error(err))
		return nil, errors.ErrorWithCode("internal error", "Could not query for todos", 500)
	}
	rows, err := tx.QueryxContext(ctx, q)
	if err != nil {
		if err == sql.ErrNoRows {
			tx.Rollback()
			return []todoitem.TodoItem{}, nil
		}
		log.Default().Debug("database query failed", zap.Error(err))
		tx.Rollback()
		return []todoitem.TodoItem{}, errors.ErrorWithCode("internal error", "Could not query for todos", 500)
	}

	var dbItems []dbTodoItem
	for rows.Next() {
		v := new(dbTodoItem)
		if err := rows.StructScan(v); err != nil {
			return nil, err
		}
		dbItems = append(dbItems, *v)
	}
	tx.Commit()
	return toCoreTodoSlice(dbItems), nil
}

func toCoreTodoSlice(dbTodoItems []dbTodoItem) []todoitem.TodoItem {
	var coreItems []todoitem.TodoItem
	for _, item := range dbTodoItems {
		coreItems = append(coreItems, toCoreItem(item))
	}
	return coreItems
}

func toCoreItem(item dbTodoItem) todoitem.TodoItem {
	coreTodoItem := todoitem.TodoItem{
		Id:        &item.Id,
		Created:   &item.DateCreated,
		Updated:   &item.DateUpdated,
		Completed: &item.Completed,
		Deleted:   &item.Deleted,
		Summary:   &item.Summary,
	}
	return coreTodoItem
}
