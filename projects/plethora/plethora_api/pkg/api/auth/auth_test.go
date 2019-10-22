package auth_test

import (
	"testing"
	"time"

	"github.com/go-pg/pg/orm"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"

	auth "github.com/Wallruzz9114/plethora_api/pkg/api/auth"
	mock "github.com/Wallruzz9114/plethora_api/pkg/util/mock"
	mockdb "github.com/Wallruzz9114/plethora_api/pkg/util/mock/mockdb"
	plethora_api "github.com/Wallruzz9114/plethora_api/pkg/util/model"
	uuid "github.com/satori/go.uuid"
)

var (
	randomID = uuid.NewV4()
)

// TestAuthenticate ...
func TestAuthenticate(t *testing.T) {
	type args struct {
		user string
		pass string
	}
	cases := []struct {
		name     string
		args     args
		wantData *plethora_api.AuthToken
		wantErr  bool
		udb      *mockdb.User
		jwt      *mock.JWT
		sec      *mock.Secure
	}{
		{
			name:    "Fail on finding user",
			args:    args{user: "juzernejm"},
			wantErr: true,
			udb: &mockdb.User{
				FindByUsernameFunction: func(db orm.DB, user string) (*plethora_api.User, error) {
					return nil, plethora_api.ErrGeneric
				},
			},
		},
		{
			name:    "Fail on wrong password",
			args:    args{user: "juzernejm", pass: "notHashedPassword"},
			wantErr: true,
			udb: &mockdb.User{
				FindByUsernameFunction: func(db orm.DB, user string) (*plethora_api.User, error) {
					return &plethora_api.User{
						Username: user,
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
			name:    "Inactive user",
			args:    args{user: "juzernejm", pass: "pass"},
			wantErr: true,
			udb: &mockdb.User{
				FindByUsernameFunction: func(db orm.DB, user string) (*plethora_api.User, error) {
					return &plethora_api.User{
						Username: user,
						Password: "pass",
						Active:   false,
					}, nil
				},
			},
			sec: &mock.Secure{
				HashMatchesPasswordFunction: func(string, string) bool {
					return true
				},
			},
		},
		{
			name:    "Fail on token generation",
			args:    args{user: "juzernejm", pass: "pass"},
			wantErr: true,
			udb: &mockdb.User{
				FindByUsernameFunction: func(db orm.DB, user string) (*plethora_api.User, error) {
					return &plethora_api.User{
						Username: user,
						Password: "pass",
						Active:   true,
					}, nil
				},
			},
			sec: &mock.Secure{
				HashMatchesPasswordFunction: func(string, string) bool {
					return true
				},
			},
			jwt: &mock.JWT{
				GenerateTokenFunction: func(u *plethora_api.User) (string, string, error) {
					return "", "", plethora_api.ErrGeneric
				},
			},
		},
		{
			name:    "Fail on updating last login",
			args:    args{user: "juzernejm", pass: "pass"},
			wantErr: true,
			udb: &mockdb.User{
				FindByUsernameFunction: func(db orm.DB, user string) (*plethora_api.User, error) {
					return &plethora_api.User{
						Username: user,
						Password: "pass",
						Active:   true,
					}, nil
				},
				UpdateFunction: func(db orm.DB, u *plethora_api.User) error {
					return plethora_api.ErrGeneric
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
			jwt: &mock.JWT{
				GenerateTokenFunction: func(u *plethora_api.User) (string, string, error) {
					return "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9", mock.TestTime(2000).Format(time.RFC3339), nil
				},
			},
		},
		{
			name: "Success",
			args: args{user: "juzernejm", pass: "pass"},
			udb: &mockdb.User{
				FindByUsernameFunction: func(db orm.DB, user string) (*plethora_api.User, error) {
					return &plethora_api.User{
						Username: user,
						Password: "password",
						Active:   true,
					}, nil
				},
				UpdateFunction: func(db orm.DB, u *plethora_api.User) error {
					return nil
				},
			},
			jwt: &mock.JWT{
				GenerateTokenFunction: func(u *plethora_api.User) (string, string, error) {
					return "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9", mock.TestTime(2000).Format(time.RFC3339), nil
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
			wantData: &plethora_api.AuthToken{
				Token:        "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
				Expires:      mock.TestTime(2000).Format(time.RFC3339),
				RefreshToken: "refreshtoken",
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			s := auth.New(nil, tt.udb, tt.jwt, tt.sec, nil)
			token, err := s.Authenticate(nil, tt.args.user, tt.args.pass)
			if tt.wantData != nil {
				tt.wantData.RefreshToken = token.RefreshToken
				assert.Equal(t, tt.wantData, token)
			}
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

// TestRefresh ...
func TestRefresh(t *testing.T) {
	type args struct {
		c     echo.Context
		token string
	}
	cases := []struct {
		name     string
		args     args
		wantData *plethora_api.RefreshToken
		wantErr  bool
		udb      *mockdb.User
		jwt      *mock.JWT
	}{
		{
			name:    "Fail on finding token",
			args:    args{token: "refreshtoken"},
			wantErr: true,
			udb: &mockdb.User{
				FindByTokenFunction: func(db orm.DB, token string) (*plethora_api.User, error) {
					return nil, plethora_api.ErrGeneric
				},
			},
		},
		{
			name:    "Fail on token generation",
			args:    args{token: "refreshtoken"},
			wantErr: true,
			udb: &mockdb.User{
				FindByTokenFunction: func(db orm.DB, token string) (*plethora_api.User, error) {
					return &plethora_api.User{
						Username: "username",
						Password: "password",
						Active:   true,
						Token:    token,
					}, nil
				},
			},
			jwt: &mock.JWT{
				GenerateTokenFunction: func(u *plethora_api.User) (string, string, error) {
					return "", "", plethora_api.ErrGeneric
				},
			},
		},
		{
			name: "Success",
			args: args{token: "refreshtoken"},
			udb: &mockdb.User{
				FindByTokenFunction: func(db orm.DB, token string) (*plethora_api.User, error) {
					return &plethora_api.User{
						Username: "username",
						Password: "password",
						Active:   true,
						Token:    token,
					}, nil
				},
			},
			jwt: &mock.JWT{
				GenerateTokenFunction: func(u *plethora_api.User) (string, string, error) {
					return "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9", mock.TestTime(2000).Format(time.RFC3339), nil
				},
			},
			wantData: &plethora_api.RefreshToken{
				Token:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
				Expires: mock.TestTime(2000).Format(time.RFC3339),
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			s := auth.New(nil, tt.udb, tt.jwt, nil, nil)
			token, err := s.Refresh(tt.args.c, tt.args.token)
			assert.Equal(t, tt.wantData, token)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

// TestMe ...
func TestMe(t *testing.T) {
	cases := []struct {
		name     string
		wantData *plethora_api.User
		udb      *mockdb.User
		rbac     *mock.RBAC
		wantErr  bool
	}{
		{
			name: "Success",
			rbac: &mock.RBAC{
				UserFunction: func(echo.Context) *plethora_api.AuthUser {
					return &plethora_api.AuthUser{ID: randomID}
				},
			},
			udb: &mockdb.User{
				ViewFunction: func(db orm.DB, id uuid.UUID) (*plethora_api.User, error) {
					return &plethora_api.User{
						Base: plethora_api.Base{
							ID:        id,
							CreatedAt: mock.TestTime(1999),
							UpdatedAt: mock.TestTime(2000),
						},
						FirstName: "John",
						LastName:  "Doe",
						Role: &plethora_api.Role{
							AccessLevel: plethora_api.UserRole,
						},
					}, nil
				},
			},
			wantData: &plethora_api.User{
				Base: plethora_api.Base{
					ID:        randomID,
					CreatedAt: mock.TestTime(1999),
					UpdatedAt: mock.TestTime(2000),
				},
				FirstName: "John",
				LastName:  "Doe",
				Role: &plethora_api.Role{
					AccessLevel: plethora_api.UserRole,
				},
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			s := auth.New(nil, tt.udb, nil, nil, tt.rbac)
			user, err := s.Me(nil)
			assert.Equal(t, tt.wantData, user)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

// TestInitialize ...
func TestInitialize(t *testing.T) {
	a := auth.Initialize(nil, nil, nil, nil)
	if a == nil {
		t.Error("auth service not initialized")
	}
}
