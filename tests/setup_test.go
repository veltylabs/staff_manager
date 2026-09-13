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

// fakeDevices es el doble de prueba de DeviceReader: un mapa en memoria ip->deviceID,
// que representa la lectura entre módulos que staff_manager realiza a
// device_manager a través de la raíz de composición de la aplicación real.
type fakeDevices struct{ byIP map[string]string }

func (f *fakeDevices) FindByIP(ip string) (string, bool) {
	id, ok := f.byIP[ip]
	return id, ok
}

var _ staffmanager.DeviceReader = (*fakeDevices)(nil)

// identityValidateRUT es el doble de prueba para Deps.ValidateRUT: transmite cualquier
// valor no vacío sin cambios, rechazando únicamente el literal "bad" — el
// algoritmo de suma de verificación real vive en webtyp.com/auth/trusted_ip y se
// omite deliberadamente aquí (lista blanca del módulo).
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
