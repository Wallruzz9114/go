package query_test

import (
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"

	plethora_api "github.com/Wallruzz9114/plethora_api/pkg/util/model"
	query "github.com/Wallruzz9114/plethora_api/pkg/util/query"
	uuid "github.com/satori/go.uuid"
)

var (
	randomID         = uuid.NewV4()
	randomLocationID = uuid.NewV4()
	randomCompanyID  = uuid.NewV4()
)

// TestList ...
func TestList(t *testing.T) {
	type args struct {
		user *plethora_api.AuthUser
	}
	cases := []struct {
		name     string
		args     args
		wantData *plethora_api.ListQuery
		wantErr  error
	}{
		{
			name: "Super admin user",
			args: args{user: &plethora_api.AuthUser{
				Role: plethora_api.SuperAdminRole,
			}},
		},
		{
			name: "Company admin user",
			args: args{user: &plethora_api.AuthUser{
				Role:      plethora_api.CompanyAdminRole,
				CompanyID: randomCompanyID,
			}},
			wantData: &plethora_api.ListQuery{
				Query: "company_id = ?",
				ID:    randomCompanyID,
			},
		},
		{
			name: "Location admin user",
			args: args{user: &plethora_api.AuthUser{
				Role:       plethora_api.LocationAdminRole,
				CompanyID:  randomCompanyID,
				LocationID: randomLocationID,
			}},
			wantData: &plethora_api.ListQuery{
				Query: "location_id = ?",
				ID:    randomLocationID,
			},
		},
		{
			name: "Normal user",
			args: args{user: &plethora_api.AuthUser{
				Role: plethora_api.UserRole,
			}},
			wantErr: echo.ErrForbidden,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			q, err := query.List(tt.args.user)
			assert.Equal(t, tt.wantData, q)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
