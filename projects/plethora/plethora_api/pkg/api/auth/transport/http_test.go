package transport_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-pg/pg/orm"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"

	auth "github.com/Wallruzz9114/plethora_api/pkg/api/auth"
	transport "github.com/Wallruzz9114/plethora_api/pkg/api/auth/transport"
	jwt "github.com/Wallruzz9114/plethora_api/pkg/util/middleware/jwt"
	mock "github.com/Wallruzz9114/plethora_api/pkg/util/mock"
	mockdb "github.com/Wallruzz9114/plethora_api/pkg/util/mock/mockdb"
	plethora_api "github.com/Wallruzz9114/plethora_api/pkg/util/model"
	server "github.com/Wallruzz9114/plethora_api/pkg/util/server"
	uuid "github.com/satori/go.uuid"
)

var (
	randomID         = uuid.NewV4()
	randomLocationID = uuid.NewV4()
	randomCompanyID  = uuid.NewV4()
)

// TestLogin ...
func TestLogin(t *testing.T) {
	cases := []struct {
		name       string
		req        string
		wantStatus int
		wantResp   *plethora_api.AuthToken
		udb        *mockdb.User
		jwt        *mock.JWT
		sec        *mock.Secure
	}{
		{
			name:       "Invalid request",
			req:        `{"username":"juzernejm"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Fail on FindByUsername",
			req:        `{"username":"juzernejm","password":"hunter123"}`,
			wantStatus: http.StatusInternalServerError,
			udb: &mockdb.User{
				FindByUsernameFunction: func(orm.DB, string) (*plethora_api.User, error) {
					return nil, plethora_api.ErrGeneric
				},
			},
		},
		{
			name:       "Success",
			req:        `{"username":"juzernejm","password":"hunter123"}`,
			wantStatus: http.StatusOK,
			udb: &mockdb.User{
				FindByUsernameFunction: func(orm.DB, string) (*plethora_api.User, error) {
					return &plethora_api.User{
						Password: "hunter123",
						Active:   true,
					}, nil
				},
				UpdateFunction: func(db orm.DB, u *plethora_api.User) error {
					return nil
				},
			},
			jwt: &mock.JWT{
				GenerateTokenFunction: func(*plethora_api.User) (string, string, error) {
					return "jwttokenstring", mock.TestTime(2018).Format(time.RFC3339), nil
				},
			},
			sec: &mock.Secure{
				HashMatchesPasswordFunction: func(string, string) bool {
					return true
				},
				TokenFunction: func(string) string {
					return "refreshtoken"
				},
			},
			wantResp: &plethora_api.AuthToken{Token: "jwttokenstring", Expires: mock.TestTime(2018).Format(time.RFC3339), RefreshToken: "refreshtoken"},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			r := server.New()
			transport.NewHTTP(auth.New(nil, tt.udb, tt.jwt, tt.sec, nil), r, nil)
			ts := httptest.NewServer(r)

			defer ts.Close()
			path := ts.URL + "/login"
			res, err := http.Post(path, "application/json", bytes.NewBufferString(tt.req))

			if err != nil {
				t.Fatal(err)
			}

			defer res.Body.Close()

			if tt.wantResp != nil {
				response := new(plethora_api.AuthToken)
				if err := json.NewDecoder(res.Body).Decode(response); err != nil {
					t.Fatal(err)
				}
				tt.wantResp.RefreshToken = response.RefreshToken
				assert.Equal(t, tt.wantResp, response)
			}

			assert.Equal(t, tt.wantStatus, res.StatusCode)
		})
	}
}

// TestRefresh ...
func TestRefresh(t *testing.T) {
	cases := []struct {
		name       string
		req        string
		wantStatus int
		wantResp   *plethora_api.RefreshToken
		udb        *mockdb.User
		jwt        *mock.JWT
	}{
		{
			name:       "Fail on FindByToken",
			req:        "refreshtoken",
			wantStatus: http.StatusInternalServerError,
			udb: &mockdb.User{
				FindByTokenFunction: func(orm.DB, string) (*plethora_api.User, error) {
					return nil, plethora_api.ErrGeneric
				},
			},
		},
		{
			name:       "Success",
			req:        "refreshtoken",
			wantStatus: http.StatusOK,
			udb: &mockdb.User{
				FindByTokenFunction: func(orm.DB, string) (*plethora_api.User, error) {
					return &plethora_api.User{
						Username: "johndoe",
						Active:   true,
					}, nil
				},
			},
			jwt: &mock.JWT{
				GenerateTokenFunction: func(*plethora_api.User) (string, string, error) {
					return "jwttokenstring", mock.TestTime(2018).Format(time.RFC3339), nil
				},
			},
			wantResp: &plethora_api.RefreshToken{Token: "jwttokenstring", Expires: mock.TestTime(2018).Format(time.RFC3339)},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			r := server.New()
			transport.NewHTTP(auth.New(nil, tt.udb, tt.jwt, nil, nil), r, nil)

			ts := httptest.NewServer(r)

			defer ts.Close()

			path := ts.URL + "/refresh/" + tt.req
			res, err := http.Get(path)

			if err != nil {
				t.Fatal(err)
			}

			defer res.Body.Close()

			if tt.wantResp != nil {
				response := new(plethora_api.RefreshToken)
				if err := json.NewDecoder(res.Body).Decode(response); err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tt.wantResp, response)
			}

			assert.Equal(t, tt.wantStatus, res.StatusCode)
		})
	}
}

// TestMe ...
func TestMe(t *testing.T) {
	cases := []struct {
		name       string
		wantStatus int
		wantResp   *plethora_api.User
		header     string
		udb        *mockdb.User
		rbac       *mock.RBAC
	}{
		{
			name:       "Fail on user view",
			wantStatus: http.StatusInternalServerError,
			udb: &mockdb.User{
				ViewFunction: func(orm.DB, uuid.UUID) (*plethora_api.User, error) {
					return nil, plethora_api.ErrGeneric
				},
			},
			rbac: &mock.RBAC{
				UserFunction: func(echo.Context) *plethora_api.AuthUser {
					return &plethora_api.AuthUser{ID: randomID}
				},
			},
			header: mock.HeaderValid(),
		},
		{
			name:       "Success",
			wantStatus: http.StatusOK,
			udb: &mockdb.User{
				ViewFunction: func(db orm.DB, i uuid.UUID) (*plethora_api.User, error) {
					return &plethora_api.User{
						Base: plethora_api.Base{
							ID: i,
						},
						CompanyID:  randomCompanyID,
						LocationID: randomLocationID,
						Email:      "john@mail.com",
						FirstName:  "John",
						LastName:   "Doe",
					}, nil
				},
			},
			rbac: &mock.RBAC{
				UserFunction: func(echo.Context) *plethora_api.AuthUser {
					return &plethora_api.AuthUser{ID: randomID}
				},
			},
			header: mock.HeaderValid(),
			wantResp: &plethora_api.User{
				Base: plethora_api.Base{
					ID: randomID,
				},
				CompanyID:  randomCompanyID,
				LocationID: randomLocationID,
				Email:      "john@mail.com",
				FirstName:  "John",
				LastName:   "Doe",
			},
		},
	}

	client := &http.Client{}
	jwtMW := jwt.New("jwtsecret", "HS256", 60)

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			r := server.New()
			transport.NewHTTP(auth.New(nil, tt.udb, nil, nil, tt.rbac), r, jwtMW.MWFunc())
			ts := httptest.NewServer(r)
			defer ts.Close()
			path := ts.URL + "/me"
			req, err := http.NewRequest("GET", path, nil)
			req.Header.Set("Authorization", tt.header)

			if err != nil {
				t.Fatal(err)
			}

			res, err := client.Do(req)

			if err != nil {
				t.Fatal(err)
			}

			defer res.Body.Close()

			if tt.wantResp != nil {
				response := new(plethora_api.User)
				if err := json.NewDecoder(res.Body).Decode(response); err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tt.wantResp, response)
			}

			assert.Equal(t, tt.wantStatus, res.StatusCode)
		})
	}
}
