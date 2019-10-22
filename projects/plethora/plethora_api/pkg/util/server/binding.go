package server

import (
	"github.com/go-playground/validator"
	"github.com/labstack/echo"
)

// CustomBinder struct
type CustomBinder struct {
	b echo.Binder
}

// CustomValidator holds custom validator
type CustomValidator struct {
	V *validator.Validate
}

// NewBinder initializes custom server binder
func NewBinder() *CustomBinder {
	return &CustomBinder{b: &echo.DefaultBinder{}}
}

// Bind tries to bind request into interface, and if it does then validate it
func (cb *CustomBinder) Bind(i interface{}, c echo.Context) error {
	if err := cb.b.Bind(i, c); err != nil && err != echo.ErrUnsupportedMediaType {
		return err
	}
	return c.Validate(i)
}

// Validate validates the request
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.V.Struct(i)
}
