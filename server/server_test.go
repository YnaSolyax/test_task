package server

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTodoAPI(t *testing.T) {
	tds := NewTodoStore()
	mux := Handler(tds)

	//POST
	t.Run("Create Success", func(t *testing.T) {
		body := `{"header": "Тест", "description": "Описание"}`
		req := httptest.NewRequest("POST", "/todos", bytes.NewBufferString(body))
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Ожидали 200, получили %d", rr.Code)
		}
	})

	//POST header
	t.Run("Create Validation Error", func(t *testing.T) {
		body := `{"header": "", "description": "Нет заголовка"}`
		req := httptest.NewRequest("POST", "/todos", bytes.NewBufferString(body))
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Ожидали 400 при пустом заголовке, получили %d", rr.Code)
		}
	})

	//GET
	t.Run("Get All", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/todos", nil)
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Ожидали 200, получили %d", rr.Code)
		}
	})

	//GET id
	t.Run("Get By ID", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/todos/1", nil)
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Ожидали 200 для ID 1, получили %d", rr.Code)
		}
	})

	//PUT
	t.Run("Update Todo", func(t *testing.T) {
		body := `{"header": "Обновлено", "status": true}`
		req := httptest.NewRequest("PUT", "/todos/1", bytes.NewBufferString(body))
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusAccepted {
			t.Errorf("Ожидали 202 Accepted, получили %d", rr.Code)
		}
	})

	//DELETE
	t.Run("Delete Todo", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/todos/1", nil)
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusAccepted {
			t.Errorf("Ожидали 202 Accepted при удалении, получили %d", rr.Code)
		}
	})
}
