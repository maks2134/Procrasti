package storage

import (
	"errors"
	"fmt"
	"io/ioutil"
	"procrastigo/internal/models"
	"procrastigo/pkg/utils"
	"sync"
	"time"

	"encoding/json"
)

// ExcuseFileFormat - структура для десериализации данных из JSON файла
type ExcuseFileFormat struct {
	Excuses []models.Excuse `json:"excuses"`
}

// MemoryStorage - простое хранилище в оперативной памяти
type MemoryStorage struct {
	excuses map[string]models.Excuse
	mu      sync.RWMutex
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		excuses: make(map[string]models.Excuse),
	}
}

// LoadFromFile загружает оправдания из JSON файла
func (s *MemoryStorage) LoadFromFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	var fileData ExcuseFileFormat
	if err := json.Unmarshal(data, &fileData); err != nil {
		return fmt.Errorf("json: cannot unmarshal object into Go value of type %T: %w", fileData, err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	for _, excuse := range fileData.Excuses {
		s.excuses[excuse.ID] = excuse
	}
	return nil
}

func (s *MemoryStorage) GetRandomExcuse() (*models.Excuse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.excuses) == 0 {
		return nil, nil
	}

	keys := make([]string, 0, len(s.excuses))
	for k := range s.excuses {
		keys = append(keys, k)
	}

	randomKey := keys[utils.RandomInt(len(keys))]
	excuse := s.excuses[randomKey]

	return &excuse, nil
}

func (s *MemoryStorage) GetExcuses(category, language string, limit int) ([]models.Excuse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []models.Excuse
	count := 0
	for _, excuse := range s.excuses {
		matches := true
		if category != "" && excuse.Category != category {
			matches = false
		}
		if language != "" && excuse.Language != language {
			matches = false
		}

		if matches {
			result = append(result, excuse)
			count++
			if limit > 0 && count >= limit {
				break
			}
		}
	}
	return result, nil
}

func (s *MemoryStorage) CreateExcuse(excuse models.Excuse) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.excuses[excuse.ID]; exists {
		return errors.New("excuse already exists")
	}

	s.excuses[excuse.ID] = excuse
	return nil
}

// RateExcuse - НОВАЯ РЕАЛИЗАЦИЯ ДЛЯ УДОВЛЕТВОРЕНИЯ ИНТЕРФЕЙСУ
func (s *MemoryStorage) RateExcuse(id string, change int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	excuse, exists := s.excuses[id]
	if !exists {
		return errors.New("excuse not found")
	}

	// Обновляем рейтинг в памяти
	excuse.Rating += change
	s.excuses[id] = excuse

	return nil
}

func (s *MemoryStorage) GetStats() (*models.Stats, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	stats := &models.Stats{
		TotalExcuses: len(s.excuses),
		ExcusesToday: 0,
	}

	now := time.Now().Truncate(24 * time.Hour)
	categoryCounts := make(map[string]int)
	maxCount := 0

	for _, excuse := range s.excuses {
		// Статистика за сегодня
		if excuse.CreatedAt.Truncate(24 * time.Hour).Equal(now) {
			stats.ExcusesToday++
		}

		// Самая популярная категория
		categoryCounts[excuse.Category]++
		if categoryCounts[excuse.Category] > maxCount {
			maxCount = categoryCounts[excuse.Category]
			stats.MostPopularCategory = excuse.Category
		}
	}

	stats.GlobalProcrastinationLevel = utils.CalculateProcrastinationLevel(stats.ExcusesToday)

	return stats, nil
}

// Утверждение, что *MemoryStorage реализует Storage
var _ Storage = (*MemoryStorage)(nil)
