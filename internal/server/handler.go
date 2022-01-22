package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	taskpkg "github.com/urvil38/todo-app/internal/task"
)

func (s *Server) createTaskHandler(w http.ResponseWriter, r *http.Request) {
	type payload struct {
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(r.Body)

	var p payload
	err := decoder.Decode(&p)
	if err != nil || p == (payload{}) {
		s.logger.Error("createTaskHandler: unable to decode request body")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	task, err := s.taskManager.CreateTask(r.Context(), p.Name)
	if err != nil {
		s.logger.Error("createTaskHandler: unable to create task: ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	s.logger.Infof("task created with id: %v", task.Id)
	w.WriteHeader(http.StatusCreated)
}

func (s *Server) listTasksHandler(w http.ResponseWriter, r *http.Request) {

	tasks, err := s.taskManager.ListTasks(r.Context())
	if err != nil {
		s.logger.Error("listTasksHandler: unable to list task: ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	encoder := json.NewEncoder(w)

	err = encoder.Encode(tasks)
	if err != nil {
		s.logger.Error("listTasksHandler: json encoding err: ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (s *Server) getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	task, err := s.taskManager.GetTask(r.Context(), id)
	if err != nil {
		if errors.Is(err, taskpkg.ErrTaskNotFound) {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		} else {
			s.logger.Error("getTaskHandler: unable to get task: ", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	encoder := json.NewEncoder(w)

	err = encoder.Encode(task)
	if err != nil {
		s.logger.Error("getTaskHandler: json encoding err: ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (s *Server) updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	type payload struct {
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(r.Body)

	var p payload
	err := decoder.Decode(&p)
	if err != nil || p == (payload{}) {
		s.logger.Error("updateTaskHandler: unable to decode request body")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	id := chi.URLParam(r, "id")

	task, err := s.taskManager.UpdateTask(r.Context(), id, p.Name)
	if err != nil {
		if errors.Is(err, taskpkg.ErrTaskNotFound) {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		} else {
			s.logger.Error("updateTaskHandler: unable to update task: ", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	s.logger.Infof("task updated with id: %v", task.Id)
	encoder := json.NewEncoder(w)
	err = encoder.Encode(task)
	if err != nil {
		s.logger.Error("updateTaskHandler: json encoding err: ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (s *Server) deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := s.taskManager.DeleteTask(r.Context(), id)
	if err != nil {
		if errors.Is(err, taskpkg.ErrTaskNotFound) {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		} else {
			s.logger.Error("deleteTaskHandler: unable to delete task: ", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	s.logger.Infof("task deleted with id: %v", id)
}
