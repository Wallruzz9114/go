package plethora_api

import (
	"context"
	"time"

	"github.com/go-pg/pg/orm"
	uuid "github.com/satori/go.uuid"
)

// Base contains common fields for all tables
type Base struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at,omitempty" pg:",soft_delete"`
}

// ListQuery holds company/location data used for list db queries
type ListQuery struct {
	Query string
	ID    uuid.UUID
}

// BeforeInsert hooks into insert operations, setting createdAt and updatedAt to current time
func (b *Base) BeforeInsert(_ context.Context, _ orm.DB) error {
	now := time.Now()
	b.CreatedAt = now
	b.UpdatedAt = now
	return nil
}

// BeforeUpdate hooks into update operations, setting updatedAt to current time
func (b *Base) BeforeUpdate(_ context.Context, _ orm.DB) error {
	b.UpdatedAt = time.Now()
	return nil
}
