package usecase_test

import (
	"testing"

	"github.com/gallardo/temp-converter/internal/adapter/memory"
	"github.com/gallardo/temp-converter/internal/port"
	"github.com/gallardo/temp-converter/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetHistory_Execute_Empty(t *testing.T) {
	repo := memory.NewHistoryRepository()
	uc := usecase.NewGetHistory(repo)
	records, err := uc.Execute()
	require.NoError(t, err)
	assert.NotNil(t, records)
	assert.Empty(t, records)
}

func TestGetHistory_Execute_WithRecords(t *testing.T) {
	repo := memory.NewHistoryRepository()
	_ = repo.Save(port.ConversionRecord{Value: 0, From: "celsius", To: "fahrenheit", Result: 32})
	uc := usecase.NewGetHistory(repo)
	records, err := uc.Execute()
	require.NoError(t, err)
	assert.Len(t, records, 1)
}
