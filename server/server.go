//go:build !wasm

// Package server es staff_manager tal como lo enchufa una app.
package server

import (
	"webtyp.com/appkit"
	trustedip "webtyp.com/auth/trusted_ip"
	"webtyp.com/model"

	staffmanager "github.com/veltylabs/staff_manager"
	"github.com/veltylabs/staff_manager/migrate"
)

// Module es staff_manager para appkit. Necesita un módulo que diga qué equipo
// tiene una IP (staffmanager.DeviceReader): en App.Modules va después de él.
//
// Ejemplo: Modules: []appkit.Module{devices.Module, staff.Module}
var Module = appkit.Module{Name: "staff_manager", Migrate: migrate.Migrate, Build: build}

func build(d appkit.Deps) (model.ModuleNaming, error) {
	devices, err := appkit.Find[staffmanager.DeviceReader](d)
	if err != nil {
		return nil, err
	}
	return staffmanager.New(d.DB, staffmanager.Deps{
		IDs: d.IDs, Publisher: d.Events, TenantID: d.Tenant,
		ValidateRUT: trustedip.ValidateRUT, Devices: devices,
	})
}
