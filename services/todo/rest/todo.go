package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	terr "github.com/stumacwastaken/todo/errors"
	"github.com/stumacwastaken/todo/todoitem"
)

type TodoHandlers struct {
	TodoItem *todoitem.Core
}

func NewTodoHandlers(todoItem *todoitem.Core) TodoHandlers {
	return TodoHandlers{
		TodoItem: todoItem,
	}
}

func (h *TodoHandlers) RegisterTodoEndpoints(parent *chi.Mux, prefix string) {
	todoRouter := chi.NewRouter()

	todoRouter.Get("/", h.GetTodos)
	todoRouter.Get("/{id}", h.GetTodo)
	todoRouter.Post("/", h.CreateTodo)
	todoRouter.Patch("/{id}", h.UpdateTodo)
	parent.Mount(fmt.Sprintf("%s/todo", prefix), todoRouter)
}

func (h *TodoHandlers) GetTodos(w http.ResponseWriter, r *http.Request) {
	todos, err := h.TodoItem.GetAll(r.Context())
	if err != nil {
		if v, ok := err.(*terr.TodoError); ok {
			w.WriteHeader(v.HttpCode)
			w.Write([]byte(err.Error()))
		} else {
			w.WriteHeader(500)
			v = terr.InternalError()
			w.Write([]byte(v.Error()))
			return
		}
	}
	jsn, err := json.Marshal(todos)
	if err != nil {
		w.WriteHeader(500)
		err = terr.InternalError()
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(200)
	w.Write(jsn)
}

func (h *TodoHandlers) GetTodo(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	w.WriteHeader(405)
	w.Write([]byte(fmt.Sprintf(`{"hello":"%s but this method isn't implemented for this demo as it's currently unused"}`, id)))
}

func (h *TodoHandlers) CreateTodo(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	var i todoitem.TodoItem
	err := dec.Decode(&i)
	if err != nil {
		figureDecodeError(err, w, r)
		return
	}
	createdItem, err := h.TodoItem.Create(r.Context(), i)
	if err != nil {
		if v, ok := err.(*terr.TodoError); ok {
			w.WriteHeader(v.HttpCode)
			w.Write([]byte(err.Error()))
			return
		} else {
			w.WriteHeader(500)
			v = terr.InternalError()
			w.Write([]byte(v.Error()))
			return
		}
	}
	jsn, err := json.Marshal(createdItem)
	if err != nil {
		w.WriteHeader(500)
		err = terr.InternalError()
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(201)
	w.Write(jsn)
}

func (h *TodoHandlers) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	var i todoitem.TodoItem
	err := dec.Decode(&i)
	if err != nil {
		figureDecodeError(err, w, r)
		return
	}
	updatedItem, err := h.TodoItem.Update(r.Context(), i, id)
	if err != nil {
		if v, ok := err.(*terr.TodoError); ok {
			w.WriteHeader(v.HttpCode)
			w.Write([]byte(err.Error()))
			return
		} else {
			w.WriteHeader(500)
			v = terr.InternalError()
			w.Write([]byte(v.Error()))
			return
		}
	}
	jsn, err := json.Marshal(updatedItem)
	if err != nil {
		w.WriteHeader(500)
		err = terr.InternalError()
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(200)
	w.Write(jsn)
}

// https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body did a far better job of explaining this logic
// so I shamelessly use it where reasonable.
func figureDecodeError(err error, w http.ResponseWriter, r *http.Request) {
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		// Catch any syntax errors in the JSON and send an error message
		// which interpolates the location of the problem to make it
		// easier for the client to fix.
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			http.Error(w, msg, http.StatusBadRequest)

		// In some circumstances Decode() may also return an
		// io.ErrUnexpectedEOF error for syntax errors in the JSON. There
		// is an open issue regarding this at
		// https://github.com/golang/go/issues/25956.
		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := fmt.Sprintf("Request body contains badly-formed JSON")
			http.Error(w, msg, http.StatusBadRequest)

		// Catch any type errors, like trying to assign a string in the
		// JSON request body to a int field in our Person struct. We can
		// interpolate the relevant field name and position into the error
		// message to make it easier for the client to fix.
		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			http.Error(w, msg, http.StatusBadRequest)

		// Catch the error caused by extra unexpected fields in the request
		// body. We extract the field name from the error message and
		// interpolate it in our custom error message. There is an open
		// issue at https://github.com/golang/go/issues/29035 regarding
		// turning this into a sentinel error.
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			http.Error(w, msg, http.StatusBadRequest)

		// An io.EOF error is returned by Decode() if the request body is
		// empty.
		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			http.Error(w, msg, http.StatusBadRequest)

		// Catch the error caused by the request body being too large. Again
		// there is an open issue regarding turning this into a sentinel
		// error at https://github.com/golang/go/issues/30715.
		case err.Error() == "http: request body too large":
			msg := "Request body must not be larger than 1MB"
			http.Error(w, msg, http.StatusRequestEntityTooLarge)

		// Otherwise default to logging the error and sending a 500 Internal
		// Server Error response.
		default:
			log.Print(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}
}
