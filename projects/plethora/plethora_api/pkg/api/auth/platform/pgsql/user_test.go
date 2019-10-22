package pgsql_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	pgsql "github.com/Wallruzz9114/plethora_api/pkg/api/auth/platform/pgsql"
	mock "github.com/Wallruzz9114/plethora_api/pkg/util/mock"
	plethora_api "github.com/Wallruzz9114/plethora_api/pkg/util/model"
	uuid "github.com/satori/go.uuid"
)

var (
	randomID         = uuid.NewV4()
	randomLocationID = uuid.NewV4()
	randomCompanyID  = uuid.NewV4()
)

// TestView ...
func TestView(t *testing.T) {
	cases := []struct {
		name     string
		wantErr  bool
		id       uuid.UUID
		wantData *plethora_api.User
	}{
		{
			name:    "User does not exist",
			wantErr: true,
			id:      randomID,
		},
		{
			name: "Success",
			id:   randomID,
			wantData: &plethora_api.User{
				Email:      "tomjones@mail.com",
				FirstName:  "Tom",
				LastName:   "Jones",
				Username:   "tomjones",
				RoleID:     1,
				CompanyID:  randomCompanyID,
				LocationID: randomLocationID,
				Password:   "newPass",
				Base: plethora_api.Base{
					ID: randomID,
				},
				Role: &plethora_api.Role{
					ID:          1,
					AccessLevel: 1,
					Name:        "SUPER_ADMIN",
				},
			},
		},
	}

	dbCon := mock.NewPGContainer(t)
	defer dbCon.Shutdown()

	db := mock.NewDB(t, dbCon, &plethora_api.Role{}, &plethora_api.User{})

	if err := mock.InsertMultiple(db, &plethora_api.Role{
		ID:          1,
		AccessLevel: 1,
		Name:        "SUPER_ADMIN"}, cases[1].wantData); err != nil {
		t.Error(err)
	}

	udb := pgsql.NewUser()

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			user, err := udb.View(db, tt.id)
			assert.Equal(t, tt.wantErr, err != nil)
			if tt.wantData != nil {
				if user == nil {
					t.Errorf("response was nil due to: %v", err)
				} else {
					tt.wantData.CreatedAt = user.CreatedAt
					tt.wantData.UpdatedAt = user.UpdatedAt
					assert.Equal(t, tt.wantData, user)
				}
			}
		})
	}
}

// TestFindByUsername ...
func TestFindByUsername(t *testing.T) {
	cases := []struct {
		name     string
		wantErr  bool
		username string
		wantData *plethora_api.User
	}{
		{
			name:     "User does not exist",
			wantErr:  true,
			username: "notExists",
		},
		{
			name:     "Success",
			username: "tomjones",
			wantData: &plethora_api.User{
				Email:      "tomjones@mail.com",
				FirstName:  "Tom",
				LastName:   "Jones",
				Username:   "tomjones",
				RoleID:     1,
				CompanyID:  randomCompanyID,
				LocationID: randomLocationID,
				Password:   "newPass",
				Base: plethora_api.Base{
					ID: randomID,
				},
				Role: &plethora_api.Role{
					ID:          1,
					AccessLevel: 1,
					Name:        "SUPER_ADMIN",
				},
			},
		},
	}

	dbCon := mock.NewPGContainer(t)
	defer dbCon.Shutdown()

	db := mock.NewDB(t, dbCon, &plethora_api.Role{}, &plethora_api.User{})

	if err := mock.InsertMultiple(db, &plethora_api.Role{
		ID:          1,
		AccessLevel: 1,
		Name:        "SUPER_ADMIN"}, cases[1].wantData); err != nil {
		t.Error(err)
	}

	udb := pgsql.NewUser()

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			user, err := udb.FindByUsername(db, tt.username)
			assert.Equal(t, tt.wantErr, err != nil)

			if tt.wantData != nil {
				tt.wantData.CreatedAt = user.CreatedAt
				tt.wantData.UpdatedAt = user.UpdatedAt
				assert.Equal(t, tt.wantData, user)

			}
		})
	}
}

