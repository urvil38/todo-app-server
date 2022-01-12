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

	cfg.Dump(log.Logger.Out)

	if cfg.DebugPort != "" {
		debugServer, err := telementry.NewServer()
		if err != nil {
			log.Logger.Fatal(ctx, err)
		}
		dAddr := cfg.Addr + ":" + cfg.DebugPort
		go http.ListenAndServe(dAddr, debugServer)
		log.Logger.Info("debug server is running on: ", dAddr)
	}

	s := server.New(ctx, *cfg)
	s.Run(ctx, *cfg)
}
