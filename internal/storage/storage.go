package storage

import "procrastigo/internal/models"

type Storage interface {
	GetRandomExcuse() (*models.Excuse, error)
	GetExcuses(category, language string, limit int) ([]models.Excuse, error)
	CreateExcuse(excuse models.Excuse) error
	GetStats() (*models.Stats, error)
	LoadFromFile(filename string) error
}
