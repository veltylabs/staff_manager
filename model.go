package staffmanager

import (
	"github.com/veltylabs/device_manager"
	"webtyp.com/input"
	"webtyp.com/model"
)

var StaffMemberModel = model.Definition{
	Name: "staff_member",
	Fields: model.Fields{
		{Name: "id", Type: model.Text(), DB: &model.FieldDB{PK: true}, OmitEmpty: true},
		{Name: "tenant_id", Type: model.Text(), NotNull: true},
		{Name: "user_id", Type: model.Text(), NotNull: true},
		{Name: "rut", Type: input.Rut(), Label: "Rut", NotNull: true, Permitted: model.Permitted{Minimum: 1, Maximum: 12}, Item: model.ItemDescription},
		{Name: "name", Type: input.Text(), Label: "Name", NotNull: true, Permitted: model.Permitted{Minimum: 1, Maximum: 255}, Item: model.ItemLabel},
		// Specialty y role son valores de texto libre sin enums ni slugs. Role es
		// puramente descriptivo (ej. título de trabajo) y nunca una entrada de autorización.
		// Ambos campos son opcionales para que los registros existentes sigan siendo válidos.
		{Name: "specialty", Type: input.Text(), Label: "Specialty", Permitted: model.Permitted{Maximum: 120}},
		{Name: "role", Type: input.Text(), Label: "Role", Permitted: model.Permitted{Maximum: 60}},
		{Name: "is_active", Type: input.Checkbox(), Label: "Active"},
		{Name: "updated_at", Type: model.NanoTime(), OmitEmpty: true},
	},
}

var StaffDeviceModel = model.Definition{
	Name: "staff_device",
	Fields: model.Fields{
		{Name: "staff_id", Type: model.Text(), Ref: &StaffMemberModel, DB: &model.FieldDB{PK: true, RefColumn: "id"}, NotNull: true},
		{Name: "device_id", Type: model.Text(), Ref: &devicemanager.DeviceModel, DB: &model.FieldDB{PK: true, RefColumn: "id"}, NotNull: true},
	},
}

var GetStaffArgsModel = model.Definition{
	Name: "get_staff_args",
	Fields: model.Fields{
		{Name: "tenant_id", Type: model.Text()},
		{Name: "id", Type: model.Text()},
	},
}

const (
	ErrNotFound   model.NotFound = "staff not found"
	ErrCannotRead model.Invalid  = "staff_manager: cannot read the request"
)

const (
	TopicStaffCreated = "staff_manager.staff.created"
	TopicStaffUpdated = "staff_manager.staff.updated"
)
