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
