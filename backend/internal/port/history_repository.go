package port

import "github.com/gallardo/temp-converter/internal/domain/temperature"

// ConversionRecord representa una conversión realizada.
type ConversionRecord struct {
	Value     float64          `json:"value"`
	From      temperature.Unit `json:"from"`
	To        temperature.Unit `json:"to"`
	Result    float64          `json:"result"`
	Timestamp int64            `json:"timestamp"`
}

// HistoryRepository es el puerto para acceder al historial de conversiones.
type HistoryRepository interface {
	Save(record ConversionRecord) error
	FindAll() ([]ConversionRecord, error)
}
