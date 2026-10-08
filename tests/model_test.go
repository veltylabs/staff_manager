package tests

import (
	"testing"

	staffmanager "github.com/veltylabs/staff_manager"
)

func TestDomainError(t *testing.T) {
	if staffmanager.ErrNotFound.Error() != "staff not found" {
		t.Errorf("expected ErrNotFound.Error() to be 'staff not found', got %q", staffmanager.ErrNotFound.Error())
	}
}
