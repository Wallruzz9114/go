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
	var user = new(plethora_api.User)

	sql := `SELECT "user".*, "role"."id" AS "role__id", "role"."access_level" AS "role__access_level", "role"."name" AS "role__name" 
	FROM "users" AS "user" LEFT JOIN "roles" AS "role" ON "role"."id" = "user"."role_id" 
	WHERE ("user"."id" = ? and deleted_at is null)`

	_, err := db.QueryOne(user, sql, id)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// FindByUsername queries for single user by username
func (u *User) FindByUsername(db orm.DB, uname string) (*plethora_api.User, error) {
	var user = new(plethora_api.User)

	sql := `SELECT "user".*, "role"."id" AS "role__id", "role"."access_level" AS "role__access_level", "role"."name" AS "role__name" 
	FROM "users" AS "user" LEFT JOIN "roles" AS "role" ON "role"."id" = "user"."role_id" 
	WHERE ("user"."username" = ? and deleted_at is null)`

	_, err := db.QueryOne(user, sql, uname)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// FindByToken queries for single user by token
func (u *User) FindByToken(db orm.DB, token string) (*plethora_api.User, error) {
	var user = new(plethora_api.User)

	sql := `SELECT "user".*, "role"."id" AS "role__id", "role"."access_level" AS "role__access_level", "role"."name" AS "role__name" 
	FROM "users" AS "user" LEFT JOIN "roles" AS "role" ON "role"."id" = "user"."role_id" 
	WHERE ("user"."token" = ? and deleted_at is null)`

	_, err := db.QueryOne(user, sql, token)

	if err != nil {
	}

	return user, err
}

// Update updates user's info
func (u *User) Update(db orm.DB, user *plethora_api.User) error {
	return db.Update(user)
}
