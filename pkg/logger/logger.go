package logger

import (
	"io"
	"log"
	"os"
	"procrastigo/internal/models"
)

var (
	Info  *log.Logger
	Warn  *log.Logger
	Error *log.Logger
	Debug *log.Logger
)

func Init(level string) {
	flags := log.Ldate | log.Ltime | log.Lshortfile

	Info = log.New(os.Stdout, "INFO: ", flags)
	Warn = log.New(os.Stdout, "WARN: ", flags)
	Error = log.New(os.Stderr, "ERROR: ", flags)
	Debug = log.New(os.Stdout, "DEBUG: ", flags)

	// В продакшене можно отключить debug логи
	if level == "production" {
		Debug.SetOutput(io.Discard)
	}
}

func LogExcuseRequest(excuse *models.Excuse, method string) {
	Info.Printf("Excuse %s: %s - Category: %s, Severity: %s",
		method, excuse.ID, excuse.Category, excuse.Severity)
}

func LogAPIRequest(method, path, remoteAddr string) {
	Debug.Printf("API Request: %s %s from %s", method, path, remoteAddr)
}

func LogStatsRequest() {
	Debug.Printf("Stats endpoint accessed")
}
