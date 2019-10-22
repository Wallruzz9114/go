package password

import (
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/labstack/echo"

	pgsql "github.com/Wallruzz9114/plethora_api/pkg/api/auth/platform/pgsql"
	plethora_api "github.com/Wallruzz9114/plethora_api/pkg/util/model"
	uuid "github.com/satori/go.uuid"
)

// Service represents password application interface
type Service interface {
	Change(echo.Context, uuid.UUID, string, string) error
}

// Password represents password application service
type Password struct {
	db   *pg.DB
	udb  UserDB
	rbac RBAC
	sec  Securer
}

// UserDB represents user repository interface
type UserDB interface {
	View(orm.DB, uuid.UUID) (*plethora_api.User, error)
	Update(orm.DB, *plethora_api.User) error
}

// Securer represents security interface
type Securer interface {
	HashPassword(string) string
	HashMatchesPassword(string, string) bool
	Password(string, ...string) bool
}

// RBAC represents role-based-access-control interface
type RBAC interface {
	EnforceUser(echo.Context, uuid.UUID) error
}

// New creates new password application service
func New(db *pg.DB, udb UserDB, rbac RBAC, sec Securer) *Password {
	return &Password{
		db:   db,
		udb:  udb,
		rbac: rbac,
		sec:  sec,
	}
}

// Initialize initalizes password application service with defaults
func Initialize(db *pg.DB, rbac RBAC, sec Securer) *Password {
	return New(db, pgsql.NewUser(), rbac, sec)
}
