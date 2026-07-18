package storage

import (
	"bookmark-keeper/internal/models"
	"sync"
	"time"
)

type Storage struct {
	bookmarks map[int]models.Bookmark
	mu        sync.RWMutex
	nextID    int
}

func NewStorage() *Storage {
	return &Storage{
		bookmarks: make(map[int]models.Bookmark),
		nextID:    1,
	}
}

func (s *Storage) Create(url, title string, tags []string) models.Bookmark {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := s.nextID

	if tags == nil {
		tags = make([]string, 0)
	}

	createMark := models.Bookmark{
		ID:        id,
		URL:       url,
		Title:     title,
		Tags:      tags,
		CreatedAt: time.Now(),
	}

	s.nextID++
	s.bookmarks[id] = createMark
	return createMark
}

func (s *Storage) GetByID(id int) (models.Bookmark, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	value, ok := s.bookmarks[id]
	return value, ok
}

func (s *Storage) GetAll() []models.Bookmark {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]models.Bookmark, 0)
	for _, bookmark := range s.bookmarks {
		result = append(result, bookmark)
	}

	return result
}

func (s *Storage) Delete(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.bookmarks[id]
	if !ok {
		return false
	}
	delete(s.bookmarks, id)
	return true
}
