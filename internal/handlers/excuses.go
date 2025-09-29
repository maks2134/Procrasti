package handlers

import (
	"net/http"
	"procrastigo/internal/models"
	"procrastigo/internal/storage"
	"procrastigo/pkg/logger"
	"procrastigo/pkg/utils"
)

type ExcuseHandler struct {
	storage storage.Storage
}

func NewExcuseHandler(storage storage.Storage) *ExcuseHandler {
	return &ExcuseHandler{storage: storage}
}

func (h *ExcuseHandler) GetRandomExcuse(w http.ResponseWriter, r *http.Request) {
	lang := r.URL.Query().Get("lang")
	category := r.URL.Query().Get("category")

	if lang != "" && !utils.ValidateLanguage(lang) {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid language")
		return
	}

	if category != "" && !utils.ValidateCategory(category) {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid category")
		return
	}

	excuse, err := h.storage.GetRandomExcuse()
	if err != nil {
		logger.Error.Printf("Failed to get random excuse: %v", err)
		utils.ErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	if excuse == nil {
		utils.ErrorResponse(w, http.StatusNotFound, "No excuses found")
		return
	}

	logger.LogExcuseRequest(excuse, "RANDOM")
	utils.JSONResponse(w, http.StatusOK, excuse)
}

func (h *ExcuseHandler) GetExcuses(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	lang := r.URL.Query().Get("lang")
	severity := r.URL.Query().Get("severity")
	limit := utils.ParseLimit(r.URL.Query().Get("limit"), 1)

	// Валидация параметров
	if category != "" && !utils.ValidateCategory(category) {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid category")
		return
	}

	if lang != "" && !utils.ValidateLanguage(lang) {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid language")
		return
	}

	if severity != "" && !utils.ValidateSeverity(severity) {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid severity")
		return
	}

	excuses, err := h.storage.GetExcuses(category, lang, limit)
	if err != nil {
		logger.Error.Printf("Failed to get excuses: %v", err)
		utils.ErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	utils.JSONResponse(w, http.StatusOK, excuses)
}

func (h *ExcuseHandler) CreateExcuse(w http.ResponseWriter, r *http.Request) {
	var req models.ExcuseRequest
	if err := utils.JSONDecode(r.Body, &req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if req.Text == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, "Text is required")
		return
	}

	// Устанавливаем значения по умолчанию
	if req.Category == "" {
		req.Category = "general"
	}
	if req.Language == "" {
		req.Language = "ru"
	}
	if req.Severity == "" {
		req.Severity = "medium"
	}

	// Валидация
	if !utils.ValidateCategory(req.Category) {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid category")
		return
	}
	if !utils.ValidateLanguage(req.Language) {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid language")
		return
	}
	if !utils.ValidateSeverity(req.Severity) {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid severity")
		return
	}

	excuse := models.Excuse{
		ID:        utils.GenerateID("exc"),
		Text:      req.Text,
		Category:  req.Category,
		Language:  req.Language,
		Severity:  req.Severity,
		CreatedAt: utils.GetStartOfDay(),
	}

	if err := h.storage.CreateExcuse(excuse); err != nil {
		logger.Error.Printf("Failed to create excuse: %v", err)
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to create excuse")
		return
	}

	logger.LogExcuseRequest(&excuse, "CREATE")
	utils.JSONResponse(w, http.StatusCreated, excuse)
}
