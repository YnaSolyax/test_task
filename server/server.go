package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type Todo struct {
	Id          int    `json:"id"`
	Header      string `json:"header"`
	Description string `json:"description"`
	Status      bool   `json:"status"`
}

type TodoStore struct {
	mu     sync.Mutex
	Todos  map[int]*Todo
	IdIndx int
}

func NewTodoStore() *TodoStore {
	return &TodoStore{
		Todos:  make(map[int]*Todo, 10),
		IdIndx: 1,
	}
}

func StartServ() {
	port := ":8080"

	if err := http.ListenAndServe(port, nil); err != nil {
		fmt.Println("Ошибка запуска HTTP сервера: ")
		return
	}
}

func Handler(tds *TodoStore) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		if r.Method != "POST" && r.Method != "GET" {
			http.Error(w, "error method", http.StatusBadRequest)
			return
		}

		todo := Todo{}

		if r.Method == "POST" {

			tds.mu.Lock()
			defer tds.mu.Unlock()
			defer r.Body.Close()

			bytes, err := io.ReadAll(r.Body)
			if err != nil {
				fmt.Fprintf(w, "error read data")
				return
			}
			err = json.Unmarshal(bytes, &todo)
			if err != nil {
				fmt.Fprintf(w, "error unmarshal")
				return
			}
			if strings.TrimSpace(todo.Header) == "" {
				http.Error(w, "Header empty", http.StatusBadRequest)
				return
			}
			todo.Id = tds.IdIndx
			tds.Todos[todo.Id] = &todo
			tds.IdIndx++
			w.Write(bytes)
		}

		if r.Method == "GET" {
			tds.mu.Lock()
			defer tds.mu.Unlock()
			bytes, err := json.Marshal(tds.Todos)
			if err != nil {
				fmt.Fprintf(w, "error marshal")
				return
			}
			w.Write(bytes)
		}
	})

	mux.HandleFunc("/todos/", func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		path := r.URL.Path
		idStr := strings.TrimPrefix(path, "/todos/")
		id, err := strconv.Atoi(idStr)

		if err != nil {
			fmt.Fprintf(w, "error Atoi")
			return
		}

		if r.Method == "GET" {
			tds.mu.Lock()
			defer tds.mu.Unlock()
			td, ok := tds.Todos[id]
			if !ok {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			bytes, err := json.Marshal(td)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			w.Write(bytes)
		}

		if r.Method == "PUT" {
			tds.mu.Lock()
			defer tds.mu.Unlock()
			defer r.Body.Close()

			bytes, err := io.ReadAll(r.Body)
			if err != nil {
				fmt.Fprintf(w, "Error reading data")
				return
			}
			td := tds.Todos[id]
			if td.Header == "" {
				http.Error(w, "Header empty", http.StatusBadRequest)
				return
			}
			err = json.Unmarshal(bytes, &td)
			if err != nil {
				fmt.Fprintf(w, "Error unmarshal id")
				return
			}
			w.WriteHeader(http.StatusAccepted)
			w.Write(bytes)
		}

		if r.Method == "DELETE" {
			tds.mu.Lock()
			defer tds.mu.Unlock()
			delete(tds.Todos, id)
			bytes, err := json.Marshal(tds.Todos)
			if err != nil {
				fmt.Fprintf(w, "Error marshal id")
				return
			}
			w.WriteHeader(http.StatusAccepted)
			w.Write(bytes)
		}

	})

	return mux
}
