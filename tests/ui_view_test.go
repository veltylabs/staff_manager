//go:build browser

package tests

import (
	"testing"

	staffmanager "github.com/veltylabs/staff_manager"
	"github.com/veltylabs/staff_manager/ui"
	"webtyp.com/model"
	"webtyp.com/router"
)

type mockCaller struct{}

func (c *mockCaller) Call(op router.Operation, args model.Encodable, res model.Decodable, cb func(error)) {
	if cb != nil {
		cb(nil)
	}
}

func (c *mockCaller) Dispatch(op router.Operation, args model.Encodable) {}

var _ router.Caller = (*mockCaller)(nil)

func TestUI_Screens(t *testing.T) {
	fake := &mockCaller{}

	screens := ui.Screens(fake)
	if len(screens) != 1 {
		t.Fatalf("len(screens) = %d, want 1", len(screens))
	}

	screen := screens[0]
	if err := screen.Validate(); err != nil {
		t.Fatalf("screen.Validate() error: %v", err)
	}

	if screen.ID != ui.ID {
		t.Errorf("ID = %q, want %q", screen.ID, ui.ID)
	}
	// La pantalla pide el mismo permiso que sus operaciones: la tabla staff_member.
	if want := model.ResourceOf(&staffmanager.StaffMemberModel); screen.Permission != want {
		t.Errorf("Permission = %q, want %q", screen.Permission, want)
	}

	if screen.Presenter.Description() == "" {
		t.Error("Description() is empty")
	}
}

func TestUI_ViewReloadReadsStaff(t *testing.T) {
	var called router.Operation
	caller := &customCaller{
		callFn: func(op router.Operation, args model.Encodable, res model.Decodable, cb func(error)) {
			called = op
			if cb != nil {
				cb(nil)
			}
		},
	}

	v := staffmanager.NewView(caller)
	v.Reload(nil)

	want := router.Operation{Module: staffmanager.ModelName, Entity: (&staffmanager.StaffMember{}).ModelName(), Action: model.Read}
	if called != want {
		t.Errorf("Reload called %+v, want %+v", called, want)
	}
}

type customCaller struct {
	callFn func(op router.Operation, args model.Encodable, res model.Decodable, cb func(error))
}

func (c *customCaller) Call(op router.Operation, args model.Encodable, res model.Decodable, cb func(error)) {
	c.callFn(op, args, res, cb)
}

func (c *customCaller) Dispatch(op router.Operation, args model.Encodable) {}
