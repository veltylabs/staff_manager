package ui

import (
	staffmanager "github.com/veltylabs/staff_manager"
	"webtyp.com/lang"
	"webtyp.com/model"
	"webtyp.com/view"
)

const descStaff lang.Text = "The clinic's staff: profiles, contact details and professional credentials."

// Screens declara las pantallas que este módulo ofrece al menú de la app.
func Screens(d view.BrowserDeps) []view.Screen {
	caller := d.Caller
	return []view.Screen{
		{
			ID:         ID,
			Label:      Label,
			Icon:       ID,
			Permission: model.ResourceOf(&staffmanager.StaffMemberModel),
			Presenter:  staffmanager.NewView(caller, view.WithDescription(descStaff)),
		},
	}
}
