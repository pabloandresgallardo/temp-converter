package usecase_test

import (
	"testing"

	"github.com/gallardo/temp-converter/internal/adapter/memory"
	"github.com/gallardo/temp-converter/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConvertTemperature_Execute(t *testing.T) {
	repo := memory.NewHistoryRepository()
	uc := usecase.NewConvertTemperature(repo)
	resp, err := uc.Execute(usecase.ConvertRequest{Value: 0, From: "celsius", To: "fahrenheit"})
	require.NoError(t, err)
	assert.InDelta(t, 32.0, resp.Result, 0.001)
	saved, _ := repo.FindAll()
	assert.Len(t, saved, 1)
	assert.InDelta(t, 32.0, saved[0].Result, 0.001)
	assert.Equal(t, "celsius", string(saved[0].From))
	assert.Equal(t, "fahrenheit", string(saved[0].To))
}

func TestConvertTemperature_Execute_InvalidUnit(t *testing.T) {
	repo := memory.NewHistoryRepository()
	uc := usecase.NewConvertTemperature(repo)
	_, err := uc.Execute(usecase.ConvertRequest{Value: 0, From: "celsius", To: "bad"})
	require.Error(t, err)
	saved, _ := repo.FindAll()
	assert.Empty(t, saved)
}
