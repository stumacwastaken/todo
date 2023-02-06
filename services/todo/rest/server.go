package rest

import (
	"context"
	"fmt"
	"net/http"

	"github.com/stumacwastaken/todo/log"
	"github.com/stumacwastaken/todo/tracing"
	"go.uber.org/zap"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type HttpServer struct {
	Router *chi.Mux
	http.Server
}

func NewServer(addr, port string) *HttpServer {
	router := chi.NewRouter()
	router.Use(middleware.AllowContentType("application/json"))
	router.Use(middleware.CleanPath)
	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	router.Use(TraceRequest)
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			next.ServeHTTP(w, r)
		})
	})

	hs := &HttpServer{
		Router: router,
		Server: http.Server{
			Addr:    fmt.Sprintf("%s:%s", addr, port),
			Handler: router,
		},
	}
	return hs
}

func (s *HttpServer) Stop() {
	log.Default().Info("shutting down todo restful api server")
	s.Shutdown(context.Background())
}

func (s *HttpServer) Start(ctx context.Context) error {
	ctx, span := tracing.Tracer().Start(ctx, "Start")
	defer span.End()
	log.Default().Info("starting todo restful api server", zap.String("address", s.Addr))
	return s.ListenAndServe()
}
