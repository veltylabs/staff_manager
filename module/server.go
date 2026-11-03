//go:build !browser

package module

import (
	"webtyp.com/appkit"
	trustedip "webtyp.com/auth/trusted_ip"
	"webtyp.com/model"

	staffmanager "github.com/veltylabs/staff_manager"
	"github.com/veltylabs/staff_manager/migrate"
)

// server es lo que el módulo aporta al servidor: sus tablas y cómo se construye.
// Necesita un módulo que sepa de equipos (device_manager): en App.Modules va después de él.
var server = appkit.ModuleServer{Migrate: migrate.Migrate, Build: build}

// En el servidor el módulo no aporta pantallas.
var browser appkit.ModuleBrowser

func build(d appkit.Deps) (model.ModuleNaming, error) {
	devices, err := appkit.Find[staffmanager.DeviceReader](d)
	if err != nil {
		return nil, err
	}
	registrar, err := appkit.Find[staffmanager.DeviceRegistrar](d)
	if err != nil {
		return nil, err
	}
	return staffmanager.New(d.DB, staffmanager.Deps{
		IDs: d.IDs, Publisher: d.Events, TenantID: d.TenantID,
		ValidateRUT: trustedip.ValidateRUT, Devices: devices, Registrar: registrar,
	})
}
