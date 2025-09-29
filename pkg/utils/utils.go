package utils

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"procrastigo/internal/models"
	"strings"
	"time"
)

func GenerateID(prefix string) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, 8)

	for i := range result {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		result[i] = charset[num.Int64()]
	}

	return prefix + "_" + string(result)
}

func ValidateCategory(category string) bool {
	validCategories := map[string]bool{
		"work":    true,
		"family":  true,
		"tech":    true,
		"urgent":  true,
		"general": true,
	}

	return validCategories[strings.ToLower(category)]
}

func ValidateSeverity(severity string) bool {
	validSeverities := map[string]bool{
		"low":      true,
		"medium":   true,
		"high":     true,
		"critical": true,
	}

	return validSeverities[strings.ToLower(severity)]
}

func ValidateLanguage(lang string) bool {
	validLanguages := map[string]bool{
		"ru": true,
		"en": true,
		"es": true,
		"fr": true,
	}

	return validLanguages[strings.ToLower(lang)]
}

func JSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func ErrorResponse(w http.ResponseWriter, status int, message string) {
	JSONResponse(w, status, map[string]string{
		"error":   http.StatusText(status),
		"message": message,
	})
}

func ParseLimit(limitStr string, defaultLimit int) int {
	if limitStr == "" {
		return defaultLimit
	}

	var limit int
	if _, err := fmt.Sscanf(limitStr, "%d", &limit); err != nil {
		return defaultLimit
	}

	if limit > 100 {
		limit = 100
	}

	return limit
}

func CalculateProcrastinationLevel(excusesToday int) string {
	switch {
	case excusesToday > 20:
		return "apocalyptic"
	case excusesToday > 10:
		return "critical"
	case excusesToday > 5:
		return "high"
	case excusesToday > 2:
		return "medium"
	default:
		return "low"
	}
}

func GetStartOfDay() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
}

func FilterExcuses(excuses []models.Excuse, category, language, severity string) []models.Excuse {
	var filtered []models.Excuse

	for _, excuse := range excuses {
		match := true

		if category != "" && !strings.EqualFold(excuse.Category, category) {
			match = false
		}

		if language != "" && !strings.EqualFold(excuse.Language, language) {
			match = false
		}

		if severity != "" && !strings.EqualFold(excuse.Severity, severity) {
			match = false
		}

		if match {
			filtered = append(filtered, excuse)
		}
	}

	return filtered
}

func JSONDecode(body io.ReadCloser, v interface{}) error {
	defer body.Close()
	return json.NewDecoder(body).Decode(v)
}

func GetClientIP(r *http.Request) string {
	ips := r.Header.Get("X-Forwarded-For")
	if ips != "" {
		return strings.Split(ips, ",")[0]
	}

	ips = r.Header.Get("X-Real-IP")
	if ips != "" {
		return ips
	}

	return r.RemoteAddr
}
