package handlers

import (
	"bookmark-keeper/internal/storage"
	"encoding/json"
	"net/http"
)

type Handler struct {
	storage *storage.Storage
}

func NewHandler(s *storage.Storage) *Handler {
	return &Handler{storage: s}
}

type input struct {
	URL   string   `json:"url"`
	Title string   `json:"title"`
	Tags  []string `json:"tags"`
}

func (h *Handler) CreateBookmark(w http.ResponseWriter, r *http.Request) {
	var req input
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !validURL(req.URL) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if req.Title == "" || len(req.Title) > 100 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	bm := h.storage.Create(req.URL, req.Title, req.Tags)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(bm)
}

func (h *Handler) ListBookmarks(w http.ResponseWriter, r *http.Request) {
	bookmarks := h.storage.GetAll()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bookmarks)
}

func (h *Handler) GetBookmark(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := parseID(idStr)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	bookmark, ok := h.storage.GetByID(id)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bookmark)
}

func (h *Handler) DeleteBookmark(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := parseID(idStr)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	ok := h.storage.Delete(id)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
