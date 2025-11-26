package notes

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

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

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
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
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

// /notes/{id}
func (h *Handler) HandleNotesByID(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	// path: /notes/123
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) != 2 {
		writeError(w, http.StatusBadRequest, "invalid URL path")
		return
	}

	id, err := strconv.Atoi(parts[1])
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid note ID")
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
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *Handler) handleListNotes(w http.ResponseWriter, r *http.Request) {
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	pageStr := strings.TrimSpace(r.URL.Query().Get("page"))
	limitStr := strings.TrimSpace(r.URL.Query().Get("limit"))

	// Kalau semua kosong -> fallback ke List() lama
	if q == "" && pageStr == "" && limitStr == "" {
		notes := h.store.List()
		writeJSON(w, http.StatusOK, notes)
		return
	}

	page := 1
	limit := 10

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	items, total := h.store.ListPage(page, limit, q)

	resp := map[string]interface{} {
		"data": items,
		"total": total,
		"page": page,
		"limit": limit,
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *Handler) handleCreateNote(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title  string `json:"title"`
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	title := strings.TrimSpace(input.Title)
	content := strings.TrimSpace(input.Content)

	if err := ValidateNoteInput(title, content); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	note := h.store.Create(title, content)
	writeJSON(w, http.StatusCreated, note)
}

func (h* Handler) handleGetNote(w http.ResponseWriter, r *http.Request, id int) {
	note, ok := h.store.Get(id)
	if !ok {
		writeError(w, http.StatusNotFound, "note not found")
		return
	}
	writeJSON(w, http.StatusOK, note)
}

func (h* Handler) handleUpdateNote(w http.ResponseWriter, r *http.Request, id int) {
	var input struct {
		Title string `json:"title"`
		Content string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	title := strings.TrimSpace(input.Title)
	content := strings.TrimSpace(input.Content)

	if err := ValidateNoteInput(title, content); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	note, ok := h.store.Update(id, title, content)
	if !ok {
		writeError(w, http.StatusNotFound, "note not found")
		return
	}

	writeJSON(w, http.StatusOK, note)
}

func (h* Handler) handleDeleteNote(w http.ResponseWriter, r *http.Request, id int) {
	ok := h.store.Delete(id)
	if !ok {
		writeError(w, http.StatusNotFound, "note not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}