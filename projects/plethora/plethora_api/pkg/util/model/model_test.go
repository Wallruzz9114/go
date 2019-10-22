package plethora_api_test

import (
	"testing"

	mock "github.com/Wallruzz9114/plethora_api/pkg/util/mock"
	plethora_api "github.com/Wallruzz9114/plethora_api/pkg/util/model"
	uuid "github.com/satori/go.uuid"
)

var (
	randomID = uuid.NewV4()
)

// TestBeforeInsert ...
func TestBeforeInsert(t *testing.T) {
	base := &plethora_api.Base{
		ID: randomID,
	}

	base.BeforeInsert(nil, nil)

	if base.CreatedAt.IsZero() {
		t.Error("CreatedAt was not changed")
	}

	if base.UpdatedAt.IsZero() {
		t.Error("UpdatedAt was not changed")
	}
}

// TestBeforeUpdate ...
func TestBeforeUpdate(t *testing.T) {
	base := &plethora_api.Base{
		ID:        randomID,
		CreatedAt: mock.TestTime(2000),
	}

	base.BeforeUpdate(nil, nil)

	if base.UpdatedAt == mock.TestTime(2001) {
		t.Error("UpdatedAt was not changed")
	}

}

// TestPaginationTransform ...
func TestPaginationTransform(t *testing.T) {
	p := &plethora_api.PaginationReq{
		Limit: 5000, Page: 5,
	}

	pag := p.Transform()

	if pag.Limit != 1000 {
		t.Error("Default limit not set")
	}

	if pag.Offset != 5000 {
		t.Error("Offset not set correctly")
	}

	p.Limit = 0
	newPag := p.Transform()

	if newPag.Limit != 100 {
		t.Error("Min limit not set")
	}

}
