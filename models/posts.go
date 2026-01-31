package models

///Модель поста

import (
	"time"

	"github.com/google/uuid"
)

// Модель поста
type Post struct {
	PostUUID       uuid.UUID
	AuthorUUID     uuid.UUID
	IdempotencyKey string
	Title          string
	Content        string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	//Поле LastEditedAt равно более позднему из CreatedAt и UpdatedAt,
	//это необходимо, чтобы отследить, что было последнее и отобразить именно
	//эту дату в публикации/изменении поста
	LastEditedAt time.Time
	Status       string //="Published" || ="Draft"
	//url-адрес картинки
	ImageUrl    string
	AuthorEmail string
}

// Ключ, уникальный для всех постов
func NewIdempotencyKey(title string, content string, userUUID string) string {
	return title + "123" + content + userUUID + "123"
}
