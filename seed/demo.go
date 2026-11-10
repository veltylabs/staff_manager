//go:build wasm
package seed
import (
	"webtyp.com/appkit"
	"webtyp.com/crud"
	"webtyp.com/events"
	"webtyp.com/fmt"
	"webtyp.com/model"
	"webtyp.com/orm"
	"webtyp.com/storage/mem"
	"webtyp.com/view"
	staffmanager "github.com/veltylabs/staff_manager"
	"github.com/veltylabs/staff_manager/ui"
)
type mockIDGen struct{ counter int }
func (g *mockIDGen) NewID() string { g.counter++; return "test-id-" + fmt.Convert(g.counter).String() }
type mockPublisher struct{ Events []events.Event }
func (p *mockPublisher) Publish(e events.Event) { p.Events = append(p.Events, e) }
func Demo() view.Demo {
	db := orm.New(mem.New())
	mod, _ := staffmanager.New(db, staffmanager.Deps{
		IDs: &mockIDGen{},
		TenantID: "test",
		Publisher: &mockPublisher{},
		ValidateRUT: func(raw string) (string, error) { return raw, nil },
	})
	Load(mod, "test")
	return view.Demo{
		Label:      "Staff",
		Timezone:   "America/Santiago",
		ActingAs:   "u1",
		Modules:    []crud.Module{mod},
		Screens: func(d appkit.BrowserDeps) []view.Screen { return ui.Screens(d.Caller) },
	}
}
