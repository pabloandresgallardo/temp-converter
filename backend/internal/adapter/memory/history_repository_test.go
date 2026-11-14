package memory_test

import (
	"testing"

	"github.com/gallardo/temp-converter/internal/adapter/memory"
	"github.com/gallardo/temp-converter/internal/port"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHistoryRepository_FindAll_Empty(t *testing.T) {
	repo := memory.NewHistoryRepository()
	all, err := repo.FindAll()
	require.NoError(t, err)
	assert.NotNil(t, all)
	assert.Empty(t, all)
}

func TestHistoryRepository_SaveAndFindAll(t *testing.T) {
	repo := memory.NewHistoryRepository()
	record := port.ConversionRecord{Value: 0, From: "celsius", To: "fahrenheit", Result: 32, Timestamp: 1}
	require.NoError(t, repo.Save(record))
	all, err := repo.FindAll()
	require.NoError(t, err)
	assert.Len(t, all, 1)
	assert.Equal(t, record, all[0])
}

func TestHistoryRepository_FindAll_LIFO(t *testing.T) {
	repo := memory.NewHistoryRepository()
	r1 := port.ConversionRecord{Value: 0, From: "celsius", To: "fahrenheit", Result: 32, Timestamp: 1}
	r2 := port.ConversionRecord{Value: 100, From: "celsius", To: "kelvin", Result: 373.15, Timestamp: 2}
	require.NoError(t, repo.Save(r1))
	require.NoError(t, repo.Save(r2))
	all, err := repo.FindAll()
	require.NoError(t, err)
	assert.Len(t, all, 2)
	assert.Equal(t, r2, all[0]) // más reciente primero
	assert.Equal(t, r1, all[1])
}
