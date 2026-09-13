package staffmanager

import (
	"webtyp.com/model"
	"webtyp.com/orm"
	"webtyp.com/router"
)

const (
	OpListStaff   = "list_staff"
	OpGetStaff    = "get_staff"
	OpUpsertStaff = "upsert_staff"
	OpDeleteStaff = "delete_staff"
)

func (m *Module) ModelName() string { return "staff_manager" }

func (m *Module) MountOperations(reg router.OperationRegistry) {
	reg.Operation(OpListStaff, m.opListStaff).Requires("staff_manager", model.Read).Accepts(&ListStaffArgs{})
	reg.Operation(OpGetStaff, m.opGetStaff).Requires("staff_manager", model.Read).Accepts(&GetStaffArgs{})
	reg.Operation(OpUpsertStaff, m.opUpsertStaff).Requires("staff_manager", model.Create|model.Update).Accepts(&StaffMember{})
	reg.Operation(OpDeleteStaff, m.opDeleteStaff).Requires("staff_manager", model.Delete).Accepts(&GetStaffArgs{})
}

var _ router.OperationModule = (*Module)(nil)

func (m *Module) opListStaff(ctx router.Context) {
	var args ListStaffArgs
	if err := ctx.Decode(&args); err != nil {
		ctx.WriteStatus(400)
		return
	}
	tenantID := args.TenantId
	if tenantID == "" {
		tenantID = m.tenantID
	}
	list, err := m.ListStaff(tenantID)
	if err != nil {
		ctx.WriteStatus(500)
		return
	}
	l := make(StaffMemberList, len(list))
	for i := range list {
		v := list[i]
		l[i] = &v
	}
	if err := ctx.Encode(&l); err != nil {
		ctx.WriteStatus(500)
	}
}

func (m *Module) opGetStaff(ctx router.Context) {
	var args GetStaffArgs
	if err := ctx.Decode(&args); err != nil {
		ctx.WriteStatus(400)
		return
	}
	s, err := m.GetStaff(args.TenantId, args.Id)
	if err != nil {
		if err == ErrNotFound {
			ctx.WriteStatus(404)
		} else {
			ctx.WriteStatus(500)
		}
		return
	}
	if err := ctx.Encode(&s); err != nil {
		ctx.WriteStatus(500)
	}
}

func (m *Module) opUpsertStaff(ctx router.Context) {
	var s StaffMember
	if err := ctx.Decode(&s); err != nil {
		ctx.WriteStatus(400)
		return
	}
	created, err := m.UpsertStaff(s)
	if err != nil {
		ctx.WriteStatus(400)
		return
	}
	if err := ctx.Encode(&created); err != nil {
		ctx.WriteStatus(500)
	}
}

func (m *Module) opDeleteStaff(ctx router.Context) {
	var args GetStaffArgs
	if err := ctx.Decode(&args); err != nil {
		ctx.WriteStatus(400)
		return
	}
	// Eliminar asignaciones staff_device primero, luego staff
	var s StaffMember
	s.Id = args.Id
	s.TenantId = args.TenantId
	if s.TenantId == "" {
		s.TenantId = m.tenantID
	}
	// Se necesita obtener y luego eliminar
	existing, err := m.GetStaff(s.TenantId, s.Id)
	if err != nil {
		ctx.WriteStatus(404)
		return
	}
	if err := m.db.Delete(&existing, orm.Eq("id", existing.Id), orm.Eq("tenant_id", existing.TenantId)); err != nil {
		ctx.WriteStatus(500)
		return
	}
	// también eliminar asignaciones
	var sd StaffDevice
	_ = m.db.Delete(&sd, orm.Eq("staff_id", existing.Id))
	ctx.WriteStatus(200)
}
