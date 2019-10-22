package transport_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-pg/pg/orm"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"

	password "github.com/Wallruzz9114/plethora_api/pkg/api/password"
	transport "github.com/Wallruzz9114/plethora_api/pkg/api/password/transport"
	mock "github.com/Wallruzz9114/plethora_api/pkg/util/mock"
	mockdb "github.com/Wallruzz9114/plethora_api/pkg/util/mock/mockdb"
	plethora_api "github.com/Wallruzz9114/plethora_api/pkg/util/model"
	server "github.com/Wallruzz9114/plethora_api/pkg/util/server"
	uuid "github.com/satori/go.uuid"
)

// TestChangePassword
func TestChangePassword(t *testing.T) {
	cases := []struct {
		name       string
		req        string
		wantStatus int
		id         string
		udb        *mockdb.User
		rbac       *mock.RBAC
		sec        *mock.Secure
	}{
		{
			name:       "NaN",
			wantStatus: http.StatusBadRequest,
			id:         "abc",
		},
		{
			name:       "Fail on Bind",
			req:        `{"new_password":"new","old_password":"my_old_password", "new_password_confirm":"new"}`,
			wantStatus: http.StatusBadRequest,
			id:         "1",
		},
		{
			name:       "Different passwords",
			req:        `{"new_password":"new_password","old_password":"my_old_password", "new_password_confirm":"new_password_cf"}`,
			wantStatus: http.StatusBadRequest,
			id:         "1",
		},
		{
			name: "Fail on RBAC",
			req:  `{"new_password":"newpassw","old_password":"oldpassw", "new_password_confirm":"newpassw"}`,
			rbac: &mock.RBAC{
				EnforceUserFunction: func(c echo.Context, id uuid.UUID) error {
					return echo.ErrForbidden
				},
			},
			id:         "1",
			wantStatus: http.StatusForbidden,
		},
		{
			name: "Success",
			req:  `{"new_password":"newpassw","old_password":"oldpassw", "new_password_confirm":"newpassw"}`,
			rbac: &mock.RBAC{
				EnforceUserFunction: func(c echo.Context, id uuid.UUID) error {
					return nil
				},
			},
			id: "1",
			udb: &mockdb.User{
				ViewFunction: func(db orm.DB, id uuid.UUID) (*plethora_api.User, error) {
					return &plethora_api.User{
						Password: "oldPassword",
					}, nil
				},
				UpdateFunction: func(db orm.DB, usr *plethora_api.User) error {
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
					return "hashedPassword"
				},
			},
			wantStatus: http.StatusOK,
		},
	}

	client := &http.Client{}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			r := server.New()
			rg := r.Group("")
			transport.NewHTTP(password.New(nil, tt.udb, tt.rbac, tt.sec), rg)
			ts := httptest.NewServer(r)

			defer ts.Close()

			path := ts.URL + "/password/" + tt.id
			req, err := http.NewRequest("PATCH", path, bytes.NewBufferString(tt.req))
			req.Header.Set("Content-Type", "application/json")

			if err != nil {
				t.Fatal(err)
			}

			res, err := client.Do(req)

			if err != nil {
				t.Fatal(err)
			}

			defer res.Body.Close()

			assert.Equal(t, tt.wantStatus, res.StatusCode)
		})
	}
}
