package mockdb

import (
	"github.com/go-pg/pg/orm"

	plethora_api "github.com/Wallruzz9114/plethora_api/pkg/util/model"
	uuid "github.com/satori/go.uuid"
)

// User database mock
type User struct {
	CreateFunction         func(orm.DB, plethora_api.User) (*plethora_api.User, error)
	ViewFunction           func(orm.DB, uuid.UUID) (*plethora_api.User, error)
	FindByUsernameFunction func(orm.DB, string) (*plethora_api.User, error)
	FindByTokenFunction    func(orm.DB, string) (*plethora_api.User, error)
	ListFunction           func(orm.DB, *plethora_api.ListQuery, *plethora_api.Pagination) ([]plethora_api.User, error)
	DeleteFunction         func(orm.DB, *plethora_api.User) error
	UpdateFunction         func(orm.DB, *plethora_api.User) error
}

// Create mock
func (user *User) Create(db orm.DB, newUser plethora_api.User) (*plethora_api.User, error) {
	return user.CreateFunction(db, newUser)
}

// View mock
func (user *User) View(db orm.DB, id uuid.UUID) (*plethora_api.User, error) {
	return user.ViewFunction(db, id)
}

// FindByUsername mock
func (user *User) FindByUsername(db orm.DB, username string) (*plethora_api.User, error) {
	return user.FindByUsernameFunction(db, username)
}

// FindByToken mock
func (user *User) FindByToken(db orm.DB, token string) (*plethora_api.User, error) {
	return user.FindByTokenFunction(db, token)
}

// List mock
func (user *User) List(db orm.DB, listQuery *plethora_api.ListQuery, pagination *plethora_api.Pagination) ([]plethora_api.User, error) {
	return user.ListFunction(db, listQuery, pagination)
}

// Delete mock
func (user *User) Delete(db orm.DB, userToDelete *plethora_api.User) error {
	return user.DeleteFunction(db, userToDelete)
}

// Update mock
func (user *User) Update(db orm.DB, userToUpdate *plethora_api.User) error {
	return user.UpdateFunction(db, userToUpdate)
}
