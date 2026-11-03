package staffmanager

// RegisterFirstAdmin completa al primer administrador que appkit ya creó (usuario y
// rol): su ficha de funcionario y el equipo desde el que se dio de alta, asignado a
// él. En el loopback registra también la otra cara (127.0.0.1 ↔ ::1): el navegador
// puede entrar por cualquiera de las dos según cómo resuelva "localhost".
// Implementa appkit.FirstAdminRegistrar.
//
// Ejemplo: err := m.RegisterFirstAdmin(user.Id, "11111111-1", "Ana", "192.168.1.50")
func (m *Module) RegisterFirstAdmin(userID, code, name, ip string) error {
	if m.devices == nil {
		return depErr("Devices")
	}
	if m.registrar == nil {
		return depErr("Registrar")
	}
	member, err := m.UpsertStaff(StaffMember{TenantId: m.tenantID, UserId: userID, Rut: code, Name: name, IsActive: true})
	if err != nil {
		return err
	}
	for _, addr := range firstAdminIPs(ip) {
		deviceID, ok := m.devices.FindByIP(addr)
		if !ok {
			if deviceID, err = m.registrar.RegisterComputer(name+firstAdminDeviceSuffix, addr); err != nil {
				return err
			}
		}
		if err := m.AssignDevice(member.Id, deviceID); err != nil {
			return err
		}
	}
	return nil
}

// firstAdminDeviceSuffix va tras el nombre del equipo del instalador. Sin ":" — el
// charset de nombres de device_manager no lo permite.
const firstAdminDeviceSuffix = " (admin)"

const (
	loopbackIPv4 = "127.0.0.1"
	loopbackIPv6 = "::1"
)

// firstAdminIPs son las IP que quedan asignadas al instalador: la suya y, en el
// loopback, también la otra cara.
func firstAdminIPs(ip string) []string {
	switch ip {
	case loopbackIPv4:
		return []string{ip, loopbackIPv6}
	case loopbackIPv6:
		return []string{ip, loopbackIPv4}
	}
	return []string{ip}
}
