package query

import (
	"github.com/labstack/echo"

	plethora_api "github.com/Wallruzz9114/plethora_api/pkg/util/model"
)

// List prepares data for list queries
func List(u *plethora_api.AuthUser) (*plethora_api.ListQuery, error) {
	switch true {
	case u.Role <= plethora_api.AdminRole: // user is SuperAdmin or Admin
		return nil, nil
	case u.Role == plethora_api.CompanyAdminRole:
		return &plethora_api.ListQuery{Query: "company_id = ?", ID: u.CompanyID}, nil
	case u.Role == plethora_api.LocationAdminRole:
		return &plethora_api.ListQuery{Query: "location_id = ?", ID: u.LocationID}, nil
	default:
		return nil, echo.ErrForbidden
	}
}
