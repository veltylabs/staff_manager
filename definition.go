package staffmanager

import (
	"webtyp.com/input"
	"webtyp.com/model"
)

var StaffMemberModel = model.Definition{
	Name: "staff_member",
	Allow: model.Create | model.Read | model.Update | model.Delete,
	Rules: model.Rules{
		Create: model.AddMyBusinessRule,
		Update: model.AddMyBusinessRule,
		Delete: model.AddMyBusinessRule,
		Read:   model.HaveNoBusinessRule,
	},
	Fields: model.Fields{
		{Name: "id", Type: model.Text(), DB: &model.FieldDB{PK: true}, OmitEmpty: true},
		{Name: "tenant_id", Type: model.Text(), NotNull: true, ServerOnly: true},
		{Name: "user_id", Type: model.Text(), NotNull: true, ServerOnly: true},
		{Name: "rut", Type: input.Rut(), Label: "Rut", NotNull: true, Permitted: model.Permitted{Minimum: 1, Maximum: 12}, Item: model.ItemDescription},
		{Name: "name", Type: input.Text(), Label: "Name", NotNull: true, Permitted: model.Permitted{Minimum: 1, Maximum: 255}, Item: model.ItemLabel},
		{Name: "specialty", Type: input.Text(), Label: "Specialty", Permitted: model.Permitted{Maximum: 120}},
		{Name: "role", Type: input.Text(), Label: "Role", Permitted: model.Permitted{Maximum: 60}},
		{Name: "is_active", Type: input.Checkbox(), Label: "Active"},
		{Name: "updated_at", Type: model.NanoTime(), OmitEmpty: true, ServerOnly: true},
	},
}

var StaffDeviceModel = model.Definition{
	Name: "staff_device",
	Fields: model.Fields{
		{Name: "staff_id", Type: model.Text(), Ref: &StaffMemberModel, DB: &model.FieldDB{PK: true, RefColumn: "id"}, NotNull: true},
		{Name: "device_id", Type: model.Text(), DB: &model.FieldDB{PK: true}, NotNull: true},
	},
}

const (
	ErrNotFound model.NotFound = "staff not found"
)

var Definition = model.ModuleDefinition{
	Name: "staff_manager",
	Entities: []model.Definition{
		StaffMemberModel,
		StaffDeviceModel,
	},
}
