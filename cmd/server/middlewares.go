package server

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/therohans/HungryLegs/internal/models"
)

// A private key for context that only this package can access. This is important
// to prevent collisions between different context uses
type contextKey struct {
	name string
}

var dbCtxKey = &contextKey{"db"}
var configCtxKey = &contextKey{"config"}

// DBFromContext finds the user from the context. REQUIRES Middleware to have run.
func DBFromContext(ctx context.Context) *sql.DB {
	raw, _ := ctx.Value(dbCtxKey).(*sql.DB)
	return raw
}

func ConfigFromContext(ctx context.Context) *models.StaticConfig {
	raw, _ := ctx.Value(configCtxKey).(*models.StaticConfig)
	return raw
}

// DBMiddleware decodes the share session cookie and packs the session into context
func DBMiddleware(db *sql.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// put it in context
			ctx := context.WithValue(r.Context(), dbCtxKey, db)
			// and call the next with our new context
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func ConfigMiddleware(config *models.StaticConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// put it in context
			ctx := context.WithValue(r.Context(), configCtxKey, config)
			// and call the next with our new context
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}
