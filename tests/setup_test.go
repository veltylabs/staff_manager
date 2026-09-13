package tests

import (
	"testing"

	"webtyp.com/events"
	"webtyp.com/fmt"
	"webtyp.com/model"
	"webtyp.com/orm"
	"webtyp.com/storage/mem"
	staffmanager "github.com/veltylabs/staff_manager"
)

const testTenant = "test-tenant"

type mockIDGen struct{ counter int }

func (g *mockIDGen) NewID() string {
	g.counter++
	return "test-id-" + fmt.Convert(g.counter).String()
}

var _ model.IDGenerator = (*mockIDGen)(nil)

type mockPublisher struct{ Events []events.Event }

func (p *mockPublisher) Publish(e events.Event) { p.Events = append(p.Events, e) }

var _ events.Publisher = (*mockPublisher)(nil)

// fakeDevices is the DeviceReader test double: an in-memory ip->deviceID map,
// standing in for the cross-module read staff_manager makes into
// device_manager over the real app's composition root.
type fakeDevices struct{ byIP map[string]string }

func (f *fakeDevices) FindByIP(ip string) (string, bool) {
	id, ok := f.byIP[ip]
	return id, ok
}

var _ staffmanager.DeviceReader = (*fakeDevices)(nil)

// identityValidateRUT is the test double for Deps.ValidateRUT: passes any
// non-empty value through unchanged, rejecting only the literal "bad" — the
// real checksum algorithm lives in webtyp.com/auth/trusted_ip and is
// deliberately not imported here (module whitelist).
func identityValidateRUT(raw string) (string, error) {
	if raw == "bad" {
		return "", fmt.Err("staff_manager/tests: invalid rut")
	}
	return raw, nil
}

func setup(t *testing.T, devices *fakeDevices) *staffmanager.Module {
	t.Helper()
	db := orm.New(mem.New())
	if devices == nil {
		devices = &fakeDevices{byIP: map[string]string{}}
	}
	m, err := staffmanager.New(db, staffmanager.Deps{
		IDs:         &mockIDGen{},
		Publisher:   &mockPublisher{},
		TenantID:    testTenant,
		ValidateRUT: identityValidateRUT,
		Devices:     devices,
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return m
}
