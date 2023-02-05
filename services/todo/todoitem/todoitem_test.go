package todoitem

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	terr "github.com/stumacwastaken/todo/errors"
)

type MockStorer struct {
	resp func(method string) ([]TodoItem, error)
}

func (m *MockStorer) Create(context.Context, TodoItem) (TodoItem, error) {
	res, err := m.resp("Create")
	if len(res) > 0 {
		return res[0], err
	}
	return TodoItem{}, err
}

func (m *MockStorer) GetAll(ctx context.Context) ([]TodoItem, error) {
	return m.resp("GetAll")

}

func (m *MockStorer) Update(context.Context, TodoItem) (TodoItem, error) {
	res, err := m.resp("Update")
	if len(res) > 0 {
		return res[0], err
	}
	return TodoItem{}, err
}

func (m *MockStorer) GetById(context.Context, string) (TodoItem, error) {
	res, err := m.resp("GetById")
	if len(res) > 0 {
		return res[0], err
	}
	return TodoItem{}, err
}

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

type test struct {
	name       string
	reqItem    TodoItem
	expect     TodoItem
	err        error
	ctx        context.Context
	mockMethod func(method string) ([]TodoItem, error)
}

func TestUpdate(t *testing.T) {
	type test struct {
		id          string
		name        string
		reqItem     TodoItem
		expect      TodoItem
		err         error
		ctx         context.Context
		mockMethod  func(method string) ([]TodoItem, error)
		dateUpdater func() time.Time
	}
	tests := []test{
		{
			name: "failed update unknown error",
			reqItem: TodoItem{
				Id:        newId("3333"),
				Summary:   newSummary("an updated summary"),
				Completed: newBool(true),
			},
			id:          "3333",
			dateUpdater: func() time.Time { return time.Date(2023, time.January, 16, 12, 12, 12, 12, time.Local) },
			mockMethod: func(method string) ([]TodoItem, error) {
				if method == "GetById" {
					return []TodoItem{
						{
							Id:        newId("3333"),
							Summary:   newSummary("a random summary"),
							Completed: newBool(false),
							Deleted:   newBool(false),
							Created:   newTime(time.Date(2023, time.January, 12, 12, 12, 12, 12, time.Local)),
							Updated:   newTime(time.Date(2023, time.January, 13, 12, 12, 12, 12, time.Local)),
						},
					}, nil
				}
				return []TodoItem{}, errors.New("this is a random unknown error from the storage interface")
			},
			expect: TodoItem{},
			err:    terr.InternalError(),
		},
		{
			name: "failed update known error",
			reqItem: TodoItem{
				Id:        newId("3333"),
				Summary:   newSummary("an updated summary"),
				Completed: newBool(true),
			},
			id:          "3333",
			dateUpdater: func() time.Time { return time.Date(2023, time.January, 16, 12, 12, 12, 12, time.Local) },
			mockMethod: func(method string) ([]TodoItem, error) {
				if method == "GetById" {
					return []TodoItem{
						{
							Id:        newId("3333"),
							Summary:   newSummary("a random summary"),
							Completed: newBool(false),
							Deleted:   newBool(false),
							Created:   newTime(time.Date(2023, time.January, 12, 12, 12, 12, 12, time.Local)),
							Updated:   newTime(time.Date(2023, time.January, 13, 12, 12, 12, 12, time.Local)),
						},
					}, nil
				}
				return []TodoItem{}, terr.InternalError()
			},
			expect: TodoItem{},
			err:    terr.InternalError(),
		},
		{
			name: "failed get unknown error",
			reqItem: TodoItem{
				Id:        newId("3333"),
				Summary:   newSummary("an updated summary"),
				Completed: newBool(true),
			},
			id:          "3333",
			dateUpdater: func() time.Time { return time.Date(2023, time.January, 16, 12, 12, 12, 12, time.Local) },
			mockMethod: func(method string) ([]TodoItem, error) {
				if method == "GetById" {
					return []TodoItem{}, errors.New("we have no idea what this error is")
				}
				return []TodoItem{
					{
						Id:        newId("3333"),
						Summary:   newSummary("an updated summary"),
						Completed: newBool(true),
						Deleted:   newBool(false),
						Created:   newTime(time.Date(2023, time.January, 12, 12, 12, 12, 12, time.Local)),
						Updated:   newTime(time.Date(2023, time.January, 16, 12, 12, 12, 12, time.Local)),
					},
				}, nil
			},
			expect: TodoItem{},
			err:    terr.InternalError(),
		},
		{
			name: "failed get known error",
			reqItem: TodoItem{
				Id:        newId("3333"),
				Summary:   newSummary("an updated summary"),
				Completed: newBool(true),
			},
			id:          "3333",
			dateUpdater: func() time.Time { return time.Date(2023, time.January, 16, 12, 12, 12, 12, time.Local) },
			mockMethod: func(method string) ([]TodoItem, error) {
				if method == "GetById" {
					return []TodoItem{}, terr.InternalError()
				}
				return []TodoItem{
					{
						Id:        newId("3333"),
						Summary:   newSummary("an updated summary"),
						Completed: newBool(true),
						Deleted:   newBool(false),
						Created:   newTime(time.Date(2023, time.January, 12, 12, 12, 12, 12, time.Local)),
						Updated:   newTime(time.Date(2023, time.January, 16, 12, 12, 12, 12, time.Local)),
					},
				}, nil
			},
			expect: TodoItem{},
			err:    terr.InternalError(),
		},
		{
			name:    "no summary",
			id:      "3333",
			reqItem: TodoItem{Summary: nil},
			expect:  TodoItem{},
			err:     terr.ErrorWithCode("bad request", "cannot have empty summary", 400),
		},
		{
			name:    "empty summary",
			id:      "3333",
			reqItem: TodoItem{Summary: newSummary("")},
			expect:  TodoItem{},
			err:     terr.ErrorWithCode("bad request", "cannot have empty summary", 400),
		},
		{
			name:    "no id",
			id:      "",
			reqItem: TodoItem{},
			expect:  TodoItem{},
			err:     terr.ErrorWithCode("no id", "no id found in request", 404),
		},
		{
			name: "happy path deleted",
			reqItem: TodoItem{
				Id:      newId("3333"),
				Summary: newSummary("an updated summary"),
				Deleted: newBool(true),
			},
			id:          "3333",
			dateUpdater: func() time.Time { return time.Date(2023, time.January, 16, 12, 12, 12, 12, time.Local) },
			mockMethod: func(method string) ([]TodoItem, error) {
				if method == "GetById" {
					return []TodoItem{
						{
							Id:        newId("3333"),
							Summary:   newSummary("a random summary"),
							Completed: newBool(false),
							Deleted:   newBool(false),
							Created:   newTime(time.Date(2023, time.January, 12, 12, 12, 12, 12, time.Local)),
							Updated:   newTime(time.Date(2023, time.January, 13, 12, 12, 12, 12, time.Local)),
						},
					}, nil
				}
				return []TodoItem{
					{
						Id:        newId("3333"),
						Summary:   newSummary("an updated summary"),
						Completed: newBool(false),
						Deleted:   newBool(true),
						Created:   newTime(time.Date(2023, time.January, 12, 12, 12, 12, 12, time.Local)),
						Updated:   newTime(time.Date(2023, time.January, 16, 12, 12, 12, 12, time.Local)),
					},
				}, nil
			},
			expect: TodoItem{
				Id:        newId("3333"),
				Summary:   newSummary("an updated summary"),
				Completed: newBool(false),
				Deleted:   newBool(true),
				Created:   newTime(time.Date(2023, time.January, 12, 12, 12, 12, 12, time.Local)),
				Updated:   newTime(time.Date(2023, time.January, 16, 12, 12, 12, 12, time.Local)),
			},
		},
		{
			name: "happy path",
			reqItem: TodoItem{
				Id:        newId("3333"),
				Summary:   newSummary("an updated summary"),
				Completed: newBool(true),
			},
			id:          "3333",
			dateUpdater: func() time.Time { return time.Date(2023, time.January, 16, 12, 12, 12, 12, time.Local) },
			mockMethod: func(method string) ([]TodoItem, error) {
				if method == "GetById" {
					return []TodoItem{
						{
							Id:        newId("3333"),
							Summary:   newSummary("a random summary"),
							Completed: newBool(false),
							Deleted:   newBool(false),
							Created:   newTime(time.Date(2023, time.January, 12, 12, 12, 12, 12, time.Local)),
							Updated:   newTime(time.Date(2023, time.January, 13, 12, 12, 12, 12, time.Local)),
						},
					}, nil
				}
				return []TodoItem{
					{
						Id:        newId("3333"),
						Summary:   newSummary("an updated summary"),
						Completed: newBool(true),
						Deleted:   newBool(false),
						Created:   newTime(time.Date(2023, time.January, 12, 12, 12, 12, 12, time.Local)),
						Updated:   newTime(time.Date(2023, time.January, 16, 12, 12, 12, 12, time.Local)),
					},
				}, nil
			},
			expect: TodoItem{
				Id:        newId("3333"),
				Summary:   newSummary("an updated summary"),
				Completed: newBool(true),
				Deleted:   newBool(false),
				Created:   newTime(time.Date(2023, time.January, 12, 12, 12, 12, 12, time.Local)),
				Updated:   newTime(time.Date(2023, time.January, 16, 12, 12, 12, 12, time.Local)),
			},
		},
	}

	for _, tt := range tests {
		tf := func(t *testing.T) {
			mocks := &MockStorer{}
			subject := NewCore(mocks)
			mocks.resp = tt.mockMethod
			res, err := subject.Update(tt.ctx, tt.reqItem, tt.id)

			assert.Equal(t, tt.err, err, "errors should match")
			assert.Equal(t, tt.expect, res)

		}
		t.Run(tt.name, tf)
	}
}

