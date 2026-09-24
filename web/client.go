//go:build wasm

package main

import (
	devicemanager "github.com/veltylabs/device_manager"
	deviceseed "github.com/veltylabs/device_manager/seed"
	staffmanager "github.com/veltylabs/staff_manager"
	"github.com/veltylabs/staff_manager/seed"
	"github.com/veltylabs/staff_manager/ui"
	"webtyp.com/auth/trusted_ip"
	. "webtyp.com/dom"
	"webtyp.com/events/mock"
	"webtyp.com/layout/platformd"
	"webtyp.com/orm"
	"webtyp.com/router/loopback"
	"webtyp.com/storage/mem"
	"webtyp.com/unixid"
)

// demoTenantID es el único tenant de esta demo ejecutable en el navegador.
const demoTenantID = "demo"

// demoUser es la identidad fija que muestra el shell de la demo: no requiere inicio de sesión.
type demoUser struct{}

func (demoUser) UserName() string    { return "Demo" }
func (demoUser) UserAvatar() string  { return "" }
func (demoUser) UserRoles() []string { return []string{"Administrador"} }

func main() {
	ids, err := unixid.NewUnixID()
	if err != nil {
		panic(err)
	}
	db := orm.New(mem.New())
	broker := &mock.Broker{}

	dm, err := devicemanager.New(db, devicemanager.Deps{
		IDs:       ids,
		Publisher: broker,
		TenantID:  demoTenantID,
	})
	if err != nil {
		panic(err)
	}

	sm, err := staffmanager.New(db, staffmanager.Deps{
		IDs:         ids,
		Publisher:   broker,
		TenantID:    demoTenantID,
		ValidateRUT: trustedip.ValidateRUT,
		Devices:     devicemanager.IPLocator{Devices: dm, TenantID: demoTenantID},
	})
	if err != nil {
		panic(err)
	}

	if _, err := deviceseed.Load(dm, demoTenantID); err != nil {
		panic(err)
	}
	if _, err := seed.Load(sm, demoTenantID); err != nil {
		panic(err)
	}

	caller := loopback.WithTenant(demoTenantID, sm, dm)
	v, err := ui.Browser(caller, ids, demoTenantID)
	if err != nil {
		panic(err)
	}

	p := &platformd.Platform{
		AppName:   ui.Label + " — demo",
		User:      demoUser{},
		Modules:   []platformd.UIModule{v},
		DefaultID: ui.ID,
	}
	Append("body", p)
	select {}
}
