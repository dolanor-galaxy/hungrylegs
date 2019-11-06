package server

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// App is a wrapper for our router and middleware
type App struct {
	*mux.Router
	log *log.Logger
	mw  []Middleware
}

// A Handler is a type that handles an http request within our own little mini framework.
type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error

// Middleware is a function designed to run some code before and/or after another Handler.
type Middleware func(Handler) Handler

// wrapMiddleware creates a new handler by wrapping middleware around a final handler.
func wrapMiddleware(mw []Middleware, handler Handler) Handler {
	for i := len(mw) - 1; i >= 0; i-- {
		h := mw[i]
		if h != nil {
			handler = h(handler)
		}
	}

	return handler
}

// NewApp creates an App value that handle a set of routes for the application.
func NewApp(log *log.Logger, mw ...Middleware) *App {
	app := App{
		Router: mux.NewRouter(),
		log:    log,
		mw:     mw,
	}
	return &app
}
