package main

import (
	"context"
	stdlog "log"
	"net/http"

	"github.com/urvil38/todo-app/internal/config"
	"github.com/urvil38/todo-app/internal/log"
	"github.com/urvil38/todo-app/internal/server"
	"github.com/urvil38/todo-app/internal/telementry"
)

func main() {

	ctx := context.Background()
	cfg, err := config.Init(ctx)
	if err != nil {
		stdlog.Fatal(err)
	}

	log.Set(log.Config{
		Format: cfg.LogFormat,
		Level:  cfg.LogLevel,
	})

	logger := log.Get()
	cfg.Dump(logger.Writer())

	if cfg.DebugPort != "" {
		debugServer, err := telementry.NewServer()
		if err != nil {
			logger.Fatal(ctx, err)
		}
		dAddr := cfg.Addr + ":" + cfg.DebugPort
		go http.ListenAndServe(dAddr, debugServer)
		logger.Info("debug server is running on: ", dAddr)
	}

	s := server.New(*cfg)
	s.Run(ctx, *cfg)
}
