package storage

import (
	"bookmark-keeper/internal/models"
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"
)

type Storage struct {
	bookmarks map[int]models.Bookmark
	mu        sync.RWMutex
	nextID    int
}

func NewStorage() *Storage {
	s := &Storage{
		bookmarks: make(map[int]models.Bookmark),
		nextID:    1,
	}

	s.load()
	return s
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
	s.save()
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
	s.save()
	return true
}

func (s *Storage) save() {
	saveBookmarks := make([]models.Bookmark, 0)
	for _, bookmark := range s.bookmarks {
		saveBookmarks = append(saveBookmarks, bookmark)
	}

	data, err := json.Marshal(saveBookmarks)
	if err != nil {
		log.Println("Ошибка сериализации")
		return
	}

	if err := os.WriteFile("bookmarks.json", data, 0644); err != nil {
		log.Println("Ошибка записи")
		return
	}
}

func (s *Storage) load() {
	data, err := os.ReadFile("bookmarks.json")
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		log.Printf("ошибка чтения файла: %v", err)
		return
	}

	var loaded []models.Bookmark
	if err := json.Unmarshal(data, &loaded); err != nil {
		log.Printf("ошибка парсинга: %v", err)
		return
	}

	maxID := 0

	for _, b := range loaded {
		s.bookmarks[b.ID] = b
		if b.ID > maxID {
			maxID = b.ID
		}
	}
	s.nextID = maxID + 1
}
