package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
	"greenlight.mkabdelrahman.net/internal/data"
	"greenlight.mkabdelrahman.net/internal/dbutil"
)

const version = "1.0.0"

func main() {
	cfg := readConfigFromFlags()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	db, err := dbutil.NewDBConnection(cfg.db)
	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()

	logger.Printf("database connection pool established")

	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/healthcheck", app.healthcheckHandler)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Printf("starting %s server on %s", cfg.env, srv.Addr)
	err = srv.ListenAndServe()
	log.Fatal(err)
}
