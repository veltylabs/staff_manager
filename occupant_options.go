package staffmanager

import (
	"webtyp.com/fmt"
)

// OccupantOptions entrega el personal activo como opciones {id, nombre}: es lo
// que room_layout pide como ocupantes de un turno (roomlayout.OccupantReader).
//
// Ejemplo: opts, err := staff.OccupantOptions(tenantID)

func (m *Module) OccupantOptions(tenantID string) ([]fmt.KeyValue, error) {
	list, err := m.ListStaff(tenantID, "", "", 1000)
	if err != nil {
		return nil, err
	}
	out := make([]fmt.KeyValue, 0, len(list))
	for _, s := range list {
		if s.IsActive {
			out = append(out, fmt.KeyValue{Key: s.Id, Value: s.Name})
		}
	}
	return out, nil
}
