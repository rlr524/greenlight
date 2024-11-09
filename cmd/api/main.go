package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"log/slog"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
}

type application struct {
	config config
	logger *slog.Logger
}

func main() {
	var cfg config

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	err := godotenv.Load()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	dbPass := os.Getenv("DB_PASS")
	dbUser := os.Getenv("DB_USER")

	flag.IntVar(&cfg.port, "port", 4000, "API Server port")
	flag.StringVar(&cfg.env, "env", "development",
		"Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn,
		"db-dsn",
		fmt.Sprintf("postgres://%s:%s@localhost/greenlight?sslmode=disable", dbUser, dbPass),
		"PostgreSQL DSN")

	flag.Parse()

	db, err := openDB(cfg)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)
	logger.Info("database connection pool established")

	app := &application{
		config: cfg,
		logger: logger,
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	logger.Info("starting server", "addr", srv.Addr, "env", cfg.env)

	err = srv.ListenAndServe()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
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
		_ = db.Close()
		return nil, err
	}

	return db, nil
}
