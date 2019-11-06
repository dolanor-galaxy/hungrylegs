package mid

import (
	"context"
	"log"
	"net/http"

	"github.com/robrohan/HungryLegs/cmd/server"
	"go.opencensus.io/trace"
)

// Errors handles errors coming out of the call chain.
func Errors(log *log.Logger) server.Middleware {

	f := func(before server.Handler) server.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {
			ctx, span := trace.StartSpan(ctx, "internal.mid.Errors")
			defer span.End()

			if err := before(ctx, w, r, params); err != nil {
				// Log the error.
				// log.Printf("%s : ERROR : %+v", v.TraceID, err)
				log.Printf("ERROR : %+v", err)

				// Respond to the error.
				// if err := web.RespondError(ctx, w, err); err != nil {
				// 	return err
				// }
			}
			return nil
		}
		return h
	}
	return f
}
