//go:build browser

package main

import (
	staffmanager "github.com/veltylabs/staff_manager"
	"github.com/veltylabs/staff_manager/seed"
	"github.com/veltylabs/staff_manager/ui"
	"webtyp.com/auth/trusted_ip"
	. "webtyp.com/dom"
	"webtyp.com/events/mock"
	"webtyp.com/html"
	"webtyp.com/orm"
	"webtyp.com/router/loopback"
	"webtyp.com/storage/mem"
	"webtyp.com/unixid"
)

// demoTenantID es el único tenant de esta demo ejecutable en el navegador.
const demoTenantID = "demo"

func main() {
	ids, err := unixid.NewUnixID()
	if err != nil {
		panic(err)
	}
	db := orm.New(mem.New())
	broker := &mock.Broker{}

	sm, err := staffmanager.New(db, staffmanager.Deps{
		IDs:         ids,
		Publisher:   broker,
		TenantID:    demoTenantID,
		ValidateRUT: trustedip.ValidateRUT,
		Devices:     demoDevices{},
	})
	if err != nil {
		panic(err)
	}

	if _, err := seed.Load(sm, demoTenantID); err != nil {
		panic(err)
	}

	caller := loopback.WithTenant(demoTenantID, sm)
	screens := ui.Screens(caller)

	appNode := html.Div().Child(
		html.H1().Text(string(screens[0].Label)),
		html.P().Text("Demo render sin webtyp.com/layout"),
	)

	Append("body", appNode)
	select {}
}

// demoDevices reemplaza a device_manager en esta demo: ningún equipo está registrado, así que
// ninguna IP inicia sesión. La app real conecta device_manager por Deps.Devices.
type demoDevices struct{}

func (demoDevices) FindByIP(string) (string, bool) { return "", false }
