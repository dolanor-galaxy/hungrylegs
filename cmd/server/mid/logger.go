package mid

import (
	"context"
	"log"
	"net/http"

	"github.com/robrohan/HungryLegs/cmd/server"
	"go.opencensus.io/trace"
)

// Logger writes some information about the request to the logs in the
func Logger(log *log.Logger) server.Middleware {
	f := func(before server.Handler) server.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {
			ctx, span := trace.StartSpan(ctx, "internal.mid.Logger")
			defer span.End()
			err := before(ctx, w, r, params)
			return err
		}
		return h
	}
	return f
}
