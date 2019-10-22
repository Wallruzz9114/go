package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Post struct {
	ID        uuid.UUID `gorm:"primary_key;unique;type:uuid; column:id;default:gen_random_uuid()" json:"id"`
	Title     string    `gorm:"size:255;not null;unique" json:"title"`
	Content   string    `gorm:"size:255;not null;" json:"content"`
	Author    User      `json:"author"`
	AuthorID  uuid.UUID `gorm:"not null" json:"author_id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
