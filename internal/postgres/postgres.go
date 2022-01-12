package postgres

import (
	"context"
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"
	"contrib.go.opencensus.io/integrations/ocsql"
	"github.com/urvil38/todo-app/internal/config"
	"github.com/urvil38/todo-app/internal/database"
	"github.com/urvil38/todo-app/internal/log"
)


func OpenDB(ctx context.Context, cfg *config.Config) (_ *DB, err error) {

	// Wrap the postgres driver with our own wrapper, which adds OpenCensus instrumentation.
	ocDriver, err := database.RegisterOCWrapper("pgx", ocsql.WithAllTraceOptions())
	if err != nil {
		return nil, fmt.Errorf("unable to register the ocsql driver: %v", err)
	}
	log.Logger.Infof("opening database on host %s", cfg.DBHost)
	ddb, err := database.Open(ocDriver, cfg.DBConnInfo())
	if err != nil {
		return nil, err
	}
	log.Logger.Infof("database open finished")
	return New(ddb), nil
}

type DB struct {
	db *database.DB
}

// New returns a new postgres DB.
func New(db *database.DB) *DB {
	return &DB{
		db: db,
	}
}

// Close closes a DB.
func (db *DB) Close() error {
	return db.db.Close()
}

// Underlying returns the *database.DB inside db.
func (db *DB) Underlying() *database.DB {
	return db.db
}
