package mock

// Secure mock
type Secure struct {
	PasswordFunction            func(string, ...string) bool
	HashPasswordFunction        func(string) string
	HashMatchesPasswordFunction func(string, string) bool
	TokenFunction               func(string) string
}

// Password mock
func (s *Secure) Password(pw string, inputs ...string) bool {
	return s.PasswordFunction(pw, inputs...)
}

// HashPassword mock
func (s *Secure) HashPassword(pw string) string {
	return s.HashPasswordFunction(pw)
}

// HashMatchesPassword mock
func (s *Secure) HashMatchesPassword(hash, pw string) bool {
	return s.HashMatchesPasswordFunction(hash, pw)
}

// Token mock
func (s *Secure) Token(token string) string {
	return s.TokenFunction(token)
}
