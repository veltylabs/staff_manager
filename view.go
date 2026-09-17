package staffmanager

import (
	"webtyp.com/model"
	"webtyp.com/router"
	"webtyp.com/view"
)

func (m *StaffMember) Item() view.Item {
	return view.Item{ID: m.Id, Label: m.Name, Description: m.Rut}
}

const titleStaff = "Funcionarios"

func NewView(caller router.Caller) view.Presenter {
	b := view.NewCallerLister(caller,
		view.Ops{Module: ModelName, List: OpListStaff, Save: OpUpsertStaff, Delete: OpDeleteStaff},
		func() model.ModelSlice { return &StaffMemberList{} })
	return view.New(b, &StaffMember{}, view.WithTitle(titleStaff))
}
