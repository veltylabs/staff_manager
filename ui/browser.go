package ui

import (
	staffmanager "github.com/veltylabs/staff_manager"
	"webtyp.com/dom"
	"webtyp.com/layout/crudview"
	"webtyp.com/layout/platformd"
	"webtyp.com/model"
	"webtyp.com/router"
	"webtyp.com/svg"
)

// StaffPanel es la vista CRUD de funcionarios como componente independiente:
// el panel que otras pantallas embeben (la pestaña "Datos y dispositivos" de
// "Personal" en appointment_booking). parentID se usa tal cual.
func StaffPanel(caller router.Caller, ids model.IDGenerator, parentID string) (dom.Component, error) {
	return crudview.New(crudview.Config{
		ParentID:  parentID,
		Presenter: staffmanager.NewView(caller),
		IDs:       ids,
	})
}

// Browser devuelve la pantalla principal de este módulo envuelta en un UIModule.
func Browser(caller router.Caller, ids model.IDGenerator, tenantID string) (platformd.UIModule, error) {
	panel, err := StaffPanel(caller, ids, ID)
	if err != nil {
		return nil, err
	}
	return platformd.NewUIModule(ID, Label, svg.Icon(ID), panel), nil
}
