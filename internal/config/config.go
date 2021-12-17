package config

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
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
}

// GetEnv return the environment variable by its key, returning its value
// if it exists, otherwise returning fallback value
func GetEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}

// Init initailize config. Config values will be read from the
// environment variables.
// Note: Call Init at the beginning of main function
func Init(ctx context.Context) (cfg *Config, err error) {
	cfg = &Config{
		Addr:      GetEnv("TODO_ADDRESS", "0.0.0.0"),
		Port:      GetEnv("TODO_PORT", "8080"),
		DebugPort: GetEnv("TODO_DEBUG_PORT", "8081"),
		LogLevel:  GetEnv("TODO_LOG_LEVEL", "info"),
		LogFormat: GetEnv("TODO_LOG_FORMAT", "text"),
	}

	if cfg.Port == cfg.DebugPort {
		return nil, fmt.Errorf("server port and debug port should be different. Both listening on port \"%v\"!", cfg.Port)
	}

	return cfg, nil
}

// Dump outputs the current config information to the given Writer.
func (c *Config) Dump(w io.Writer) error {
	fmt.Fprint(w, "config: ")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "    ")
	return enc.Encode(c)
}
