package http_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	httpadapter "github.com/gallardo/temp-converter/internal/adapter/http"
	"github.com/gallardo/temp-converter/internal/adapter/memory"
	"github.com/gallardo/temp-converter/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestHandler() *httpadapter.Handler {
	repo := memory.NewHistoryRepository()
	return httpadapter.NewHandler(
		usecase.NewConvertTemperature(repo),
		usecase.NewGetHistory(repo),
		usecase.NewConvertOnly(),
	)
}

// GET /convert tests (existing)

func TestHandleGetConvert_OK(t *testing.T) {
	h := newTestHandler()
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest("GET", "/convert?from=celsius&to=fahrenheit&value=0", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var resp map[string]float64
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
	assert.InDelta(t, 32.0, resp["result"], 0.001)
}

func TestHandleGetConvert_MissingParam(t *testing.T) {
	h := newTestHandler()
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest("GET", "/convert?from=celsius&to=fahrenheit", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestHandleGetConvert_InvalidUnit(t *testing.T) {
	h := newTestHandler()
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest("GET", "/convert?from=rankine&to=celsius&value=100", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestHandleGetConvert_InvalidValue(t *testing.T) {
	h := newTestHandler()
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest("GET", "/convert?from=celsius&to=fahrenheit&value=abc", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// POST /convert tests

func TestHandlePostConvert_OK(t *testing.T) {
	h := newTestHandler()
	mux := http.NewServeMux()
	h.Register(mux)

	body := `{"value":0,"from":"celsius","to":"fahrenheit"}`
	req := httptest.NewRequest("POST", "/convert", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var resp map[string]float64
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
	assert.InDelta(t, 32.0, resp["result"], 0.001)
}

func TestHandlePostConvert_InvalidUnit(t *testing.T) {
	h := newTestHandler()
	mux := http.NewServeMux()
	h.Register(mux)

	body := `{"value":0,"from":"celsius","to":"bad"}`
	req := httptest.NewRequest("POST", "/convert", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestHandlePostConvert_MissingUnit(t *testing.T) {
	h := newTestHandler()
	mux := http.NewServeMux()
	h.Register(mux)

	body := `{"value":0,"from":"","to":"fahrenheit"}`
	req := httptest.NewRequest("POST", "/convert", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "missing")
}

func TestHandlePostConvert_MalformedBody(t *testing.T) {
	h := newTestHandler()
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest("POST", "/convert", strings.NewReader("not-json"))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// GET /history tests

func TestHandleGetHistory_Empty(t *testing.T) {
	h := newTestHandler()
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest("GET", "/history", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "[]")
}

func TestHandleGetHistory_AfterPost(t *testing.T) {
	h := newTestHandler()
	mux := http.NewServeMux()
	h.Register(mux)

	body := `{"value":0,"from":"celsius","to":"fahrenheit"}`
	postReq := httptest.NewRequest("POST", "/convert", strings.NewReader(body))
	postReq.Header.Set("Content-Type", "application/json")
	mux.ServeHTTP(httptest.NewRecorder(), postReq)

	getReq := httptest.NewRequest("GET", "/history", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, getReq)

	assert.Equal(t, http.StatusOK, rec.Code)
	var records []map[string]interface{}
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&records))
	assert.Len(t, records, 1)
}
