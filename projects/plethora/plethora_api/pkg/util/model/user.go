package plethora_api

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// User represents user domain model
type User struct {
	Base
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Password  string `json:"-"`
	Email     string `json:"email"`

	Mobile  string `json:"mobile,omitempty"`
	Phone   string `json:"phone,omitempty"`
	Address string `json:"address,omitempty"`

	Active bool `json:"active"`

	LastLogin          time.Time `json:"last_login,omitempty"`
	LastPasswordChange time.Time `json:"last_password_change,omitempty"`

	Token string `json:"-"`

	Role *Role `json:"role,omitempty"`

	RoleID     AccessRole `json:"-"`
	CompanyID  uuid.UUID  `json:"company_id"`
	LocationID uuid.UUID  `json:"location_id"`
}

// AuthUser represents data stored in JWT token for user
type AuthUser struct {
	ID         uuid.UUID
	CompanyID  uuid.UUID
	LocationID uuid.UUID
	Username   string
	Email      string
	Role       AccessRole
}

// ChangePasseword ...
func (user *User) ChangePassword(hashedPassword string) {
	user.Password = hashedPassword
	user.LastPasswordChange = time.Now()
}

// UpdateLastLogin updates last login field
func (user *User) UpdateLastLogin(token string) {
	user.Token = token
	user.LastLogin = time.Now()
}
