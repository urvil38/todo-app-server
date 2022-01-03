package server

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	taskpkg "github.com/urvil38/todo-app/internal/task"
)

func (s *Server) createTaskHandler(w http.ResponseWriter, r *http.Request) {
	type payload struct {
		Name string `json:"task_name"`
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

	tasks := s.taskManager.ListTasks(r.Context())

	encoder := json.NewEncoder(w)

	err := encoder.Encode(tasks)
	if err != nil {
		s.logger.Error("listTasksHandler: json encoding err: ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (s *Server) getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	task, err := s.taskManager.GetTask(r.Context(), id)
	if err == taskpkg.ErrTaskNotFound {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
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
		Name string `json:"task_name"`
	}
	decoder := json.NewDecoder(r.Body)

	var p payload
	err := decoder.Decode(&p)
	if err != nil || p == (payload{}) {
		s.logger.Error("updateTaskHandler: unable to decode request body")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	id := mux.Vars(r)["id"]

	task, err := s.taskManager.UpdateTask(r.Context(), id, p.Name)
	if err == taskpkg.ErrTaskNotFound {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
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
	id := mux.Vars(r)["id"]

	err := s.taskManager.DeleteTask(r.Context(), id)
	if err == taskpkg.ErrTaskNotFound {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	s.logger.Infof("task deleted with id: %v", id)
}
