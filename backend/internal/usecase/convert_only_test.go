package usecase_test

import (
	"testing"

	"github.com/gallardo/temp-converter/internal/domain/temperature"
	"github.com/gallardo/temp-converter/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConvertOnly_Execute(t *testing.T) {
	uc := usecase.NewConvertOnly()
	resp, err := uc.Execute(usecase.ConvertOnlyRequest{
		Value: 0,
		From:  temperature.Celsius,
		To:    temperature.Fahrenheit,
	})
	require.NoError(t, err)
	assert.InDelta(t, 32.0, resp.Result, 0.001)
}

func TestConvertOnly_Execute_InvalidUnit(t *testing.T) {
	uc := usecase.NewConvertOnly()
	_, err := uc.Execute(usecase.ConvertOnlyRequest{
		Value: 0,
		From:  "rankine",
		To:    temperature.Celsius,
	})
	assert.Error(t, err)
}
