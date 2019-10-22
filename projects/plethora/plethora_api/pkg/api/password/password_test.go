package password_test

import (
	"testing"

	"github.com/go-pg/pg/orm"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"

	password "github.com/Wallruzz9114/plethora_api/pkg/api/password"
	mock "github.com/Wallruzz9114/plethora_api/pkg/util/mock"
	mockdb "github.com/Wallruzz9114/plethora_api/pkg/util/mock/mockdb"
	plethora_api "github.com/Wallruzz9114/plethora_api/pkg/util/model"
	uuid "github.com/satori/go.uuid"
)

var (
	randomID         = uuid.NewV4()
	randomLocationID = uuid.NewV4()
	randomCompanyID  = uuid.NewV4()
)

// TestChange ...
func TestChange(t *testing.T) {
	type args struct {
		oldpass string
		newpass string
		id      uuid.UUID
	}
	cases := []struct {
		name    string
		args    args
		wantErr bool
		udb     *mockdb.User
		rbac    *mock.RBAC
		sec     *mock.Secure
	}{
		{
			name: "Fail on EnforceUser",
			args: args{id: randomID},
			rbac: &mock.RBAC{
				EnforceUserFunction: func(c echo.Context, id uuid.UUID) error {
					return plethora_api.ErrGeneric
				}},
			wantErr: true,
		},
		{
			name:    "Fail on ViewUser",
			args:    args{id: randomID},
			wantErr: true,
			rbac: &mock.RBAC{
				EnforceUserFunction: func(c echo.Context, id uuid.UUID) error {
					return nil
				}},
			udb: &mockdb.User{
				ViewFunction: func(db orm.DB, id uuid.UUID) (*plethora_api.User, error) {
					if id != randomID {
						return nil, nil
					}
					return nil, plethora_api.ErrGeneric
				},
			},
		},
		{
			name: "Fail on PasswordMatch",
			args: args{id: randomID, oldpass: "hunter123"},
			rbac: &mock.RBAC{
				EnforceUserFunction: func(c echo.Context, id uuid.UUID) error {
					return nil
				}},
			wantErr: true,
			udb: &mockdb.User{
				ViewFunction: func(db orm.DB, id uuid.UUID) (*plethora_api.User, error) {
					return &plethora_api.User{
						Password: "HashedPassword",
					}, nil
				},
			},
			sec: &mock.Secure{
				HashMatchesPasswordFunction: func(string, string) bool {
					return false
				},
			},
		},
		{
			name: "Fail on InsecurePassword",
			args: args{id: randomID, oldpass: "hunter123"},
			rbac: &mock.RBAC{
				EnforceUserFunction: func(c echo.Context, id uuid.UUID) error {
					return nil
				}},
			wantErr: true,
			udb: &mockdb.User{
				ViewFunction: func(db orm.DB, id uuid.UUID) (*plethora_api.User, error) {
					return &plethora_api.User{
						Password: "HashedPassword",
					}, nil
				},
			},
			sec: &mock.Secure{
				HashMatchesPasswordFunction: func(string, string) bool {
					return true
				},
				PasswordFunction: func(string, ...string) bool {
					return false
				},
			},
		},
		{
			name: "Success",
			args: args{id: randomID, oldpass: "hunter123", newpass: "password"},
			rbac: &mock.RBAC{
				EnforceUserFunction: func(c echo.Context, id uuid.UUID) error {
					return nil
				}},
			udb: &mockdb.User{
				ViewFunction: func(db orm.DB, id uuid.UUID) (*plethora_api.User, error) {
					return &plethora_api.User{
						Password: "$2a$10$udRBroNGBeOYwSWCVzf6Lulg98uAoRCIi4t75VZg84xgw6EJbFunctionsG",
					}, nil
				},
				UpdateFunction: func(orm.DB, *plethora_api.User) error {
					return nil
				},
			},
			sec: &mock.Secure{
				HashMatchesPasswordFunction: func(string, string) bool {
					return true
				},
				PasswordFunction: func(string, ...string) bool {
					return true
				},
				HashPasswordFunction: func(string) string {
					return "hash3d"
				},
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			s := password.New(nil, tt.udb, tt.rbac, tt.sec)
			err := s.Change(nil, tt.args.id, tt.args.oldpass, tt.args.newpass)
			assert.Equal(t, tt.wantErr, err != nil)
			// Check whether password was changed
		})
	}
}

// TestInitialize ...
func TestInitialize(t *testing.T) {
	p := password.Initialize(nil, nil, nil)

	if p == nil {
		t.Error("password service not initialized")
	}
}
