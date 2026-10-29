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

func (c *mockCaller) Call(route string, args model.Encodable, res model.Decodable, cb func(error)) {
	if cb != nil {
		cb(nil)
	}
}

func (c *mockCaller) Dispatch(route string, args model.Encodable) {}

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

func TestUI_ViewReloadCallsListStaff(t *testing.T) {
	var calledRoute string
	caller := &customCaller{
		callFn: func(route string, args model.Encodable, res model.Decodable, cb func(error)) {
			calledRoute = route
			if cb != nil {
				cb(nil)
			}
		},
	}

	v := staffmanager.NewView(caller)
	v.Reload(nil)

	if calledRoute != "staff_manager.list_staff" {
		t.Errorf("Reload called route %q, want %q", calledRoute, "staff_manager.list_staff")
	}
}

type customCaller struct {
	callFn func(route string, args model.Encodable, res model.Decodable, cb func(error))
}

func (c *customCaller) Call(route string, args model.Encodable, res model.Decodable, cb func(error)) {
	c.callFn(route, args, res, cb)
}

func (c *customCaller) Dispatch(route string, args model.Encodable) {}
