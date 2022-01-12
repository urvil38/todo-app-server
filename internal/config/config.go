package config

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/urvil38/todo-app/internal/log"
)

// Config hold the configuration for todo server.
type Config struct {
	// Address on which server is running
	Addr string

	// TCP port on which server is listening
	Port string

	// TCP port on which debug server is listening
	// Debug port should be different from the server port
	DebugPort string

	// LogLevel can be [info, debug, error, fatal]
	// If invalid log level is provided then info log level will be used as default
	LogLevel string

	// LogFormat can be [json, json-pretty, text]
	// Default will be text
	LogFormat string

	DBUser, DBHost, DBPort, DBName string
	DBPassword                     string `json:"-"`

	// Whether to use postgres to store tasks or not.
	// Default value is false, in that case tasks will be store in memory.
	UseDB bool
}

// StatementTimeout is the value of the Postgres statement_timeout parameter.
// Statements that run longer than this are terminated.
// 10 minutes is the App Engine standard request timeout,
// but we set this longer for the worker.
const StatementTimeout = 30 * time.Minute

// GetEnv return the environment variable by its key, returning its value
// if it exists, otherwise returning fallback value
func GetEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}

// DBConnInfo returns a PostgreSQL connection string constructed from
// environment variables, using the primary database host.
func (c *Config) DBConnInfo() string {
	return c.dbConnInfo(c.DBHost)
}

// dbConnInfo returns a PostgresSQL connection string for the given host.
func (c *Config) dbConnInfo(host string) string {
	// For the connection string syntax, see
	// https://www.postgresql.org/docs/current/libpq-connect.html#LIBPQ-CONNSTRING.
	// Set the statement_timeout config parameter for this session.
	// See https://www.postgresql.org/docs/current/runtime-config-client.html.
	timeoutOption := fmt.Sprintf("-c statement_timeout=%d", StatementTimeout/time.Millisecond)
	return fmt.Sprintf("user='%s' password='%s' host='%s' port=%s dbname='%s' sslmode=disable options='%s'",
		c.DBUser, c.DBPassword, host, c.DBPort, c.DBName, timeoutOption)
}

// Init initailize config. Config values will be read from the
// environment variables.
// Note: Call Init at the beginning of main function
func Init(ctx context.Context) (cfg *Config, err error) {
	cfg = &Config{
		Addr:       GetEnv("TODO_ADDRESS", "localhost"),
		Port:       GetEnv("TODO_PORT", "8080"),
		DebugPort:  GetEnv("TODO_DEBUG_PORT", "8081"),
		LogLevel:   GetEnv("TODO_LOG_LEVEL", "info"),
		LogFormat:  GetEnv("TODO_LOG_FORMAT", "text"),
		DBHost:     GetEnv("TODO_DATABASE_HOST", "localhost"),
		DBUser:     GetEnv("TODO_DATABASE_USER", "postgres"),
		DBPassword: os.Getenv("TODO_DATABASE_PASSWORD"),
		DBPort:     GetEnv("TODO_DATABASE_PORT", "5432"),
		DBName:     GetEnv("TODO_DATABASE_NAME", "todo-db"),
		UseDB:      os.Getenv("TODO_USE_DB") == "true",
	}

	if cfg.Port == cfg.DebugPort {
		return nil, fmt.Errorf("server port and debug port should be different. Both listening on port \"%v\"!", cfg.Port)
	}

	err = log.Set(log.Config{
		Format: cfg.LogFormat,
		Level:  cfg.LogLevel,
	})
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// Dump outputs the current config information to the given Writer.
func (cfg *Config) Dump(w io.Writer) error {
	fmt.Fprint(w, "config: ")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "    ")
	return enc.Encode(cfg)
}
