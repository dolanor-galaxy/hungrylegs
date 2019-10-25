package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/handler"
	"github.com/gorilla/mux"
	"github.com/robrohan/HungryLegs/cmd/server"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/robrohan/HungryLegs/internal/importer"
	"github.com/robrohan/HungryLegs/internal/models"
	"github.com/robrohan/HungryLegs/internal/repository"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	////////////////////////////////////////////
	config, err := models.LoadConfig("config.json")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Creating root athlete...")
	rootAthlete := models.NewAthlete(&config.RootAthlete, &config.RootAthlete)
	db, err := repository.OpenDatabase(config, rootAthlete)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Put the API on top of the connection
	repo := repository.Attach(*rootAthlete.Alterego, db, config)

	log.Printf("Starting import...")
	start := time.Now()
	// Launch the activity importer
	importer.ImportActivites(config.ImportDir, repo)
	t := time.Now()
	elapsed := t.Sub(start)
	log.Printf("Full import took %v", elapsed)
	////////////////////////////////////////////

	router := mux.NewRouter()
	router.Use(server.DBMiddleware(db))
	router.Use(server.ConfigMiddleware(config))

	router.Handle("/", handler.Playground("GraphQL playground", "/query"))
	router.Handle("/query", handler.GraphQL(
		server.NewExecutableSchema(server.Config{Resolvers: &server.Resolver{}})))

	http.Handle("/", router)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
