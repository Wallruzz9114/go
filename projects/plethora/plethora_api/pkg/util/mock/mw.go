package mock

import (
	plethora_api "github.com/Wallruzz9114/plethora_api/pkg/util/model"
)

// JWT mock
type JWT struct {
	GenerateTokenFunction func(*plethora_api.User) (string, string, error)
}

// GenerateToken mock
func (j *JWT) GenerateToken(user *plethora_api.User) (string, string, error) {
	return j.GenerateTokenFunction(user)
}
