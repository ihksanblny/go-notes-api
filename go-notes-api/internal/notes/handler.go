package notes

import (
	"encoding/json"
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
	store Store
}

func NewHandler(store Store) *Handler {
	return &Handler{store: store}
}

// util response
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
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

	// Jika ada slash lagi, maka invalid
	if strings.Contains(path, "/") || path == "" {
		writeAPIError(w, http.StatusBadRequest, "INVALID_NOTE_PATH", "path must be /notes/{id}")
		return
	}
	
	id, err := strconv.Atoi(path)
	if err != nil {
		writeAPIError(w, http.StatusBadRequest, "INVALID_NOTE_ID", "id must be number")
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
	limit := 10

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

	items, total := h.store.ListPage(page, limit, q)

	resp := map[string]interface{}{
		"data":  items,
		"total": total,
		"page":  page,
		"limit": limit,
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *Handler) handleCreateNote(w http.ResponseWriter, r *http.Request) {
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

	if err := ValidateNoteInput(title, content); err != nil {
		writeAPIError(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	note := h.store.Create(title, content)
	writeJSON(w, http.StatusCreated, note)
}

func (h *Handler) handleGetNote(w http.ResponseWriter, r *http.Request, id int) {
	note, ok := h.store.Get(id)
	if !ok {
		writeAPIError(w, http.StatusNotFound, "NOT_FOUND", "note not found")
		return
	}
	writeJSON(w, http.StatusOK, note)
}

func (h *Handler) handleUpdateNote(w http.ResponseWriter, r *http.Request, id int) {
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

	if err := ValidateNoteInput(title, content); err != nil {
		writeAPIError(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	note, ok := h.store.Update(id, title, content)
	if !ok {
		writeAPIError(w, http.StatusNotFound, "NOT_FOUND", "note not found")
		return
	}

	writeJSON(w, http.StatusOK, note)
}

func (h *Handler) handleDeleteNote(w http.ResponseWriter, r *http.Request, id int) {
	ok := h.store.Delete(id)
	if !ok {
		writeAPIError(w, http.StatusNotFound, "NOT_FOUND", "note not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
