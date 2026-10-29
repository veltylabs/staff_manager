package staffmanager

import (
	"webtyp.com/model"
	"webtyp.com/router"
	"webtyp.com/rpc"
	"webtyp.com/view"
)

const (
	OpListStaff   = "list_staff"
	OpGetStaff    = "get_staff"
	OpUpsertStaff = "upsert_staff"
	OpDeleteStaff = "delete_staff"
)

// ModelName is this module's identity: mcp.HarvestOps qualifies every op
// below as "staff_manager.<name>" on the wire, and view.go's NewView passes
// this same constant as view.Ops.Module so the client composes the
// identical qualified name.
const ModelName = "staff_manager"

func (m *Module) ModelName() string { return ModelName }

func (m *Module) MountOperations(reg router.OperationRegistry) {
	res := model.ResourceOf(&StaffMemberModel)
	reg.Operation(OpListStaff, m.opListStaff).Requires(res, model.Read).Accepts(nil)
	reg.Operation(OpGetStaff, m.opGetStaff).Requires(res, model.Read).Accepts(&GetStaffArgs{})

	view.MountSave(reg, OpUpsertStaff, func(s *StaffMember) error {
		s.TenantId = m.tenantID // el tenant lo pone el servidor, nunca la petición (D32)
		out, err := m.UpsertStaff(*s)
		*s = out
		return err
	}).Requires(res, model.Create|model.Update)

	view.MountDelete(reg, OpDeleteStaff, func(ids []string) error {
		for _, id := range ids {
			if err := m.DeleteStaff(m.tenantID, id); err != nil {
				return err
			}
		}
		return nil
	}).Requires(res, model.Delete)
}

var _ router.OperationModule = (*Module)(nil)

func (m *Module) opListStaff(ctx router.Context) {
	list, err := m.ListStaff(m.tenantID)
	if err != nil {
		rpc.WriteError(ctx, err)
		return
	}
	l := make(StaffMemberList, len(list))
	for i := range list {
		v := list[i]
		l[i] = &v
	}
	if err := ctx.Encode(&l); err != nil {
		rpc.WriteError(ctx, err)
	}
}

func (m *Module) opGetStaff(ctx router.Context) {
	var args GetStaffArgs
	if err := ctx.Decode(&args); err != nil {
		rpc.WriteError(ctx, ErrCannotRead)
		return
	}
	s, err := m.GetStaff(m.tenantID, args.Id)
	if err != nil {
		rpc.WriteError(ctx, err)
		return
	}
	if err := ctx.Encode(&s); err != nil {
		rpc.WriteError(ctx, err)
	}
}
