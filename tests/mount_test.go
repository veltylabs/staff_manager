package tests

import (
	"testing"

	staffmanager "github.com/veltylabs/staff_manager"
	"webtyp.com/crud"
	"webtyp.com/model"
	"webtyp.com/orm"
	"webtyp.com/router"
	"webtyp.com/storage/mem"
)

func TestMount_MatchesDefinition(t *testing.T) {
	m := setup(t, nil)
	if err := crud.Mount(nil, nil, nil, m); err != nil {
		t.Fatalf("what MountOperations mounts must match Definition(), got %v", err)
	}
}

func TestCreate_PublishesChange(t *testing.T) {
	devices := &fakeDevices{byIP: map[string]string{}}
	pub := &mockPublisher{}
	db := orm.New(mem.New())
	m, err := staffmanager.New(db, staffmanager.Deps{
		IDs: &mockIDGen{}, Publisher: pub, TenantID: testTenant,
		ValidateRUT: identityValidateRUT, Devices: devices, Registrar: devices,
	})
	if err != nil {
		t.Fatal(err)
	}
	var callErr error
	crud.LoopbackActingAs(pub, "u1", m).Call(
		router.Operation{Module: staffmanager.ModelName, Entity: (&staffmanager.StaffMember{}).ModelName(), Action: model.Create},
		&crud.CreateArgs{Records: []model.Model{&staffmanager.StaffMember{Id: "s1", Rut: "11111111-1", Name: "Ana", IsActive: true}}},
		nil, func(e error) { callErr = e })
	if callErr != nil {
		t.Fatalf("create: %v", callErr)
	}
	topic := crud.Topic(staffmanager.ModelName, (&staffmanager.StaffMember{}).ModelName(), model.Create)
	count := 0
	for _, e := range pub.Events {
		if e.Topic != topic {
			continue
		}
		count++
		if _, ok := e.Payload.(*crud.Change); !ok {
			t.Fatalf("the created event must carry *crud.Change, got %T", e.Payload)
		}
	}
	if count != 1 {
		t.Fatalf("one created event on %q expected, got %d", topic, count)
	}
}