// TestFindByToken ...
func TestFindByToken(t *testing.T) {
	cases := []struct {
		name     string
		wantErr  bool
		token    string
		wantData *plethora_api.User
	}{
		{
			name:    "User does not exist",
			wantErr: true,
			token:   "notExists",
		},
		{
			name:  "Success",
			token: "loginrefresh",
			wantData: &plethora_api.User{
				Email:      "johndoe@mail.com",
				FirstName:  "John",
				LastName:   "Doe",
				Username:   "johndoe",
				RoleID:     1,
				CompanyID:  randomCompanyID,
				LocationID: randomLocationID,
				Password:   "hunter2",
				Base: plethora_api.Base{
					ID: randomID,
				},
				Role: &plethora_api.Role{
					ID:          1,
					AccessLevel: 1,
					Name:        "SUPER_ADMIN",
				},
				Token: "loginrefresh",
			},
		},
	}

	dbCon := mock.NewPGContainer(t)
	defer dbCon.Shutdown()

	db := mock.NewDB(t, dbCon, &plethora_api.Role{}, &plethora_api.User{})

	if err := mock.InsertMultiple(db, &plethora_api.Role{
		ID:          1,
		AccessLevel: 1,
		Name:        "SUPER_ADMIN"}, cases[1].wantData); err != nil {
		t.Error(err)
	}

	udb := pgsql.NewUser()

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			user, err := udb.FindByToken(db, tt.token)
			assert.Equal(t, tt.wantErr, err != nil)

			if tt.wantData != nil {
				tt.wantData.CreatedAt = user.CreatedAt
				tt.wantData.UpdatedAt = user.UpdatedAt
				assert.Equal(t, tt.wantData, user)

			}
		})
	}
}

// TestUpdate ...
func TestUpdate(t *testing.T) {
	cases := []struct {
		name     string
		wantErr  bool
		usr      *plethora_api.User
		wantData *plethora_api.User
	}{
		{
			name: "Success",
			usr: &plethora_api.User{
				Base: plethora_api.Base{
					ID: randomID,
				},
				FirstName: "Z",
				LastName:  "Freak",
				Address:   "Address",
				Phone:     "123456",
				Mobile:    "345678",
				Username:  "newUsername",
			},
			wantData: &plethora_api.User{
				Email:      "tomjones@mail.com",
				FirstName:  "Z",
				LastName:   "Freak",
				Username:   "tomjones",
				RoleID:     1,
				CompanyID:  randomCompanyID,
				LocationID: randomLocationID,
				Password:   "newPass",
				Address:    "Address",
				Phone:      "123456",
				Mobile:     "345678",
				Base: plethora_api.Base{
					ID: randomID,
				},
			},
		},
	}

	dbCon := mock.NewPGContainer(t)
	defer dbCon.Shutdown()

	db := mock.NewDB(t, dbCon, &plethora_api.Role{}, &plethora_api.User{})

	if err := mock.InsertMultiple(db, &plethora_api.Role{
		ID:          1,
		AccessLevel: 1,
		Name:        "SUPER_ADMIN"}, cases[0].usr); err != nil {
		t.Error(err)
	}

	udb := pgsql.NewUser()

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			err := udb.Update(db, tt.wantData)
			assert.Equal(t, tt.wantErr, err != nil)
			if tt.wantData != nil {
				user := &plethora_api.User{
					Base: plethora_api.Base{
						ID: tt.usr.ID,
					},
				}
				if err := db.Select(user); err != nil {
					t.Error(err)
				}
				tt.wantData.UpdatedAt = user.UpdatedAt
				tt.wantData.CreatedAt = user.CreatedAt
				tt.wantData.LastLogin = user.LastLogin
				tt.wantData.DeletedAt = user.DeletedAt
				assert.Equal(t, tt.wantData, user)
			}
		})
	}
}
