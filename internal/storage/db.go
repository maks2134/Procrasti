package storage

import (
	"database/sql"
	"fmt"
	"procrastigo/internal/models"
	"procrastigo/pkg/utils"
	"strings"
	_ "time"

	_ "github.com/lib/pq"
)

// PostgresStorage реализует Storage с использованием PostgreSQL
type PostgresStorage struct {
	db *sql.DB
}

// NewPostgresStorage создает новое хранилище PostgreSQL и проверяет подключение.
func NewPostgresStorage(dsn string) (*PostgresStorage, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Установка миграций/создание таблицы
	if err := createExcusesTable(db); err != nil {
		return nil, fmt.Errorf("failed to create excuses table: %w", err)
	}

	return &PostgresStorage{db: db}, nil
}

// createExcusesTable создает таблицу excuses, если она не существует
func createExcusesTable(db *sql.DB) error {
	query := `
    CREATE TABLE IF NOT EXISTS excuses (
        id VARCHAR(50) PRIMARY KEY,
        text TEXT NOT NULL,
        category VARCHAR(50) NOT NULL,
        language VARCHAR(10) NOT NULL,
        severity VARCHAR(20) NOT NULL,
        created_at TIMESTAMP WITH TIME ZONE NOT NULL,
        rating INTEGER DEFAULT 0
    );`
	_, err := db.Exec(query)
	return err
}

// LoadFromFile в PostgresStorage не используется, но должен быть реализован для интерфейса.
func (s *PostgresStorage) LoadFromFile(filename string) error {
	// В реальном приложении здесь может быть логика импорта данных
	return nil
}

// GetRandomExcuse получает случайное оправдание из БД
func (s *PostgresStorage) GetRandomExcuse() (*models.Excuse, error) {
	var excuse models.Excuse
	query := `
    SELECT id, text, category, language, severity, created_at, rating
    FROM excuses
    ORDER BY RANDOM()
    LIMIT 1`

	row := s.db.QueryRow(query)
	// Обязательно сканируем все поля, включая Rating
	err := row.Scan(&excuse.ID, &excuse.Text, &excuse.Category, &excuse.Language, &excuse.Severity, &excuse.CreatedAt, &excuse.Rating)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to scan random excuse: %w", err)
	}

	return &excuse, nil
}

// GetExcuses получает список оправданий с фильтрацией и сортировкой по рейтингу
func (s *PostgresStorage) GetExcuses(category, language string, limit int) ([]models.Excuse, error) {
	var args []interface{}
	argCounter := 1

	// Динамическое построение WHERE части запроса
	whereClauses := []string{}

	if category != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("category = $%d", argCounter))
		args = append(args, category)
		argCounter++
	}
	if language != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("language = $%d", argCounter))
		args = append(args, language)
		argCounter++
	}

	// Собираем запрос
	sqlQuery := "SELECT id, text, category, language, severity, created_at, rating FROM excuses"
	if len(whereClauses) > 0 {
		sqlQuery += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	// Сортируем по рейтингу и ограничиваем
	sqlQuery += fmt.Sprintf(" ORDER BY rating DESC, created_at DESC LIMIT $%d", argCounter)
	args = append(args, limit)

	rows, err := s.db.Query(sqlQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query excuses: %w", err)
	}
	defer rows.Close()

	var excuses []models.Excuse
	for rows.Next() {
		var excuse models.Excuse
		// Обязательно сканируем все поля, включая Rating
		if err := rows.Scan(&excuse.ID, &excuse.Text, &excuse.Category, &excuse.Language, &excuse.Severity, &excuse.CreatedAt, &excuse.Rating); err != nil {
			return nil, fmt.Errorf("failed to scan excuse row: %w", err)
		}
		excuses = append(excuses, excuse)
	}

	return excuses, nil
}

// CreateExcuse создает новое оправдание
func (s *PostgresStorage) CreateExcuse(excuse models.Excuse) error {
	query := `
    INSERT INTO excuses (id, text, category, language, severity, created_at, rating)
    VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := s.db.Exec(query,
		excuse.ID,
		excuse.Text,
		excuse.Category,
		excuse.Language,
		excuse.Severity,
		excuse.CreatedAt,
		excuse.Rating, // <--- Вставляем рейтинг
	)
	if err != nil {
		return fmt.Errorf("failed to insert excuse: %w", err)
	}
	return nil
}

// RateExcuse обновляет рейтинг оправдания
func (s *PostgresStorage) RateExcuse(id string, change int) error {
	query := `
    UPDATE excuses
    SET rating = rating + $1
    WHERE id = $2`

	res, err := s.db.Exec(query, change, id)
	if err != nil {
		return fmt.Errorf("failed to update excuse rating: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows // Оправдание не найдено
	}

	return nil
}

// GetStats вычисляет и возвращает статистику
func (s *PostgresStorage) GetStats() (*models.Stats, error) {
	stats := &models.Stats{}

	// 1. Общее количество
	if err := s.db.QueryRow("SELECT COUNT(*) FROM excuses").Scan(&stats.TotalExcuses); err != nil {
		return nil, fmt.Errorf("failed to get total excuses: %w", err)
	}

	// 2. Самая популярная категория (используем COALESCE для случая, когда нет данных)
	var mostPopular sql.NullString
	err := s.db.QueryRow(`
    SELECT category FROM excuses
    GROUP BY category
    ORDER BY COUNT(*) DESC
    LIMIT 1`).Scan(&mostPopular)

	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get most popular category: %w", err)
	}
	stats.MostPopularCategory = mostPopular.String

	// 3. Оправдания за сегодня
	today := utils.GetStartOfDay()
	if err := s.db.QueryRow(`
    SELECT COUNT(*) FROM excuses
    WHERE created_at >= $1`, today).Scan(&stats.ExcusesToday); err != nil {
		return nil, fmt.Errorf("failed to get excuses today: %w", err)
	}

	// 4. Уровень прокрастинации
	stats.GlobalProcrastinationLevel = utils.CalculateProcrastinationLevel(stats.ExcusesToday)

	return stats, nil
}

var _ Storage = (*PostgresStorage)(nil)
