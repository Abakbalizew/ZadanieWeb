package models

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	PostUUID       uuid.UUID
	AuthorUUID     uuid.UUID
	IdempotencyKey string
	Title          string
	Content        string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Status         string
}

// Ключ, уникальный для всех постов
func NewIdempotencyKey(title string, content string, userUUID string) string {
	return title + "123" + content + userUUID + "123"
}
