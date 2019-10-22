package rbac_test

import (
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"

	mock "github.com/Wallruzz9114/plethora_api/pkg/util/mock"
	plethora_api "github.com/Wallruzz9114/plethora_api/pkg/util/model"
	rbac "github.com/Wallruzz9114/plethora_api/pkg/util/rbac"
	uuid "github.com/satori/go.uuid"
)

var (
	randomID         = uuid.NewV4()
	randomLocationID = uuid.NewV4()
	randomCompanyID  = uuid.NewV4()
)

// TestUser ...
func TestUser(t *testing.T) {
	ctx := mock.EchoContextWithKeys(
		[]string{"id", "company_id", "location_id", "username", "email", "role"},
		randomID, randomCompanyID, randomLocationID, "jpinto", "pintojose.benedicto2@gmail.com", plethora_api.SuperAdminRole,
	)
	wantUser := &plethora_api.AuthUser{
		ID:         randomID,
		Username:   "ribice",
		CompanyID:  randomCompanyID,
		LocationID: randomLocationID,
		Email:      "ribice@gmail.com",
		Role:       plethora_api.SuperAdminRole,
	}
	rbacSvc := rbac.New()

	assert.Equal(t, wantUser, rbacSvc.User(ctx))
}

// TestEnforceRole ...
func TestEnforceRole(t *testing.T) {
	type args struct {
		ctx  echo.Context
		role plethora_api.AccessRole
	}
	cases := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Not authorized",
			args:    args{ctx: mock.EchoContextWithKeys([]string{"role"}, plethora_api.CompanyAdminRole), role: plethora_api.SuperAdminRole},
			wantErr: true,
		},
		{
			name:    "Authorized",
			args:    args{ctx: mock.EchoContextWithKeys([]string{"role"}, plethora_api.SuperAdminRole), role: plethora_api.CompanyAdminRole},
			wantErr: false,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			rbacSvc := rbac.New()
			res := rbacSvc.EnforceRole(tt.args.ctx, tt.args.role)
			assert.Equal(t, tt.wantErr, res == echo.ErrForbidden)
		})
	}
}

// TestEnforceUser ...
func TestEnforceUser(t *testing.T) {
	type args struct {
		ctx echo.Context
		id  uuid.UUID
	}
	cases := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Not same user, not an admin",
			args:    args{ctx: mock.EchoContextWithKeys([]string{"id", "role"}, 15, plethora_api.LocationAdminRole), id: randomID},
			wantErr: true,
		},
		{
			name:    "Not same user, but admin",
			args:    args{ctx: mock.EchoContextWithKeys([]string{"id", "role"}, 22, plethora_api.SuperAdminRole), id: randomID},
			wantErr: false,
		},
		{
			name:    "Same user",
			args:    args{ctx: mock.EchoContextWithKeys([]string{"id", "role"}, 8, plethora_api.AdminRole), id: randomID},
			wantErr: false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			rbacSvc := rbac.New()
			res := rbacSvc.EnforceUser(tt.args.ctx, tt.args.id)
			assert.Equal(t, tt.wantErr, res == echo.ErrForbidden)
		})
	}
}

// TestEnforceCompany ...
func TestEnforceCompany(t *testing.T) {
	type args struct {
		ctx echo.Context
		id  uuid.UUID
	}
	cases := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Not same company, not an admin",
			args:    args{ctx: mock.EchoContextWithKeys([]string{"company_id", "role"}, 7, plethora_api.UserRole), id: randomCompanyID},
			wantErr: true,
		},
		{
			name:    "Same company, not company admin or admin",
			args:    args{ctx: mock.EchoContextWithKeys([]string{"company_id", "role"}, 22, plethora_api.UserRole), id: randomCompanyID},
			wantErr: true,
		},
		{
			name:    "Same company, company admin",
			args:    args{ctx: mock.EchoContextWithKeys([]string{"company_id", "role"}, 5, plethora_api.CompanyAdminRole), id: randomCompanyID},
			wantErr: false,
		},
		{
			name:    "Not same company but admin",
			args:    args{ctx: mock.EchoContextWithKeys([]string{"company_id", "role"}, 8, plethora_api.AdminRole), id: randomCompanyID},
			wantErr: false,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			rbacSvc := rbac.New()
			res := rbacSvc.EnforceCompany(tt.args.ctx, tt.args.id)
			assert.Equal(t, tt.wantErr, res == echo.ErrForbidden)
		})
	}
}

