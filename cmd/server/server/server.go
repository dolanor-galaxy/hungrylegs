package main

import (
	"database/sql"
	"expvar"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/handler"
	"github.com/pkg/errors"
	"github.com/robrohan/HungryLegs/cmd/server"
	"github.com/robrohan/HungryLegs/cmd/server/mid"

	"github.com/ardanlabs/conf"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/robrohan/HungryLegs/internal/importer"
	"github.com/robrohan/HungryLegs/internal/models"
	"github.com/robrohan/HungryLegs/internal/repository"
)

// will be replaced with git hash
var build = "develop"

func main() {
	if err := run(); err != nil {
		log.Println("error :", err)
		os.Exit(1)
	}
}

func run() error {
	// =========================================================================
	// Logging
	log := log.New(os.Stdout, "HUNGRYLEGS : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	// =========================================================================
	// Configuration
	cfg := models.Config

	if err := conf.Parse(os.Args[1:], "HL", &cfg); err != nil {
		if err == conf.ErrHelpWanted {
			usage, err := conf.Usage("HL", &cfg)
			if err != nil {
				return errors.Wrap(err, "generating config usage")
			}
			fmt.Println(usage)
			return nil
		}
		return errors.Wrap(err, "parsing config")
	}

	// =========================================================================
	// App Starting
	expvar.NewString("build").Set(build)
	log.Printf("Started : Application initializing : version %q", build)
	defer log.Println("Completed")

	out, err := conf.String(&cfg)
	if err != nil {
		return errors.Wrap(err, "generating config for output")
	}
	log.Printf("Config :\n%v\n", out)

	// =========================================================================
	// Root athlete
	log.Printf("Creating root athlete...")
	rootAthlete := models.NewAthlete(&cfg.Base.Root, &cfg.Base.Root)

	// =========================================================================
	// Start Database
	log.Println("Initializing database support")

	db, err := repository.OpenDatabase(
		cfg.DB.Driver, cfg.DB.Connection, cfg.DB.Post, rootAthlete)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		log.Printf("Database Stopping : %s", cfg.DB.Connection)
		db.Close()
	}()

	// =========================================================================
	// Start Debug Service
	//
	// /debug/pprof - Added to the default mux by importing the net/http/pprof package.
	// /debug/vars - Added to the default mux by importing the expvar package.
	//
	// Not concerned with shutting this down when the application is shutdown.
	log.Println("Initializing debugging support")
	go func() {
		log.Printf("Debug Listening %s", cfg.Web.DebugHost)
		log.Printf("Debug Listener closed : %v", http.ListenAndServe(cfg.Web.DebugHost, http.DefaultServeMux))
	}()

	// Put the API on top of the connection
	repo := repository.Attach(*rootAthlete.Alterego, db, cfg.DB.Driver)

	// =========================================================================
	// Start Initial DB import
	log.Println("Starting import...")
	start := time.Now()
	// Launch the activity importer
	importer.ImportActivites(cfg.Base.Import, repo)
	t := time.Now()
	elapsed := t.Sub(start)
	log.Printf("Full import took %v", elapsed)

	// =========================================================================
	// Start API Service
	log.Println("Initializing API support")

	api := http.Server{
		Addr:         cfg.Web.APIHost,
		Handler:      newAPI(log, db, cfg.DB.Driver),
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
	}

	log.Printf("API listening on %s", api.Addr)
	api.ListenAndServe()

	return nil
}

func newAPI(log *log.Logger, db *sql.DB, driver string) http.Handler {
	// Construct the web.App which holds all routes as well as common Middleware.
	app := server.NewApp(log, mid.Logger(log), mid.Errors(log))

	app.Handle("/", handler.Playground("GraphQL playground", "/query"))
	app.Handle("/query", handler.GraphQL(
		server.NewExecutableSchema(server.Config{Resolvers: &server.Resolver{
			DB:     db,
			Driver: driver,
		}})))

	// Register health check endpoint. This route is not authenticated.
	// check := Check{
	// 	db: db,
	// }
	// app.Handle("GET", "/v1/health", check.Health)
	return app
}
