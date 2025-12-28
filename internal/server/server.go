package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"test_task/internal/storage"
)

type Todo struct {
	Id          int     `json:"id"`
	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`
	Status      string  `json:"status,omitempty"`
}

func (t Todo) toIntStatus() (int, error) {
	switch t.Status {
	case "", "created":
		return 0, nil
	case "finished":
		return 1, nil
	default:
		return 0, errors.New("invalid status")
	}
}

func fromIntStatus(status int) string {
	switch status {
	case 0:
		return "created"
	case 1:
		return "finished"
	default:
		return ""
	}
}

type Server struct {
	srv     *http.Server
	storage *storage.Inmemory
}

func NewServer(addr string, storage *storage.Inmemory) *Server {
	s := &Server{
		srv: &http.Server{
			Addr: addr,
		},
		storage: storage,
	}
	s.srv.Handler = s.handler()
	return s
}

func (s *Server) Start() error {
	if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	if err := s.srv.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}

func (s *Server) handler() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method != "POST" && r.Method != "GET" {
			http.Error(w, "unsupported method", http.StatusBadRequest)
			return
		}

		todo := &Todo{}

		if r.Method == "POST" {
			defer r.Body.Close()

			bytes, err := io.ReadAll(r.Body)
			if err != nil {
				fmt.Printf("Error: can't read request body: %v\n", err)
				http.Error(w, "can't read request body", http.StatusInternalServerError)
				return
			}
			err = json.Unmarshal(bytes, todo)
			if err != nil {
				fmt.Printf("Error: can't unmarshal request body: %v\n", err)
				http.Error(w, "can't unmarshal request body", http.StatusInternalServerError)
				return
			}
			if strings.TrimSpace(todo.Title) == "" {
				http.Error(w, "title is empty", http.StatusBadRequest)
				return
			}
			status := 0
			if todo.Status != "" {
				status, err = todo.toIntStatus()
				if err != nil {
					http.Error(w, "status is invalid, must be created, finished or empty", http.StatusBadRequest)
					return
				}
			}
			description := ""
			if todo.Description != nil {
				description = *todo.Description
			}
			err = s.storage.Create(&storage.Todo{
				Title: todo.Title,
				Description: storage.OptinonalString{
					Value: description,
					IsSet: todo.Description != nil,
				},
				Status: storage.OptinonalInt{
					Value: status,
					IsSet: todo.Status != "",
				},
			})
			if err != nil {
				fmt.Printf("Error: can't create task: %v\n", err)
				http.Error(w, "can't create task", http.StatusInternalServerError)
				return
			}
		}

		if r.Method == "GET" {
			items, err := s.storage.GetAll()
			if err != nil {
				fmt.Printf("Error: can't get all tasks: %v\n", err)
				http.Error(w, "can't get all tasks", http.StatusInternalServerError)
				return
			}

			resp := make([]Todo, 0, len(items))
			for _, item := range items {
				var desc *string
				if item.Description.IsSet {
					desc = &item.Description.Value
				}
				var status int
				if item.Status.IsSet {
					status = item.Status.Value
				}
				resp = append(resp, Todo{
					Id:          item.Id,
					Title:       item.Title,
					Description: desc,
					Status:      fromIntStatus(status),
				})
			}

			bytes, err := json.Marshal(resp)
			if err != nil {
				fmt.Printf("Error: can't encode response body: %v\n", err)
				http.Error(w, "can't encode response body", http.StatusInternalServerError)
				return
			}
			w.Write(bytes)
		}
	})

	mux.HandleFunc("/todos/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		if r.Method == "GET" {
			t, err := s.storage.Get(id)
			if err != nil {
				fmt.Printf("Error: can't get task: %v\n", err)
				http.Error(w, "can't get task", http.StatusInternalServerError)
			}
			if t == nil {
				http.Error(w, "no task with such id", http.StatusNotFound)
			}

			var desc *string
			if t.Description.IsSet {
				desc = &t.Description.Value
			}
			var status int
			if t.Status.IsSet {
				status = t.Status.Value
			}

			resp := Todo{
				Id:          t.Id,
				Title:       t.Title,
				Description: desc,
				Status:      fromIntStatus(status),
			}

			bytes, err := json.Marshal(resp)
			if err != nil {
				fmt.Printf("Error: can't encode resp: %v\n", err)
				http.Error(w, "can't encode resp", http.StatusInternalServerError)
				return
			}
			w.Write(bytes)
		}

		if r.Method == "PUT" {
			defer r.Body.Close()

			bytes, err := io.ReadAll(r.Body)
			if err != nil {
				fmt.Printf("Error: can't read request: %v\n", err)
				http.Error(w, "can't read request", http.StatusInternalServerError)
				return
			}
			todo := &Todo{}
			err = json.Unmarshal(bytes, todo)
			if err != nil {
				fmt.Printf("Error: can't decode request: %v\n", err)
				http.Error(w, "can't decode request", http.StatusInternalServerError)
				return
			}

			if todo.Title == "" {
				http.Error(w, "can't decode request", http.StatusBadRequest)
				return
			}

			description := ""
			if todo.Description != nil {
				description = *todo.Description
			}
			status := 0
			if todo.Status != "" {
				status, err = todo.toIntStatus()
				if err != nil {
					http.Error(w, "status is invalid, must be created, finished or empty", http.StatusBadRequest)
					return
				}
			}
			err = s.storage.Update(id, storage.Todo{
				Title: todo.Title,
				Description: storage.OptinonalString{
					Value: description,
					IsSet: todo.Description != nil,
				},
				Status: storage.OptinonalInt{
					Value: status,
					IsSet: todo.Status != "",
				},
			})
			if err != nil {
				if err == storage.ErrNotFound {
					http.Error(w, "task not found", http.StatusBadRequest)
					return
				}
				fmt.Printf("Error: can't update task: %v\n", err)
				http.Error(w, "can't update task", http.StatusInternalServerError)
				return
			}
		}

		if r.Method == "DELETE" {
			err := s.storage.Delete(id)
			if err != nil {
				if err == storage.ErrNotFound {
					http.Error(w, "task not found", http.StatusBadRequest)
					return
				}
				fmt.Printf("Error: can't delete task: %v\n", err)
				http.Error(w, "can't delete task", http.StatusInternalServerError)
				return
			}
		}
	})

	return mux
}
