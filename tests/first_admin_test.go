package tests

import "testing"

// TestRegisterFirstAdmin recorre las ramas del alta: IP de LAN (un equipo), loopback
// (las dos caras) y un equipo que ya existía (se reutiliza, no se registra otro).
func TestRegisterFirstAdmin(t *testing.T) {
	cases := []struct {
		name           string
		known          map[string]string
		ip             string
		wantRegistered []string
		trusted        []string
	}{
		{"lan", map[string]string{}, "192.168.1.50", []string{"192.168.1.50"}, []string{"192.168.1.50"}},
		{"loopback v4", map[string]string{}, "127.0.0.1", []string{"127.0.0.1", "::1"}, []string{"127.0.0.1", "::1"}},
		{"loopback v6, v4 already known", map[string]string{"127.0.0.1": "pc-1"}, "::1", []string{"::1"}, []string{"::1", "127.0.0.1"}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			devices := &fakeDevices{byIP: c.known}
			m := setup(t, devices)
			if err := m.RegisterFirstAdmin("user-1", "12345678-5", "Ana", c.ip); err != nil {
				t.Fatalf("RegisterFirstAdmin: %v", err)
			}
			if len(devices.registered) != len(c.wantRegistered) {
				t.Fatalf("registered = %v, want %v", devices.registered, c.wantRegistered)
			}
			for i, ip := range c.wantRegistered {
				if devices.registered[i] != ip {
					t.Errorf("registered[%d] = %q, want %q", i, devices.registered[i], ip)
				}
			}
			for _, ip := range c.trusted {
				if !m.IsTrustedIP("user-1", ip) {
					t.Errorf("IsTrustedIP(user-1, %q) = false, want true", ip)
				}
			}
		})
	}
}
