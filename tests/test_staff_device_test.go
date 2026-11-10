package tests

import (
	"github.com/veltylabs/staff_manager"
	"strings"
	"testing"
	"webtyp.com/crud"
	"webtyp.com/model"
	"webtyp.com/router"
)

func TestStaffDevice_NotOffered(t *testing.T) {
	m := setup(t, nil)
	caller := crud.LoopbackActingAs(nil, "u1", m)
	caller.Call(router.Operation{Module: staffmanager.ModelName, Entity: "staff_device", Action: model.Read}, &crud.ReadArgs{}, nil, func(err error) {
		if err == nil || !strings.Contains(err.Error(), "unknown operation") {
			t.Errorf("expected unknown operation error, got %v", err)
		}
	})
}
