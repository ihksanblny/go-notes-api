package notes

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

func writeAPIError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := map[string]interface{}{
		"error": map[string]string{
			"code":    code,
			"message": message,
		},
	}

	_ = json.NewEncoder(w).Encode(resp)
}

type Handler struct {
	service *Service
}

func NewHandler(store Store) *Handler {
	return &Handler{
		service: NewService(store),
	}
}

// util response
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}

func ensureJSON(w http.ResponseWriter, r *http.Request) bool {
	ct := r.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "application/json") {
		writeAPIError(w, http.StatusUnsupportedMediaType, "UNSUPPORTED_CONTENT_TYPE", "Content-Type must be application/json")
		return false
	}
	return true
}

// /notes
func (h *Handler) HandleNotes(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.handleListNotes(w, r)
	case http.MethodPost:
		h.handleCreateNote(w, r)
	case http.MethodOptions:
		w.WriteHeader(http.StatusNoContent)
	default:
		writeAPIError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
	}
}

// /notes/{id}
func (h *Handler) HandleNotesByID(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/notes/")
	path = strings.TrimSuffix(path, "/")

	// Jika ada slash lagi, maka invalid
	if strings.Contains(path, "/") || path == "" {
		writeAPIError(w, http.StatusBadRequest, "INVALID_NOTE_PATH", "path must be /notes/{id}")
		return
	}

	id, err := strconv.Atoi(path)
	if err != nil {
		writeAPIError(w, http.StatusBadRequest, "INVALID_NOTE_ID", "id must be a valid integer")
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.handleGetNote(w, r, id)
	case http.MethodPut:
		h.handleUpdateNote(w, r, id)
	case http.MethodDelete:
		h.handleDeleteNote(w, r, id)
	default:
		writeAPIError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
	}
}

func (h *Handler) handleListNotes(w http.ResponseWriter, r *http.Request) {
	q := strings.TrimSpace(r.URL.Query().Get("q"))

	page := 1
	limit := 100

	// parse page
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		p, err := strconv.Atoi(pageStr)
		if err != nil || p < 1 {
			writeAPIError(w, http.StatusBadRequest, "INVALID_PAGE", "page must be a positive integer")
			return
		}
		page = p
	}

	// parse limit
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		l, err := strconv.Atoi(limitStr)
		if err != nil || l < 1 || l > 100 {
			writeAPIError(w, http.StatusBadRequest, "INVALID_LIMIT", "limit must be a positive integer between 1 and 100")
			return
		}
		limit = l
	}

	sort := r.URL.Query().Get("sort")
	order := r.URL.Query().Get("order")

	items, total, err := h.service.ListNotes(r.Context(), q, page, limit, sort, order)
	if err != nil {
		writeAPIError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to list notes")
		return
	}

	resp := map[string]interface{}{
		"data":  items,
		"total": total,
		"page":  page,
		"limit": limit,
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *Handler) handleCreateNote(w http.ResponseWriter, r *http.Request) {
	if !ensureJSON(w, r) {
		return
	}

	var input struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeAPIError(w, http.StatusBadRequest, "INVALID_REQUEST_BODY", "invalid request body")
		return
	}

	title := strings.TrimSpace(input.Title)
	content := strings.TrimSpace(input.Content)

	note, err := h.service.CreateNote(r.Context(), title, content)
	if err != nil {
		writeAPIError(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, note)
}

func (h *Handler) handleGetNote(w http.ResponseWriter, r *http.Request, id int) {
	note, err := h.service.GetNote(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNoteNotFound) {
			writeAPIError(w, http.StatusNotFound, "NOT_FOUND", "note not found")
		} else {
			writeAPIError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to get note")
		}
		return
	}
	writeJSON(w, http.StatusOK, note)
}

func (h *Handler) handleUpdateNote(w http.ResponseWriter, r *http.Request, id int) {
	if !ensureJSON(w, r) {
		return
	}

	var input struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeAPIError(w, http.StatusBadRequest, "INVALID_REQUEST_BODY", "invalid request body")
		return
	}

	title := strings.TrimSpace(input.Title)
	content := strings.TrimSpace(input.Content)

	note, err := h.service.UpdateNote(r.Context(), id, title, content)
	if err != nil {
		switch {
		case errors.Is(err, ErrNoteNotFound):
			writeAPIError(w, http.StatusNotFound, "NOT_FOUND", "note not found")
		default:
			writeAPIError(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		}
		return
	}

	writeJSON(w, http.StatusOK, note)
}

func (h *Handler) handleDeleteNote(w http.ResponseWriter, r *http.Request, id int) {
	err := h.service.DeleteNote(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNoteNotFound) {
			writeAPIError(w, http.StatusNotFound, "NOT_FOUND", "note not found")
		} else {
			writeAPIError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to delete note")
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
