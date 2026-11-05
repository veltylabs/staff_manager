package staffmanager

type DeviceReader interface {
	FindByIP(ip string) (string, bool)
}

// DeviceRegistrar registra un computador con su IP y devuelve su id: lo usa el alta
// del primer administrador (RegisterFirstAdmin). device_manager lo implementa.
//
// Ejemplo: id, err := registrar.RegisterComputer("Ana (admin)", "192.168.1.50")
type DeviceRegistrar interface {
	RegisterComputer(name, ip string) (deviceID string, err error)
}
