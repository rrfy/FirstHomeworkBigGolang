package server

import (
	"bytes"
	"encoding/json"
	"homework1/internal/pkg/storage"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServerHealth(t *testing.T) {
	st, err := storage.NewStorage()
	assert.NoError(t, err)

	s := New("localhost:8080", &st)
	engine := s.newAPI()

	t.Run("Health Check", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/health", nil)
		resp := httptest.NewRecorder()

		engine.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
	})
}
func TestServerSetGetSring(t *testing.T) {
	st, err := storage.NewStorage()
	assert.NoError(t, err)

	s := New("localhost:8080", &st)
	engine := s.newAPI()
	t.Run("Set and Get String Value", func(t *testing.T) {
		// Устанавливаем значение для ключа
		setBody := Entry{Value: "Hello, World!"}
		bodyBytes, _ := json.Marshal(setBody)
		req, _ := http.NewRequest(http.MethodPut, "/scalar/set/greeting", bytes.NewBuffer(bodyBytes))
		resp := httptest.NewRecorder()

		engine.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusOK, resp.Code)

		// Получаем значение по ключу
		req, _ = http.NewRequest(http.MethodGet, "/scalar/get/greeting", nil)
		resp = httptest.NewRecorder()

		engine.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusOK, resp.Code)

		var getResp Entry
		json.NewDecoder(resp.Body).Decode(&getResp)

		assert.Equal(t, "Hello, World!", getResp.Value)
	})
}

func TestServerSetGetInt(t *testing.T) {
	st, err := storage.NewStorage()
	assert.NoError(t, err)

	s := New("localhost:8080", &st)
	engine := s.newAPI()
	t.Run("Set and Get Int Value", func(t *testing.T) {
		// Устанавливаем значение для ключа
		setBody := Entry{Value: 42}
		bodyBytes, _ := json.Marshal(setBody)
		req, _ := http.NewRequest(http.MethodPut, "/scalar/set/number", bytes.NewBuffer(bodyBytes))
		resp := httptest.NewRecorder()

		engine.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusOK, resp.Code)

		// Получаем значение по ключу
		req, _ = http.NewRequest(http.MethodGet, "/scalar/get/number", nil)
		resp = httptest.NewRecorder()

		engine.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusOK, resp.Code)

		var getResp Entry
		json.NewDecoder(resp.Body).Decode(&getResp)

		assert.Equal(t, float64(42), getResp.Value) // JSON всегда возвращает числа как float64
	})
}

func TestServerSetGetBool(t *testing.T) {
	st, err := storage.NewStorage()
	assert.NoError(t, err)

	s := New("localhost:8080", &st)
	engine := s.newAPI()
	t.Run("Set and Get Bool Value", func(t *testing.T) {
		// Устанавливаем значение для ключа
		setBody := Entry{Value: true}
		bodyBytes, _ := json.Marshal(setBody)
		req, _ := http.NewRequest(http.MethodPut, "/scalar/set/flag", bytes.NewBuffer(bodyBytes))
		resp := httptest.NewRecorder()

		engine.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusOK, resp.Code)

		// Получаем значение по ключу
		req, _ = http.NewRequest(http.MethodGet, "/scalar/get/flag", nil)
		resp = httptest.NewRecorder()

		engine.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusOK, resp.Code)

		var getResp Entry
		json.NewDecoder(resp.Body).Decode(&getResp)

		assert.Equal(t, true, getResp.Value)
	})
}

func TestServerSetGetFloat(t *testing.T) {
	st, err := storage.NewStorage()
	assert.NoError(t, err)

	s := New("localhost:8080", &st)
	engine := s.newAPI()
	t.Run("Set and Get Float Value", func(t *testing.T) {
		// Устанавливаем значение для ключа
		setBody := Entry{Value: 3.14159}
		bodyBytes, _ := json.Marshal(setBody)
		req, _ := http.NewRequest(http.MethodPut, "/scalar/set/pi", bytes.NewBuffer(bodyBytes))
		resp := httptest.NewRecorder()

		engine.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusOK, resp.Code)

		// Получаем значение по ключу
		req, _ = http.NewRequest(http.MethodGet, "/scalar/get/pi", nil)
		resp = httptest.NewRecorder()

		engine.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusOK, resp.Code)

		var getResp Entry
		json.NewDecoder(resp.Body).Decode(&getResp)

		assert.Equal(t, 3.14159, getResp.Value)
	})
}

func TestServerGetNone(t *testing.T) {
	st, err := storage.NewStorage()
	assert.NoError(t, err)

	s := New("localhost:8080", &st)
	engine := s.newAPI()
	t.Run("Get Non-Existing Key", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/scalar/get/nonexistent", nil)
		resp := httptest.NewRecorder()

		engine.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusNotFound, resp.Code)
	})
}
