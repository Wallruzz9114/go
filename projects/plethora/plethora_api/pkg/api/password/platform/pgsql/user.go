package pgsql

import (
	"github.com/go-pg/pg/orm"

	plethora_api "github.com/Wallruzz9114/plethora_api/pkg/util/model"
	uuid "github.com/satori/go.uuid"
)

// User represents the client for user table
type User struct{}

// NewUser returns a new user database instance
func NewUser() *User {
	return &User{}
}

// View returns single user by ID
func (u *User) View(db orm.DB, id uuid.UUID) (*plethora_api.User, error) {
	user := &plethora_api.User{Base: plethora_api.Base{ID: id}}
	err := db.Select(user)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// Update updates user's info
func (u *User) Update(db orm.DB, user *plethora_api.User) error {
	return db.Update(user)
}
