package mock

import (
	"github.com/labstack/echo"

	plethora_api "github.com/Wallruzz9114/plethora_api/pkg/util/model"
	uuid "github.com/satori/go.uuid"
)

// RBAC Mock
type RBAC struct {
	UserFunction            func(echo.Context) *plethora_api.AuthUser
	EnforceRoleFunction     func(echo.Context, plethora_api.AccessRole) error
	EnforceUserFunction     func(echo.Context, uuid.UUID) error
	EnforceCompanyFunction  func(echo.Context, uuid.UUID) error
	EnforceLocationFunction func(echo.Context, uuid.UUID) error
	AccountCreateFunction   func(echo.Context, plethora_api.AccessRole, uuid.UUID, uuid.UUID) error
	IsLowerRoleFunction     func(echo.Context, plethora_api.AccessRole) error
}

// User mock
func (a *RBAC) User(c echo.Context) *plethora_api.AuthUser {
	return a.UserFunction(c)
}

// EnforceRole mock
func (a *RBAC) EnforceRole(c echo.Context, role plethora_api.AccessRole) error {
	return a.EnforceRoleFunction(c, role)
}

// EnforceUser mock
func (a *RBAC) EnforceUser(c echo.Context, id uuid.UUID) error {
	return a.EnforceUserFunction(c, id)
}

// EnforceCompany mock
func (a *RBAC) EnforceCompany(c echo.Context, id uuid.UUID) error {
	return a.EnforceCompanyFunction(c, id)
}

// EnforceLocation mock
func (a *RBAC) EnforceLocation(c echo.Context, id uuid.UUID) error {
	return a.EnforceLocationFunction(c, id)
}

// AccountCreate mock
func (a *RBAC) AccountCreate(c echo.Context, roleID plethora_api.AccessRole, companyID, locationID uuid.UUID) error {
	return a.AccountCreateFunction(c, roleID, companyID, locationID)
}

// IsLowerRole mock
func (a *RBAC) IsLowerRole(c echo.Context, role plethora_api.AccessRole) error {
	return a.IsLowerRoleFunction(c, role)
}
