package rest

import (
	"fmt"
	"net/http"

	"github.com/stumacwastaken/todo/tracing"
)

// TraceRequest a middlware to our http calls to provide a super basic trace should we ever forget to add in a trace
// span at the handler level.
func TraceRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, span := tracing.Tracer().Start(r.Context(), fmt.Sprintf("%s:%s", r.Method, r.URL.Path))
		defer span.End()

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
