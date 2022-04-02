package main

import (
	"backend/models"
	"context"
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"time"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
}

type AppStatus struct {
	Status      string `json:"status"`
	Environment string `json:"environment"`
	Version     string `json:"version"`
}

//application struct: have all the information we need to share with our handler and other
//components of our application
type application struct {
	config config
	logger *log.Logger
	models models.Models
}

func main() {
	var cfg config
	flag.IntVar(&cfg.port, "port", 4000, "Server port to listen on")
	flag.StringVar(&cfg.env, "env", "development", "Application environment (development|production")
	flag.StringVar(&cfg.db.dsn, "dsn",
		"postgres://postgres:password@localhost/go-movies?sslmode=disable",
		"Postgres connection")
	// parse the flags
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	db, err := openDB(cfg)
	if err != nil {
		logger.Fatal(err)
	}
	app := &application{
		config: cfg,
		logger: logger,
		models: models.NewModels(db),
	}
	defer db.Close()

	// accomplish basic web server running and serve some json to the browser

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  time.Minute,
	}
	logger.Println("Starting server on port", cfg.port)
	err = srv.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
