package memory

import (
	"sync"

	"github.com/gallardo/temp-converter/internal/port"
)

// HistoryRepository implementa port.HistoryRepository con almacenamiento in-memory.
type HistoryRepository struct {
	mu      sync.RWMutex
	records []port.ConversionRecord
}

func NewHistoryRepository() *HistoryRepository {
	return &HistoryRepository{}
}

func (r *HistoryRepository) Save(record port.ConversionRecord) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.records = append([]port.ConversionRecord{record}, r.records...)
	return nil
}

func (r *HistoryRepository) FindAll() ([]port.ConversionRecord, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]port.ConversionRecord, len(r.records))
	copy(result, r.records)
	return result, nil
}
