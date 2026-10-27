package ui

import (
	staffmanager "github.com/veltylabs/staff_manager"
	"webtyp.com/lang"
	"webtyp.com/router"
	"webtyp.com/view"
)

const descStaff lang.Text = "Gestión de personal: perfiles, datos de contacto y credenciales profesionales."

// Screens declara las pantallas que este módulo ofrece al menú de la app.
func Screens(caller router.Caller) []view.Screen {
	return []view.Screen{
		{
			ID:         ID,
			Label:      Label,
			Icon:       ID,
			Permission: "staff_manager",
			Presenter:  staffmanager.NewView(caller, view.WithDescription(descStaff)),
		},
	}
}
