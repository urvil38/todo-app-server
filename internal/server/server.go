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

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/urvil38/todo-app/internal/config"
	"github.com/urvil38/todo-app/internal/log"
	"github.com/urvil38/todo-app/internal/middleware"
	"github.com/urvil38/todo-app/internal/task"
	"github.com/urvil38/todo-app/internal/telementry"
)

type Server struct {
	listenAddr  string
	server      *http.Server
	logger      *logrus.Logger
	taskManager task.Manager
}

func New(cfg config.Config) *Server {
	return &Server{
		listenAddr:  cfg.Addr + ":" + cfg.Port,
		logger:      log.Get(),
		taskManager: task.NewInMemoryTaskManager(),
	}
}

func (s *Server) Run(ctx context.Context) {

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	router := telementry.NewRouter(nil)
	s.Install(router.Handle)

	views := append(ServerViews,
		task.TaskCreatedCountView,
		task.TaskUpdatedCountView,
		task.TaskDeletedCountView,
	)

	if err := telementry.Init(views...); err != nil {
		s.logger.Fatal(ctx, err)
	}

	mw := middleware.Chain(
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

type MuxHandler func(route string, handler http.Handler) *mux.Route

func (s *Server) Install(handle MuxHandler) {
	handle("/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	}))
	handle("/task", middleware.Json()(http.HandlerFunc(s.createTaskHandler))).Methods("POST")
	handle("/tasks", middleware.Json()(http.HandlerFunc(s.listTasksHandler))).Methods("GET")
	handle("/task/{id}", middleware.Json()(http.HandlerFunc(s.getTaskHandler))).Methods("GET")
	handle("/task/{id}", middleware.Json()(http.HandlerFunc(s.updateTaskHandler))).Methods("POST")
	handle("/task/{id}", middleware.Json()(http.HandlerFunc(s.deleteTaskHandler))).Methods("DELETE")
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