// TestEnforceLocation ...
func TestEnforceLocation(t *testing.T) {
	type args struct {
		ctx echo.Context
		id  uuid.UUID
	}
	cases := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Not same location, not an admin",
			args:    args{ctx: mock.EchoContextWithKeys([]string{"location_id", "role"}, 7, plethora_api.UserRole), id: randomLocationID},
			wantErr: true,
		},
		{
			name:    "Same location, not company admin or admin",
			args:    args{ctx: mock.EchoContextWithKeys([]string{"location_id", "role"}, 22, plethora_api.UserRole), id: randomLocationID},
			wantErr: true,
		},
		{
			name:    "Same location, company admin",
			args:    args{ctx: mock.EchoContextWithKeys([]string{"location_id", "role"}, 5, plethora_api.CompanyAdminRole), id: randomLocationID},
			wantErr: false,
		},
		{
			name:    "Location admin",
			args:    args{ctx: mock.EchoContextWithKeys([]string{"location_id", "role"}, 5, plethora_api.LocationAdminRole), id: randomLocationID},
			wantErr: false,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			rbacSvc := rbac.New()
			res := rbacSvc.EnforceLocation(tt.args.ctx, tt.args.id)
			assert.Equal(t, tt.wantErr, res == echo.ErrForbidden)
		})
	}
}

// TestAccountCreate ...
func TestAccountCreate(t *testing.T) {
	type args struct {
		ctx         echo.Context
		roleID      plethora_api.AccessRole
		company_id  uuid.UUID
		location_id uuid.UUID
	}
	cases := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Different location, company, creating user role, not an admin",
			args:    args{ctx: mock.EchoContextWithKeys([]string{"company_id", "location_id", "role"}, 2, 3, plethora_api.UserRole), roleID: 500, company_id: randomCompanyID, location_id: randomLocationID},
			wantErr: true,
		},
		{
			name:    "Same location, not company, creating user role, not an admin",
			args:    args{ctx: mock.EchoContextWithKeys([]string{"company_id", "location_id", "role"}, 2, 3, plethora_api.UserRole), roleID: 500, company_id: randomCompanyID, location_id: randomLocationID},
			wantErr: true,
		},
		{
			name:    "Different location, company, creating user role, not an admin",
			args:    args{ctx: mock.EchoContextWithKeys([]string{"company_id", "location_id", "role"}, 2, 3, plethora_api.CompanyAdminRole), roleID: 400, company_id: randomCompanyID, location_id: randomLocationID},
			wantErr: false,
		},
		{
			name:    "Same location, company, creating user role, not an admin",
			args:    args{ctx: mock.EchoContextWithKeys([]string{"company_id", "location_id", "role"}, 2, 3, plethora_api.CompanyAdminRole), roleID: 500, company_id: randomCompanyID, location_id: randomLocationID},
			wantErr: false,
		},
		{
			name:    "Same location, company, creating user role, admin",
			args:    args{ctx: mock.EchoContextWithKeys([]string{"company_id", "location_id", "role"}, 2, 3, plethora_api.CompanyAdminRole), roleID: 500, company_id: randomCompanyID, location_id: randomLocationID},
			wantErr: false,
		},
		{
			name:    "Different everything, admin",
			args:    args{ctx: mock.EchoContextWithKeys([]string{"company_id", "location_id", "role"}, 2, 3, plethora_api.AdminRole), roleID: 200, company_id: randomCompanyID, location_id: randomLocationID},
			wantErr: false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			rbacSvc := rbac.New()
			res := rbacSvc.AccountCreate(tt.args.ctx, tt.args.roleID, tt.args.company_id, tt.args.location_id)
			assert.Equal(t, tt.wantErr, res == echo.ErrForbidden)
		})
	}
}

// TestIsLowerRole ...
func TestIsLowerRole(t *testing.T) {
	ctx := mock.EchoContextWithKeys([]string{"role"}, plethora_api.CompanyAdminRole)
	rbacSvc := rbac.New()

	if rbacSvc.IsLowerRole(ctx, plethora_api.LocationAdminRole) != nil {
		t.Error("The requested user is higher role than the user requesting it")
	}

	if rbacSvc.IsLowerRole(ctx, plethora_api.AdminRole) == nil {
		t.Error("The requested user is lower role than the user requesting it")
	}
}
