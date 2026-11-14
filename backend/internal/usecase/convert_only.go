package usecase

import (
	"fmt"

	"github.com/gallardo/temp-converter/internal/domain/temperature"
)

type ConvertOnlyRequest struct {
	Value float64
	From  temperature.Unit
	To    temperature.Unit
}

type ConvertOnlyResponse struct {
	Result float64
}

type ConvertOnly struct{}

func NewConvertOnly() *ConvertOnly {
	return &ConvertOnly{}
}

func (uc *ConvertOnly) Execute(req ConvertOnlyRequest) (ConvertOnlyResponse, error) {
	result, err := temperature.Convert(req.Value, req.From, req.To)
	if err != nil {
		return ConvertOnlyResponse{}, fmt.Errorf("convert_only: %w", err)
	}
	return ConvertOnlyResponse{Result: result}, nil
}
