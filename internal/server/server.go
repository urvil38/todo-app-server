package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	chi_middleware "github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
	"github.com/urvil38/todo-app/internal/config"
	"github.com/urvil38/todo-app/internal/log"
	"github.com/urvil38/todo-app/internal/memory"
	"github.com/urvil38/todo-app/internal/middleware"
	"github.com/urvil38/todo-app/internal/postgres"
	"github.com/urvil38/todo-app/internal/task"
	"github.com/urvil38/todo-app/internal/telementry"
	"github.com/urvil38/todo-app/internal/version"
)

type Server struct {
	listenAddr  string
	server      *http.Server
	logger      *logrus.Logger
	taskManager task.Manager
}

func New(ctx context.Context,cfg config.Config) *Server {
	s := Server{
		listenAddr:  cfg.Addr + ":" + cfg.Port,
		logger:      log.Logger,
	}

	if cfg.UseDB {
		s.taskManager = postgres.NewTaskManager(ctx, cfg)
	} else {
		s.taskManager = memory.NewTaskManager()
	}

	return &s
}

func (s *Server) Run(ctx context.Context, cfg config.Config) {

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	router := telementry.NewRouter(nil)
	s.Install(router.Handle)

	views := append(ServerViews,
		task.TaskCreatedCountView,
		task.TaskUpdatedCountView,
		task.TaskDeletedCountView,
	)

	if err := telementry.Init(cfg, views...); err != nil {
		s.logger.Fatal(ctx, err)
	}

	mw := middleware.Chain(
		chi_middleware.SetHeader("content-type", "application/json"),
		middleware.RequestLog(s.logger),
		middleware.Timeout(10*time.Second),
	)

	s.server = &http.Server{
		Addr:         s.listenAddr,
		Handler:      mw(router),
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	go s.start()

	<-signalCh
	s.logger.Info("Received SIGINT Signal")

	s.shutdown()
}

func (s *Server) Install(handle func(string, string, http.Handler)) {
	handle(http.MethodGet, "/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	}))
	handle(http.MethodGet, "/version", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, fmt.Sprintf("version: %v\ncommit: %v", version.Version, version.Commit))
	}))
	handle(http.MethodPost, "/task", http.HandlerFunc(s.createTaskHandler))
	handle(http.MethodGet, "/tasks", http.HandlerFunc(s.listTasksHandler))
	handle(http.MethodGet, "/task/{id}", http.HandlerFunc(s.getTaskHandler))
	handle(http.MethodPost, "/task/{id}", http.HandlerFunc(s.updateTaskHandler))
	handle(http.MethodDelete, "/task/{id}", http.HandlerFunc(s.deleteTaskHandler))
}

func (s *Server) start() {
	s.logger.Infof("Server is running on %s", s.listenAddr)
	err := s.server.ListenAndServe()
	if err != http.ErrServerClosed && !errors.Is(err, http.ErrServerClosed) {
		s.logger.Fatal(err)
	}
}

func (s *Server) shutdown() {
	s.logger.Info("shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := s.server.Shutdown(ctx)
	if err != nil {
		if err == context.DeadlineExceeded {
			s.logger.Error("Error while closing server: ", err)
		} else {
			s.logger.Infof("Error while shutting down server: %s", err)
		}
	} else {
		s.logger.Info("server shutdown successfully")
	}
}
