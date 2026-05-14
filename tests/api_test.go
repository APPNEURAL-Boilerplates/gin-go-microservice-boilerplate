package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/example/gin-microservice-boilerplate/internal/app"
	"github.com/example/gin-microservice-boilerplate/internal/config"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func newTestRouter() *gin.Engine {
	cfg := config.Config{
		ServiceName:            "gin-microservice-test",
		Environment:            "test",
		Host:                   "127.0.0.1",
		Port:                   "0",
		GinMode:                gin.TestMode,
		LogLevelName:           "error",
		ReadTimeoutSeconds:     5,
		WriteTimeoutSeconds:    5,
		ShutdownTimeoutSeconds: 5,
	}

	return app.NewRouter(cfg, nil)
}

func performRequest(router http.Handler, method string, path string, body []byte) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(method, path, bytes.NewReader(body))
	if body != nil {
		request.Header.Set("Content-Type", "application/json")
	}
	router.ServeHTTP(recorder, request)
	return recorder
}

func decodeBody(t *testing.T, recorder *httptest.ResponseRecorder) map[string]any {
	t.Helper()

	var body map[string]any
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	return body
}

func TestRootReturnsServiceMetadata(t *testing.T) {
	router := newTestRouter()

	response := performRequest(router, http.MethodGet, "/", nil)
	body := decodeBody(t, response)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}
	if body["ok"] != true {
		t.Fatalf("expected ok=true, got %v", body["ok"])
	}
}

func TestHealthReturnsHealthy(t *testing.T) {
	router := newTestRouter()

	response := performRequest(router, http.MethodGet, "/api/v1/health", nil)
	body := decodeBody(t, response)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}

	data := body["data"].(map[string]any)
	if data["status"] != "healthy" {
		t.Fatalf("expected healthy status, got %v", data["status"])
	}
}

func TestReadyReturnsReady(t *testing.T) {
	router := newTestRouter()

	response := performRequest(router, http.MethodGet, "/api/v1/ready", nil)
	body := decodeBody(t, response)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}

	data := body["data"].(map[string]any)
	if data["status"] != "ready" {
		t.Fatalf("expected ready status, got %v", data["status"])
	}
}

func TestCreateItemReturnsCreatedItem(t *testing.T) {
	router := newTestRouter()
	payload := []byte(`{"name":"Keyboard","description":"Mechanical keyboard","price":99.99}`)

	response := performRequest(router, http.MethodPost, "/api/v1/items", payload)
	body := decodeBody(t, response)

	if response.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d and body %s", http.StatusCreated, response.Code, response.Body.String())
	}

	data := body["data"].(map[string]any)
	item := data["item"].(map[string]any)
	if item["name"] != "Keyboard" {
		t.Fatalf("expected item name Keyboard, got %v", item["name"])
	}
}

func TestGetItemReturnsNotFound(t *testing.T) {
	router := newTestRouter()

	response := performRequest(router, http.MethodGet, "/api/v1/items/missing", nil)
	body := decodeBody(t, response)

	if response.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, response.Code)
	}

	errorPayload := body["error"].(map[string]any)
	if errorPayload["code"] != "NOT_FOUND" {
		t.Fatalf("expected NOT_FOUND, got %v", errorPayload["code"])
	}
}

func TestInvalidJSONReturnsBadRequest(t *testing.T) {
	router := newTestRouter()

	response := performRequest(router, http.MethodPost, "/api/v1/items", []byte(`{"name":`))
	body := decodeBody(t, response)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, response.Code)
	}

	errorPayload := body["error"].(map[string]any)
	if errorPayload["code"] != "VALIDATION_ERROR" {
		t.Fatalf("expected VALIDATION_ERROR, got %v", errorPayload["code"])
	}
}

func TestUnknownRouteReturnsNotFound(t *testing.T) {
	router := newTestRouter()

	response := performRequest(router, http.MethodGet, "/unknown", nil)
	body := decodeBody(t, response)

	if response.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, response.Code)
	}

	errorPayload := body["error"].(map[string]any)
	if errorPayload["code"] != "NOT_FOUND" {
		t.Fatalf("expected NOT_FOUND, got %v", errorPayload["code"])
	}
}

func TestUnsupportedMethodReturnsMethodNotAllowed(t *testing.T) {
	router := newTestRouter()

	response := performRequest(router, http.MethodDelete, "/api/v1/items", nil)
	body := decodeBody(t, response)

	if response.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status %d, got %d", http.StatusMethodNotAllowed, response.Code)
	}

	errorPayload := body["error"].(map[string]any)
	if errorPayload["code"] != "METHOD_NOT_ALLOWED" {
		t.Fatalf("expected METHOD_NOT_ALLOWED, got %v", errorPayload["code"])
	}
}

func TestRequestIDIsReturned(t *testing.T) {
	router := newTestRouter()

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	request.Header.Set("X-Request-Id", "test-request-id")

	router.ServeHTTP(recorder, request)

	if recorder.Header().Get("X-Request-Id") != "test-request-id" {
		t.Fatalf("expected propagated request id")
	}
}
