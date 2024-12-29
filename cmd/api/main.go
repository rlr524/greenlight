package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/rlr524/greenlight/internal/dal"
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
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  time.Duration
	}
}

type application struct {
	config           config
	logger           *slog.Logger
	dataAccessLayers dal.DataAccessLayers
}

func main() {
	var cfg config

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	err := godotenv.Load()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	dbDSN := os.Getenv("DB_DSN")

	flag.IntVar(&cfg.port, "port", 4000, "API Server port")
	flag.StringVar(&cfg.env, "env", "development",
		"Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", dbDSN, "PostgreSQL DSN")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25,
		"PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25,
		"PostgreSQL max idle connections")
	flag.DurationVar(&cfg.db.maxIdleTime, "db-max-idle-time", 15*time.Minute,
		"PostgreSQL max connection idle time")

	flag.Parse()

	// The DB connection pool is established and if there is an error, it is
	// logged, and we exit the application with a code 1 immediately.
	db, err := openDB(cfg)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// Defer a call to db.Close() so that the connection pool is closed before main() exits and
	// log a message that the connection pool has been successfully established. We use a
	// lambda to do this so we can handle any errors thrown upon closing the pool. Because
	// this function is only called upon exit of the main() function, we go ahead and just exit
	// the app in case of an error with a code 1 which will close any open connections.
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}
	}(db)
	logger.Info("database connection pool established")

	app := &application{
		config:           cfg,
		logger:           logger,
		dataAccessLayers: dal.NewDALs(db),
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
