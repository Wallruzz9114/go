package plethora_api

import (
	uuid "github.com/satori/go.uuid"
)

// Location represents company location model
type Location struct {
	Base
	Name      string    `json:"name"`
	Active    bool      `json:"active"`
	Address   string    `json:"address"`
	CompanyID uuid.UUID `json:"company_id"`
}
