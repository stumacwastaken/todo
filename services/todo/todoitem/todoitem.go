package todoitem

import (
	"context"
	"time"

	terr "github.com/stumacwastaken/todo/errors"
)

type Storer interface {
	Create(context.Context, TodoItem) (TodoItem, error)
	GetAll(ctx context.Context) ([]TodoItem, error) //could be improved to return additional metadata/page/row/etc
	Update(context.Context, TodoItem) (TodoItem, error)
	GetById(context.Context, string) (TodoItem, error)
}

type Core struct {
	storer Storer
}

func NewCore(storer Storer) *Core {
	return &Core{
		storer: storer,
	}
}

//pull out so we can change give a custom time at testing.
var dateUpdateFn = time.Now

//Creates and Inserts a new Todo item into the database after basic validation
func (c *Core) Create(ctx context.Context, newTodo TodoItem) (TodoItem, error) {
	if newTodo.Id != nil {
		return TodoItem{}, terr.ErrorWithCode("invalid param", "cannot create a todo item with an already existing id", 400)
	}
	if newTodo.Summary == nil || *newTodo.Summary == "" {
		return TodoItem{}, terr.ErrorWithCode("invalid param", "summary cannot be empty", 400)
	}
	return c.storer.Create(ctx, newTodo)
}

func (c *Core) Update(ctx context.Context, newItem TodoItem, id string) (TodoItem, error) {
	if id == "" {
		//should never really get here from restful api
		return TodoItem{}, terr.ErrorWithCode("no id", "no id found in request", 404)
	}
	if newItem.Summary == nil || *newItem.Summary == "" {
		return TodoItem{}, terr.ErrorWithCode("bad request", "cannot have empty summary", 400)
	}
	//getItem
	oldItem, err := c.storer.GetById(ctx, id)
	if err != nil {
		if v, ok := err.(*terr.TodoError); ok {
			return TodoItem{}, v
		} else {
			return TodoItem{}, terr.InternalError()
		}
	}
	//merge new into old
	toSave := mergeItems(oldItem, newItem)
	t := dateUpdateFn()
	toSave.Updated = &t

	//update
	saved, err := c.storer.Update(ctx, toSave)
	if err != nil {
		if v, ok := err.(*terr.TodoError); ok {
			return TodoItem{}, v
		} else {
			return TodoItem{}, terr.InternalError()
		}
	}
	return saved, nil
}

func (c *Core) GetAll(ctx context.Context) ([]TodoItem, error) {
	return c.storer.GetAll(ctx)
}

func mergeItems(old, new TodoItem) TodoItem {
	if new.Completed != nil {
		old.Completed = new.Completed
	}
	if new.Deleted != nil {
		old.Deleted = new.Deleted
	}
	if new.Summary != nil {
		old.Summary = new.Summary
	}
	return old
}
