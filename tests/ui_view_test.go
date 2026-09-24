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

func TestUI_StaffPanelAndBrowser(t *testing.T) {
	fake := &mockCaller{}
	ids := &mockIDGen{}

	panel, err := ui.StaffPanel(fake, ids, "parent")
	if err != nil {
		t.Fatalf("StaffPanel error: %v", err)
	}
	if panel == nil {
		t.Fatal("StaffPanel returned nil component")
	}

	module, err := ui.Browser(fake, ids, testTenant)
	if err != nil {
		t.Fatalf("Browser error: %v", err)
	}
	if module == nil {
		t.Fatal("Browser returned nil UIModule")
	}

	if got := module.ModelName(); got != ui.ID {
		t.Errorf("ModelName = %q, want %q", got, ui.ID)
	}
	if got := module.Label(); got != ui.Label {
		t.Errorf("Label = %q, want %q", got, ui.Label)
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
