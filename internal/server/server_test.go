package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"test_task/internal/storage"
	"testing"
)

func setupTestServer() (*Server, *storage.Inmemory) {
	st := storage.NewInmemory()
	srv := NewServer(":8080", st)
	return srv, st
}

func TestTodoFlow(t *testing.T) {
	srv, _ := setupTestServer()
	handler := srv.handler()

	// 1. POST /todos
	t.Run("Create Todo Success", func(t *testing.T) {
		body := `{"title": "Learn Go", "description": "Finish the course", "status": "created"}`
		req := httptest.NewRequest("POST", "/todos", strings.NewReader(body))
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})

	// 2. POST /todos - valid
	t.Run("Create Todo Validation Error", func(t *testing.T) {
		body := `{"title": " ", "description": "No title here"}`
		req := httptest.NewRequest("POST", "/todos", strings.NewReader(body))
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400 for empty title, got %d", w.Code)
		}
	})

	// 3. GET /todos
	t.Run("Get All Todos", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/todos", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		var todos []Todo
		json.NewDecoder(w.Body).Decode(&todos)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
		if len(todos) == 0 {
			t.Error("expected at least one todo in list")
		}
	})

	// 4. GET /todos/{id}
	t.Run("Get Todo By ID", func(t *testing.T) {

		req := httptest.NewRequest("GET", "/todos/1", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
	})

	// 5. PUT /todos/{id}
	t.Run("Update Todo Success", func(t *testing.T) {
		body := `{"title": "Updated Title", "status": "finished"}`
		req := httptest.NewRequest("PUT", "/todos/1", strings.NewReader(body))
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200 for update, got %d", w.Code)
		}
	})

	// 6. DELETE /todos/{id}
	t.Run("Delete Todo Success", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/todos/1", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200 for delete, got %d", w.Code)
		}
	})
}
