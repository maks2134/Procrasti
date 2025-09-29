package storage

import (
	"encoding/json"
	"math/rand"
	"os"
	"sync"
	"time"

	"procrastigo/internal/models"
	"procrastigo/pkg/utils"
)

type MemoryStorage struct {
	mu      sync.RWMutex
	excuses []models.Excuse
}

func NewMemoryStorage() *MemoryStorage {
	rand.Seed(time.Now().UnixNano())
	return &MemoryStorage{excuses: make([]models.Excuse, 0)}
}

func (m *MemoryStorage) LoadFromFile(filename string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	var data []models.Excuse
	if err := json.NewDecoder(f).Decode(&data); err != nil {
		return err
	}

	m.excuses = data
	return nil
}

func (m *MemoryStorage) GetRandomExcuse() (*models.Excuse, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if len(m.excuses) == 0 {
		return nil, nil
	}

	idx := rand.Intn(len(m.excuses))
	e := m.excuses[idx]
	return &e, nil
}

func (m *MemoryStorage) GetExcuses(category, language string, limit int) ([]models.Excuse, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	filtered := utils.FilterExcuses(m.excuses, category, language, "")
	if limit > len(filtered) {
		limit = len(filtered)
	}
	return filtered[:limit], nil
}

func (m *MemoryStorage) CreateExcuse(excuse models.Excuse) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.excuses = append(m.excuses, excuse)
	return nil
}

func (m *MemoryStorage) GetStats() (*models.Stats, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if len(m.excuses) == 0 {
		return &models.Stats{}, nil
	}

	counts := make(map[string]int)
	today := utils.GetStartOfDay()
	excusesToday := 0
	for _, e := range m.excuses {
		counts[e.Category]++
		if !e.CreatedAt.Before(today) {
			excusesToday++
		}
	}

	most := ""
	max := 0
	for k, v := range counts {
		if v > max {
			max = v
			most = k
		}
	}

	return &models.Stats{
		TotalExcuses:               len(m.excuses),
		MostPopularCategory:        most,
		ExcusesToday:               excusesToday,
		GlobalProcrastinationLevel: utils.CalculateProcrastinationLevel(excusesToday),
	}, nil
}

var _ Storage = (*MemoryStorage)(nil)
