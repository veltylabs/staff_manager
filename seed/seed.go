package seed

import (
	staffmanager "github.com/veltylabs/staff_manager"
)

// Data contiene los datos semilla creados para este módulo.
type Data struct {
	Staff []staffmanager.StaffMember
}

// Nombres de especialidades para las semillas.
const (
	specialtyMedicinaGeneral = "Medicina General"
	specialtyTraumatologia    = "Traumatología"
	specialtyEcografia        = "Ecografía"
)

// Load inserta datos de demostración a través de los métodos del módulo.
func Load(m *staffmanager.Module, tenantID string) (Data, error) {
	members := []staffmanager.StaffMember{
		{
			TenantId:  tenantID,
			Rut:       "17654321-3",
			Name:      "Dra. Ana Rojas",
			UserId:    "demo-user-ana",
			Specialty: specialtyMedicinaGeneral,
			Role:      "Médico",
			IsActive:  true,
		},
		{
			TenantId:  tenantID,
			Rut:       "9876543-3",
			Name:      "Dr. Luis Fuentes",
			UserId:    "demo-user-luis",
			Specialty: specialtyTraumatologia,
			Role:      "Médico",
			IsActive:  true,
		},
		{
			TenantId:  tenantID,
			Rut:       "20123456-5",
			Name:      "Dra. Carla Díaz",
			UserId:    "demo-user-carla",
			Specialty: specialtyEcografia,
			Role:      "Tecnólogo",
			IsActive:  true,
		},
	}

	created := make([]staffmanager.StaffMember, 0, len(members))
	for _, member := range members {
		res, err := m.UpsertStaff(member)
		if err != nil {
			return Data{}, err
		}
		created = append(created, res)
	}

	return Data{Staff: created}, nil
}
