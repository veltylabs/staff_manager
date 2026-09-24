package tests

import (
	"testing"

	"github.com/veltylabs/staff_manager/seed"
)

func TestSeed_Load(t *testing.T) {
	m := setup(t, nil)

	data, err := seed.Load(m, testTenant)
	if err != nil {
		t.Fatalf("seed.Load: %v", err)
	}

	if len(data.Staff) != 3 {
		t.Fatalf("len(data.Staff) = %d, want 3", len(data.Staff))
	}

	for _, member := range data.Staff {
		if member.Id == "" {
			t.Errorf("Staff member %s has empty Id", member.Name)
		}
		if member.TenantId != testTenant {
			t.Errorf("Staff member %s has TenantId = %q, want %q", member.Name, member.TenantId, testTenant)
		}
		if member.Specialty == "" {
			t.Errorf("Staff member %s has empty Specialty", member.Name)
		}
	}
}
