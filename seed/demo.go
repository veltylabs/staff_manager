//go:build wasm

package seed

import (
	staffmanager "github.com/veltylabs/staff_manager"
	"github.com/veltylabs/staff_manager/ui"
	"webtyp.com/crud"
	"webtyp.com/input"
	"webtyp.com/orm"
	"webtyp.com/storage/mem"
	"webtyp.com/unixid"
	"webtyp.com/view"
)

const demoTenantID = "demo"

// demoDevices reemplaza a device_manager en esta demo: ningún equipo está registrado, así que
// ninguna IP inicia sesión. La app real conecta device_manager por Deps.Devices.
type demoDevices struct{}

func (demoDevices) FindByIP(string) (string, bool) { return "", false }

// Demo crea una demostración ejecutable del módulo staff_manager sobre mem db con datos semilla.
//
// Ejemplo: demo, err := seed.Demo()
func Demo() (view.Demo, error) {
	db := orm.New(mem.New())
	ids, err := unixid.NewUnixID()
	if err != nil {
		return view.Demo{}, err
	}
	mod, err := staffmanager.New(db, staffmanager.Deps{
		IDs:         ids,
		TenantID:    demoTenantID,
		ValidateRUT: input.NormalizeRUT,
		Devices:     demoDevices{},
	})
	if err != nil {
		return view.Demo{}, err
	}
	if _, err := Load(mod, demoTenantID); err != nil {
		return view.Demo{}, err
	}
	return view.Demo{
		Label:    ui.Label,
		Timezone: "America/Santiago",
		ActingAs: "u1",
		Modules:  []crud.Module{mod},
		Screens: func(d view.BrowserDeps) []view.Screen {
			return ui.Screens(d.Caller)
		},
	}, nil
}
