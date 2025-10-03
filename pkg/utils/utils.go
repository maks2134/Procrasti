package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Инициализация генератора случайных чисел
func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomInt генерирует случайное целое число в диапазоне [0, max-1].
func RandomInt(max int) int {
	if max <= 0 {
		return 0
	}
	return rand.Intn(max)
}

// GetStartOfDay возвращает начало текущего дня (UTC).
func GetStartOfDay() time.Time {
	now := time.Now().UTC()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
}

// GenerateID создает уникальный ID с префиксом.
func GenerateID(prefix string) string {
	timestamp := time.Now().UnixNano()
	return fmt.Sprintf("%s-%d", prefix, timestamp)
}

// --- Validation and Parsing ---

// ValidateLanguage проверяет, является ли язык допустимым.
func ValidateLanguage(lang string) bool {
	validLangs := map[string]bool{"ru": true, "en": true}
	return validLangs[strings.ToLower(lang)]
}

// ValidateCategory проверяет, является ли категория допустимой.
func ValidateCategory(category string) bool {
	validCats := map[string]bool{
		"general": true,
		"work":    true,
		"study":   true,
		"social":  true,
		"health":  true,
	}
	return validCats[strings.ToLower(category)]
}

// ValidateSeverity проверяет, является ли важность допустимой.
func ValidateSeverity(severity string) bool {
	validSeverities := map[string]bool{
		"low":    true,
		"medium": true,
		"high":   true,
	}
	return validSeverities[strings.ToLower(severity)]
}

// ParseLimit парсит строковый лимит в целое число, используя значение по умолчанию.
func ParseLimit(limitStr string, defaultLimit int) int {
	if limitStr == "" {
		return defaultLimit
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 0 {
		return defaultLimit
	}
	return limit
}

// CalculateProcrastinationLevel оценивает уровень прокрастинации.
func CalculateProcrastinationLevel(excusesToday int) string {
	if excusesToday > 100 {
		return "Critical"
	} else if excusesToday > 50 {
		return "High"
	} else if excusesToday > 10 {
		return "Medium"
	}
	return "Low"
}

// --- HTTP Helpers ---

// JSONResponse отправляет ответ в формате JSON.
func JSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding JSON: %v", err), http.StatusInternalServerError)
	}
}

// ErrorResponse отправляет ответ об ошибке в формате JSON.
func ErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	JSONResponse(w, statusCode, map[string]string{"error": message})
}

// JSONDecode декодирует JSON из тела запроса.
func JSONDecode(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}
