package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	terr "github.com/stumacwastaken/todo/errors"
	"github.com/stumacwastaken/todo/todoitem"
)

func TestNewTodoHandler(t *testing.T) {
	handler := NewTodoHandlers(nil)
	assert.NotNil(t, handler)
}

type MockStorer struct {
	resp func(method string) ([]todoitem.TodoItem, error)
}

func (m *MockStorer) Create(context.Context, todoitem.TodoItem) (todoitem.TodoItem, error) {
	res, err := m.resp("Create")
	if len(res) > 0 {
		return res[0], err
	}
	return todoitem.TodoItem{}, err
}

func (m *MockStorer) GetAll(ctx context.Context) ([]todoitem.TodoItem, error) {
	return m.resp("GetAll")

}

func (m *MockStorer) Update(context.Context, todoitem.TodoItem) (todoitem.TodoItem, error) {
	res, err := m.resp("Update")
	if len(res) > 0 {
		return res[0], err
	}
	return todoitem.TodoItem{}, err
}

func (m *MockStorer) GetById(context.Context, string) (todoitem.TodoItem, error) {
	res, err := m.resp("GetById")
	if len(res) > 0 {
		return res[0], err
	}
	return todoitem.TodoItem{}, err
}

func NewCore(m *MockStorer) *todoitem.Core {
	return todoitem.NewCore(m)
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

type expectErr struct {
	statusCode   int
	bodyContains string
}
type test struct {
	name       string
	mockMethod func(method string) ([]todoitem.TodoItem, error)
	expect     []todoitem.TodoItem
	expectErr  *expectErr
}

func TestGetTodos(t *testing.T) {

	tests := []test{
		{
			name: "HappyGetAll",
			mockMethod: func(method string) ([]todoitem.TodoItem, error) {
				testTodo := todoitem.TodoItem{
					Id:      newId("343434"),
					Created: newTime(time.Now()),
					Updated: newTime(time.Now()),
					Deleted: newBool(false), Completed: newBool(true),
					Summary: newSummary("test summary"),
				}
				return []todoitem.TodoItem{testTodo}, nil
			},
			expect: []todoitem.TodoItem{{
				Id:      newId("343434"),
				Created: newTime(time.Now()),
				Updated: newTime(time.Now()),
				Deleted: newBool(false), Completed: newBool(true),
				Summary: newSummary("test summary"),
			}},
			expectErr: nil,
		},
		{
			name: "Unknown error",
			mockMethod: func(method string) ([]todoitem.TodoItem, error) {
				return []todoitem.TodoItem{}, terr.UnknownError()
			},
			expect:    []todoitem.TodoItem{},
			expectErr: &expectErr{500, "unknown"},
		},
		{
			name: "internal",
			mockMethod: func(method string) ([]todoitem.TodoItem, error) {
				return []todoitem.TodoItem{}, errors.New("test unknown error")
			},
			expect:    []todoitem.TodoItem{},
			expectErr: &expectErr{500, "internal"},
		},
	}
	for _, tt := range tests {
		tf := func(t *testing.T) {
			parent := chi.NewRouter()
			mocks := &MockStorer{}
			subject := NewTodoHandlers(NewCore(mocks))
			subject.RegisterTodoEndpoints(parent, "/api")
			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/api/todo", nil)

			mocks.resp = tt.mockMethod
			parent.ServeHTTP(rr, req)

			if tt.expectErr != nil {
				assert.Equal(t, tt.expectErr.statusCode, rr.Result().StatusCode, "Should have correct status code")
				defer rr.Result().Body.Close()
				b, _ := io.ReadAll(rr.Result().Body)
				assert.Contains(t, string(b), tt.expectErr.bodyContains, "error code should contain reference to details")
			} else {
				var resitem []todoitem.TodoItem
				if err := json.NewDecoder(rr.Body).Decode(&resitem); err != nil {
					assert.Fail(t, "failed to decode body", err)
				}
				assert.EqualValues(t, tt.expect[0].Summary, resitem[0].Summary, "response body should marshal to the expected item array")
			}
		}
		t.Run(tt.name, tf)
	}

}

func TestUpdateTodo(t *testing.T) {
	type putTest struct {
		test
		reqBody todoitem.TodoItem
	}
	tests := []putTest{
		{
			reqBody: todoitem.TodoItem{
				Id:        newId("343434"),
				Completed: newBool(true),
				Summary:   newSummary("updated test summary"),
			},
			test: test{
				name: "HappyPost",
				mockMethod: func(method string) ([]todoitem.TodoItem, error) {
					if method == "GetById" {
						testTodo := todoitem.TodoItem{
							Id:        newId("343434"),
							Completed: newBool(false),
							Summary:   newSummary("an uncompleted test summary"),
						}
						return []todoitem.TodoItem{testTodo}, nil
					} else if method == "Update" {

					}
					testTodo := todoitem.TodoItem{
						Id:      newId("343434"),
						Created: newTime(time.Now()),
						Updated: newTime(time.Now()),
						Deleted: newBool(false), Completed: newBool(true),
						Summary: newSummary("updated test summary"),
					}
					return []todoitem.TodoItem{testTodo}, nil
				},

				expect: []todoitem.TodoItem{{
					Id:      newId("343434"),
					Created: newTime(time.Now()),
					Updated: newTime(time.Now()),
					Deleted: newBool(false), Completed: newBool(true),
					Summary: newSummary("updated test summary"),
				}},
				expectErr: nil,
			},
		},
		{
			reqBody: todoitem.TodoItem{Summary: newSummary("TESTING ALL THE THINGS"), Id: newId("343434")},
			test: test{
				name: "Unknown error",
				mockMethod: func(method string) ([]todoitem.TodoItem, error) {
					return []todoitem.TodoItem{}, terr.UnknownError()
				},
				expect:    []todoitem.TodoItem{},
				expectErr: &expectErr{500, "unknown"},
			},
		},
		{
			reqBody: todoitem.TodoItem{Summary: newSummary("TESTING ALL THE THINGS"), Id: newId("343434")},
			test: test{
				name: "internal",
				mockMethod: func(method string) ([]todoitem.TodoItem, error) {
					return []todoitem.TodoItem{}, errors.New("test unknown error")
				},
				expect:    []todoitem.TodoItem{},
				expectErr: &expectErr{500, "internal"},
			},
		},
	}
	for _, tt := range tests {
		tf := func(t *testing.T) {
			parent := chi.NewRouter()
			mocks := &MockStorer{}
			subject := NewTodoHandlers(NewCore(mocks))
			subject.RegisterTodoEndpoints(parent, "/api")
			rr := httptest.NewRecorder()
			body, _ := json.Marshal(tt.reqBody)
			req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/api/todo/%s", *tt.reqBody.Id), bytes.NewReader(body))

			mocks.resp = tt.mockMethod
			parent.ServeHTTP(rr, req)

			if tt.expectErr != nil {
				assert.Equal(t, tt.expectErr.statusCode, rr.Result().StatusCode, "Should have correct status code")
				defer rr.Result().Body.Close()
				b, _ := io.ReadAll(rr.Result().Body)
				assert.Contains(t, string(b), tt.expectErr.bodyContains, "error code should contain reference to details")
			} else {
				var resitem todoitem.TodoItem
				if err := json.NewDecoder(rr.Body).Decode(&resitem); err != nil {
					assert.Fail(t, "failed to decode body", err)
				}
				assert.EqualValues(t, tt.expect[0].Summary, resitem.Summary, "response body should marshal to the expected item")
			}
		}
		t.Run(tt.name, tf)
	}
}

func TestCreateTodo(t *testing.T) {
	type postTest struct {
		test
		reqBody todoitem.TodoItem
	}
	tests := []postTest{
		{
			reqBody: todoitem.TodoItem{
				Summary: newSummary("test summary"),
			},
			test: test{
				name: "HappyPost",
				mockMethod: func(method string) ([]todoitem.TodoItem, error) {
					testTodo := todoitem.TodoItem{
						Id:      newId("343434"),
						Created: newTime(time.Now()),
						Updated: newTime(time.Now()),
						Deleted: newBool(false), Completed: newBool(true),
						Summary: newSummary("test summary"),
					}
					return []todoitem.TodoItem{testTodo}, nil
				},

				expect: []todoitem.TodoItem{{
					Id:      newId("343434"),
					Created: newTime(time.Now()),
					Updated: newTime(time.Now()),
					Deleted: newBool(false), Completed: newBool(true),
					Summary: newSummary("test summary"),
				}},
				expectErr: nil,
			},
		},
		{
			reqBody: todoitem.TodoItem{Summary: newSummary("TESTING ALL THE THINGS")},
			test: test{
				name: "Unknown error",
				mockMethod: func(method string) ([]todoitem.TodoItem, error) {
					return []todoitem.TodoItem{}, terr.UnknownError()
				},
				expect:    []todoitem.TodoItem{},
				expectErr: &expectErr{500, "unknown"},
			},
		},
		{
			reqBody: todoitem.TodoItem{Summary: newSummary("TESTING ALL THE THINGS")},
			test: test{
				name: "internal",
				mockMethod: func(method string) ([]todoitem.TodoItem, error) {
					return []todoitem.TodoItem{}, errors.New("test unknown error")
				},
				expect:    []todoitem.TodoItem{},
				expectErr: &expectErr{500, "internal"},
			},
		},
	}
	for _, tt := range tests {
		tf := func(t *testing.T) {
			parent := chi.NewRouter()
			mocks := &MockStorer{}
			subject := NewTodoHandlers(NewCore(mocks))
			subject.RegisterTodoEndpoints(parent, "/api")
			rr := httptest.NewRecorder()
			body, _ := json.Marshal(tt.reqBody)
			req := httptest.NewRequest(http.MethodPost, "/api/todo", bytes.NewReader(body))

			mocks.resp = tt.mockMethod
			parent.ServeHTTP(rr, req)

			if tt.expectErr != nil {
				assert.Equal(t, tt.expectErr.statusCode, rr.Result().StatusCode, "Should have correct status code")
				defer rr.Result().Body.Close()
				b, _ := io.ReadAll(rr.Result().Body)
				assert.Contains(t, string(b), tt.expectErr.bodyContains, "error code should contain reference to details")
			} else {
				var resitem todoitem.TodoItem
				if err := json.NewDecoder(rr.Body).Decode(&resitem); err != nil {
					assert.Fail(t, "failed to decode body", err)
				}
				assert.EqualValues(t, tt.expect[0].Summary, resitem.Summary, "response body should marshal to the expected item array")
			}
		}
		t.Run(tt.name, tf)
	}
}

func TestEmptyBodies(t *testing.T) {
	type emptyBodyTest struct {
		path    string
		reqType string
	}
	tests := []emptyBodyTest{
		{
			path:    "/api/todo",
			reqType: http.MethodPost,
		},
		{
			path:    "/api/todo/4545454",
			reqType: http.MethodPatch,
		},
	}

	for _, tt := range tests {
		tf := func(t *testing.T) {
			parent := chi.NewRouter()
			mocks := &MockStorer{}
			subject := NewTodoHandlers(NewCore(mocks))
			subject.RegisterTodoEndpoints(parent, "/api")
			rr := httptest.NewRecorder()
			req := httptest.NewRequest(tt.reqType, tt.path, nil)

			parent.ServeHTTP(rr, req)
			assert.Equal(t, http.StatusBadRequest, rr.Result().StatusCode, "Should have correct status code")

		}
		t.Run(tt.reqType, tf)

	}
}

func TestFigureDecodeError(t *testing.T) {
	type dtest struct {
		body      []byte
		err       error
		errPrefix string
		expect    int
		name      string
	}
	tests := []dtest{
		{
			name:   "syntaxerr",
			body:   []byte("hey hey hey hey hey"),
			err:    &json.SyntaxError{Offset: 13},
			expect: http.StatusBadRequest,
		},
		{
			name:   "unmarshal",
			body:   []byte("hey hey hey hey hey"),
			err:    &json.UnmarshalTypeError{Field: "unknown", Offset: 1},
			expect: http.StatusBadRequest,
		},
		{
			name:   "eofErr",
			body:   []byte("hey hey hey hey hey"),
			err:    io.ErrUnexpectedEOF,
			expect: http.StatusBadRequest,
		},
		{
			name:   "unknown fields",
			body:   []byte("hey hey hey hey hey"),
			err:    errors.New("json: unknown field hey hey hey"),
			expect: http.StatusBadRequest,
		},
		{
			name:   "too large",
			body:   []byte("an overly large body"),
			err:    errors.New("http: request body too large"),
			expect: http.StatusRequestEntityTooLarge,
		},
		{
			name:   "default unknown",
			body:   []byte("unknown decode error"),
			err:    errors.New("we're at a point where this is unknowalbe"),
			expect: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		tf := func(t *testing.T) {
			rr := httptest.NewRecorder()

			req := httptest.NewRequest(http.MethodPost, "/api/todo/", bytes.NewReader(tt.body))
			figureDecodeError(tt.err, rr, req)

			assert.Equal(t, tt.expect, rr.Result().StatusCode)

		}
		t.Run(tt.name, tf)
	}

}