func TestGetAll(t *testing.T) {
	type test struct {
		name       string
		reqItem    TodoItem
		expect     []TodoItem
		err        error
		ctx        context.Context
		mockMethod func(method string) ([]TodoItem, error)
	}
	tests := []test{
		{
			name:   "error on retrieval from store",
			expect: []TodoItem{},
			err:    terr.UnknownError(),
			mockMethod: func(method string) ([]TodoItem, error) {
				return []TodoItem{}, terr.UnknownError()
			},
			ctx: context.Background(),
		},
		{
			name: "happy path",
			expect: []TodoItem{
				{
					Summary:   newSummary("test summary"),
					Id:        newId("112233445566"),
					Created:   newTime(time.Date(2023, time.January, 12, 12, 12, 12, 12, time.Local)),
					Updated:   newTime(time.Date(2023, time.January, 13, 12, 12, 12, 12, time.Local)),
					Completed: newBool(false),
					Deleted:   newBool(false),
				},
			},
			ctx: context.Background(),
			mockMethod: func(method string) ([]TodoItem, error) {
				return []TodoItem{
					{
						Summary:   newSummary("test summary"),
						Id:        newId("112233445566"),
						Created:   newTime(time.Date(2023, time.January, 12, 12, 12, 12, 12, time.Local)),
						Updated:   newTime(time.Date(2023, time.January, 13, 12, 12, 12, 12, time.Local)),
						Completed: newBool(false),
						Deleted:   newBool(false),
					},
				}, nil
			},
		},
	}
	for _, tt := range tests {
		tf := func(t *testing.T) {
			mocks := &MockStorer{}
			subject := NewCore(mocks)
			mocks.resp = tt.mockMethod
			res, err := subject.GetAll(tt.ctx)

			assert.Equal(t, tt.err, err, "errors should match")
			assert.Equal(t, tt.expect, res)

		}
		t.Run(tt.name, tf)
	}
}
func TestCreate(t *testing.T) {
	type test struct {
		name       string
		reqItem    TodoItem
		expect     TodoItem
		err        error
		ctx        context.Context
		mockMethod func(method string) ([]TodoItem, error)
	}
	tests := []test{
		{
			name: "id not nil",
			reqItem: TodoItem{
				Summary: newSummary("summary string"),
				Id:      newId("3434343"),
			},
			expect: TodoItem{},
			err:    terr.ErrorWithCode("invalid param", "cannot create a todo item with an already existing id", 400),
		},
		{
			name: "summary nil",
			reqItem: TodoItem{
				Summary: nil,
			},
			expect: TodoItem{},
			err:    terr.ErrorWithCode("invalid param", "summary cannot be empty", 400),
		},
		{
			name: "summary empty",
			reqItem: TodoItem{
				Summary: newSummary(""),
			},
			expect: TodoItem{},
			err:    terr.ErrorWithCode("invalid param", "summary cannot be empty", 400),
		},
		{
			name: "happy path",
			mockMethod: func(method string) ([]TodoItem, error) {
				return []TodoItem{
					{
						Summary:   newSummary("test summary"),
						Id:        newId("112233445566"),
						Created:   newTime(time.Date(2023, time.January, 12, 12, 12, 12, 12, time.Local)),
						Updated:   newTime(time.Date(2023, time.January, 12, 12, 12, 12, 12, time.Local)),
						Completed: newBool(false),
						Deleted:   newBool(false),
					},
				}, nil
			},
			reqItem: TodoItem{Summary: newSummary("test summary")},
			expect: TodoItem{
				Summary:   newSummary("test summary"),
				Id:        newId("112233445566"),
				Created:   newTime(time.Date(2023, time.January, 12, 12, 12, 12, 12, time.Local)),
				Updated:   newTime(time.Date(2023, time.January, 12, 12, 12, 12, 12, time.Local)),
				Completed: newBool(false),
				Deleted:   newBool(false),
			},
			err: nil,
		},
	}
	for _, tt := range tests {
		tf := func(t *testing.T) {
			mocks := &MockStorer{}
			subject := NewCore(mocks)
			mocks.resp = tt.mockMethod
			res, err := subject.Create(tt.ctx, tt.reqItem)

			assert.Equal(t, tt.err, err, "errors should match")
			assert.Equal(t, tt.expect, res)

		}
		t.Run(tt.name, tf)
	}

}
