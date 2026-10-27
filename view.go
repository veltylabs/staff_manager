package staffmanager

import (
	"webtyp.com/model"
	"webtyp.com/router"
	"webtyp.com/view"
)

const titleStaff = "Funcionarios"

func NewView(caller router.Caller, opts ...view.Option) view.Presenter {
	caller = view.NewRecords().Caller(caller)
	b := view.NewCallerLister(caller,
		view.Ops{Module: ModelName, List: OpListStaff, Save: OpUpsertStaff, Delete: OpDeleteStaff},
		func() model.ModelSlice { return &StaffMemberList{} })
	allOpts := append([]view.Option{view.WithTitle(titleStaff)}, opts...)
	return view.New(b, &StaffMember{}, allOpts...)
}
