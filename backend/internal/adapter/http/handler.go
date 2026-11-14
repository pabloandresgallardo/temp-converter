package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gallardo/temp-converter/internal/domain/temperature"
	"github.com/gallardo/temp-converter/internal/usecase"
)

// Handler agrupa los handlers HTTP del servidor.
type Handler struct {
	convertUC     *usecase.ConvertTemperature
	historyUC     *usecase.GetHistory
	convertOnlyUC *usecase.ConvertOnly
}

func NewHandler(convert *usecase.ConvertTemperature, history *usecase.GetHistory, convertOnly *usecase.ConvertOnly) *Handler {
	return &Handler{convertUC: convert, historyUC: history, convertOnlyUC: convertOnly}
}

// Register registra las rutas en el mux.
func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /convert", h.handleConvert)
	mux.HandleFunc("GET /convert", h.handleGetConvert)
	mux.HandleFunc("GET /history", h.handleHistory)
}

func (h *Handler) handleConvert(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Value float64 `json:"value"`
		From  string  `json:"from"`
		To    string  `json:"to"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if body.From == "" || body.To == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing fields: from and to are required"})
		return
	}
	resp, err := h.convertUC.Execute(usecase.ConvertRequest{
		Value: body.Value,
		From:  temperature.Unit(body.From),
		To:    temperature.Unit(body.To),
	})
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]float64{"result": resp.Result})
}

func (h *Handler) handleGetConvert(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	from := q.Get("from")
	to := q.Get("to")
	valueStr := q.Get("value")

	if from == "" || to == "" || valueStr == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing param: from, to and value are required"})
		return
	}

	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid value: must be a number"})
		return
	}

	resp, err := h.convertOnlyUC.Execute(usecase.ConvertOnlyRequest{
		Value: value,
		From:  temperature.Unit(from),
		To:    temperature.Unit(to),
	})
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]float64{"result": resp.Result})
}

func (h *Handler) handleHistory(w http.ResponseWriter, r *http.Request) {
	records, err := h.historyUC.Execute()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, records)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v) //nolint:errcheck
}
