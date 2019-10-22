package server_test

import (
	"testing"

	server "github.com/Wallruzz9114/plethora_api/pkg/util/server"
)

// TestNew - Improve tests
func TestNew(t *testing.T) {
	e := server.New()
	if e == nil {
		t.Errorf("Server should not be nil")
	}
}