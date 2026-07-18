package models

import "time"

type Bookmark struct {
	ID        int       `json:"id"`
	URL       string    `json:"url"`
	Title     string    `json:"title"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
}
