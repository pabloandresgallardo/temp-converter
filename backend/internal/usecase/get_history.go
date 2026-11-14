package usecase

import (
	"fmt"

	"github.com/gallardo/temp-converter/internal/port"
)

// GetHistory retorna el historial volátil de conversiones.
type GetHistory struct {
	repo port.HistoryRepository
}

func NewGetHistory(repo port.HistoryRepository) *GetHistory {
	return &GetHistory{repo: repo}
}

func (uc *GetHistory) Execute() ([]port.ConversionRecord, error) {
	records, err := uc.repo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("get_history: %w", err)
	}
	if records == nil {
		return []port.ConversionRecord{}, nil
	}
	return records, nil
}
