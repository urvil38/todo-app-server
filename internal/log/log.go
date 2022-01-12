package log

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

// initialize logger
var Logger = logrus.New()

// Config defines format and level of logging for logger
type Config struct {
	// Format can be text | json | json-pretty. Default format is text.
	Format string
	// Level can be info | debug | error. Default level is info.
	Level string
}

// Set setting up logger using given configuration
func Set(c Config) error {

	Logger.SetOutput(os.Stdout)

	switch c.Format {
	case "text":
		Logger.SetFormatter(&logrus.TextFormatter{TimestampFormat: "2006-01-02 15:04:05", DisableSorting: true, FullTimestamp: true, DisableLevelTruncation: true, DisableColors: true})
	case "json-pretty":
		Logger.SetFormatter(&logrus.JSONFormatter{PrettyPrint: true, TimestampFormat: "2006-01-02 15:04:05"})
	case "json":
		fallthrough
	default:
		Logger.SetFormatter(&logrus.JSONFormatter{TimestampFormat: "2006-01-02 15:04:05"})
	}

	level := logrus.InfoLevel
	if c.Level != "" {
		var err error
		level, err = logrus.ParseLevel(c.Level)
		if err != nil {
			return fmt.Errorf("Unable to parse log level: %w", err)
		}
	}
	Logger.SetLevel(level)
	return nil
}

// Get returns an instance of logger
func Get() *logrus.Logger {
	return Logger
}
