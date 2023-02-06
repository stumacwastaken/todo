package rest

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

//this acts as a mini integration test. I feel ok using Small sleeps here
//for that purpose
func TestNewServerStartupAndShutdown(t *testing.T) {
	srv := NewServer("localhost", "10901")
	srv.Router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(`{"hello":"test"}`))
	})
	assert.NotNil(t, srv, "Server should not be nil")
	assert.Equal(t, "localhost:10901", srv.Addr)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	go srv.Start(context.Background())
	time.Sleep(10 * time.Millisecond)
	res, err := http.Get("http://localhost:10901/")
	assert.Nil(t, err)
	assert.Equal(t, 200, res.StatusCode, "server should be running")
	assert.Equal(t, []string{"application/json"}, res.Header["Content-Type"])
	srv.Router.ServeHTTP(rr, req)
	assert.Equal(t, 200, rr.Result().StatusCode, "should have same dummy code")
	go srv.Stop()
	res, err = http.Get("localhost:10901/")
	assert.Error(t, err)
}
