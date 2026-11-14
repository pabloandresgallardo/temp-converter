package usecase

import (
	"fmt"
	"time"

	"github.com/gallardo/temp-converter/internal/domain/temperature"
	"github.com/gallardo/temp-converter/internal/port"
)

// ConvertTemperature es el caso de uso principal: convierte y persiste en historial.
type ConvertTemperature struct {
	repo port.HistoryRepository
}

func NewConvertTemperature(repo port.HistoryRepository) *ConvertTemperature {
	return &ConvertTemperature{repo: repo}
}

type ConvertRequest struct {
	Value float64
	From  temperature.Unit
	To    temperature.Unit
}

type ConvertResponse struct {
	Result float64
}

func (uc *ConvertTemperature) Execute(req ConvertRequest) (ConvertResponse, error) {
	result, err := temperature.Convert(req.Value, req.From, req.To)
	if err != nil {
		return ConvertResponse{}, fmt.Errorf("convert_temperature: %w", err)
	}
	record := port.ConversionRecord{
		Value:     req.Value,
		From:      req.From,
		To:        req.To,
		Result:    result,
		Timestamp: time.Now().Unix(),
	}
	if err := uc.repo.Save(record); err != nil {
		return ConvertResponse{}, fmt.Errorf("convert_temperature: %w", err)
	}
	return ConvertResponse{Result: result}, nil
}
