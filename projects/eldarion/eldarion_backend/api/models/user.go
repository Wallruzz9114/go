package models

import (
	"errors"
	"html"
	"log"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

// User ...
type User struct {
	ID        uuid.UUID `gorm:"primary_key;unique;type:uuid; column:id;default:gen_random_uuid()" json:"id"`
	Nickname  string    `gorm:"size:255;not null;unique" json:"nickname"`
	Email     string    `gorm:"size:100;not null;unique" json:"email"`
	Password  string    `gorm:"size:100;not null;" json:"password"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// HashPassword ...
func HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

// VerifyPassword ...
func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// BeforeSave ...
func (user *User) BeforeSave() error {
	hashedPassword, err := HashPassword(user.Password)

	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	return nil
}

// Prepare ...
func (user *User) Prepare() {
	user.ID = uuid.Nil
	user.Nickname = html.EscapeString(strings.TrimSpace(user.Nickname))
	user.Email = html.EscapeString(strings.TrimSpace(user.Email))
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
}

// Validate ...
func (user *User) Validate(action string) error {
	switch strings.ToLower(action) {
	case "update":
		if user.Nickname == "" {
			return errors.New("Required Nickname")
		}
		if user.Password == "" {
			return errors.New("Required Password")
		}
		if user.Email == "" {
			return errors.New("Required Email")
		}
		if err := checkmail.ValidateFormat(user.Email); err != nil {
			return errors.New("Invalid Email")
		}

		return nil
	case "login":
		if user.Password == "" {
			return errors.New("Required Password")
		}
		if user.Email == "" {
			return errors.New("Required Email")
		}
		if err := checkmail.ValidateFormat(user.Email); err != nil {
			return errors.New("Invalid Email")
		}
		return nil

	default:
		if user.Nickname == "" {
			return errors.New("Required Nickname")
		}
		if user.Password == "" {
			return errors.New("Required Password")
		}
		if user.Email == "" {
			return errors.New("Required Email")
		}
		if err := checkmail.ValidateFormat(user.Email); err != nil {
			return errors.New("Invalid Email")
		}
		return nil
	}
}

// SaveUser ...
func (user *User) SaveUser(db *gorm.DB) (*User, error) {
	var err error

	err = db.Debug().Create(&user).Error

	if err != nil {
		return &User{}, err
	}

	return user, nil
}

// FindAllUsers ...
func (user *User) FindAllUsers(db *gorm.DB) (*[]User, error) {
	var err error
	users := []User{}
	err = db.Debug().Model(&User{}).Limit(100).Find(&users).Error

	if err != nil {
		return &[]User{}, err
	}

	return &users, err
}

// FindUserByID ...
func (user *User) FindUserByID(db *gorm.DB, uid uuid.UUID) (*User, error) {
	var err error
	err = db.Debug().Model(User{}).Where("id = ?", uid).Take(&user).Error

	if err != nil {
		return &User{}, err
	}

	if gorm.IsRecordNotFoundError(err) {
		return &User{}, errors.New("User Not Found")
	}

	return user, err
}

// UpdateAUser
func (user *User) UpdateAUser(db *gorm.DB, uid uuid.UUID) (*User, error) {
	// To hash the password
	err := user.BeforeSave()

	if err != nil {
		log.Fatal(err)
	}

	db = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&User{}).UpdateColumns(
		map[string]interface{}{
			"password":  user.Password,
			"nickname":  user.Nickname,
			"email":     user.Email,
			"update_at": time.Now(),
		},
	)

	if db.Error != nil {
		return &User{}, db.Error
	}

	// This is the display the updated user
	err = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&user).Error

	if err != nil {
		return &User{}, err
	}

	return user, nil
}

// DeleteAUser ...
func (user *User) DeleteAUser(db *gorm.DB, uid uuid.UUID) (int64, error) {

	db = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&User{}).Delete(&User{})

	if db.Error != nil {
		return 0, db.Error
	}

	return db.RowsAffected, nil
}
