package auth

import (
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/labstack/echo"

	pgsql "github.com/Wallruzz9114/plethora_api/pkg/api/auth/platform/pgsql"
	plethora_api "github.com/Wallruzz9114/plethora_api/pkg/util/model"
	uuid "github.com/satori/go.uuid"
)

// Service represents auth service interface
type Service interface {
	Authenticate(echo.Context, string, string) (*plethora_api.AuthToken, error)
	Refresh(echo.Context, string) (*plethora_api.RefreshToken, error)
	Me(echo.Context) (*plethora_api.User, error)
}

// Auth represents auth application service
type Auth struct {
	db   *pg.DB
	udb  UserDB
	tg   TokenGenerator
	sec  Securer
	rbac RBAC
}

// UserDB represents user repository interface
type UserDB interface {
	View(orm.DB, uuid.UUID) (*plethora_api.User, error)
	FindByUsername(orm.DB, string) (*plethora_api.User, error)
	FindByToken(orm.DB, string) (*plethora_api.User, error)
	Update(orm.DB, *plethora_api.User) error
}

// TokenGenerator represents token generator (jwt) interface
type TokenGenerator interface {
	GenerateToken(*plethora_api.User) (string, string, error)
}

// Securer represents security interface
type Securer interface {
	HashMatchesPassword(string, string) bool
	Token(string) string
}

// RBAC represents role-based-access-control interface
type RBAC interface {
	User(echo.Context) *plethora_api.AuthUser
}

// New creates new iam service
func New(db *pg.DB, udb UserDB, j TokenGenerator, sec Securer, rbac RBAC) *Auth {
	return &Auth{
		db:   db,
		udb:  udb,
		tg:   j,
		sec:  sec,
		rbac: rbac,
	}
}

// Initialize initializes auth application service
func Initialize(db *pg.DB, j TokenGenerator, sec Securer, rbac RBAC) *Auth {
	return New(db, pgsql.NewUser(), j, sec, rbac)
}
