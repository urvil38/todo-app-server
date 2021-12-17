package log

import (
	"os"

	"github.com/sirupsen/logrus"
)

// initialize logger
var logger = logrus.New()

func init() {
	// setting logger output to stdout
	logger.SetOutput(os.Stdout)
}

// Config defines format and level of logging for logger
type Config struct {
	// Format can be text | json | json-pretty. Default format is text.
	Format string
	// Level can be info | debug | error. Default level is info.
	Level string
}

// Set setting up logger using given configuration
func Set(c Config) {

	switch c.Format {
	case "text":
		logger.SetFormatter(&logrus.TextFormatter{TimestampFormat: "2006-01-02 15:04:05", DisableSorting: true, FullTimestamp: true, DisableLevelTruncation: true, DisableColors: true})
	case "json-pretty":
		logger.SetFormatter(&logrus.JSONFormatter{PrettyPrint: true, TimestampFormat: "2006-01-02 15:04:05"})
	case "json":
		fallthrough
	default:
		logger.SetFormatter(&logrus.JSONFormatter{TimestampFormat: "2006-01-02 15:04:05"})
	}

	level := logrus.InfoLevel
	if c.Level != "" {
		var err error
		level, err = logrus.ParseLevel(c.Level)
		if err != nil {
			logrus.Fatalf("Unable to parse log level: %v", err)
		}
	}
	logger.SetLevel(level)
}

// Get returns an instance of logger
func Get() *logrus.Logger {
	return logger
}
