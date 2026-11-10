# staff_manager
<img src="docs/img/badges.svg">

## ¿Por qué existe?

`staff_manager` es el registro del personal de una clínica: vincula a cada persona con su cuenta de usuario (por RUT), su
especialidad, su rol y los equipos desde los que trabaja. Lo encuentras cuando otro módulo, como `clinical_encounter`, necesita
saber quién es el profesional de una cuenta, o cuando el inicio de sesión por red local pregunta si una IP es de confianza.

- **Sin esto:** cada módulo guarda su propia lista de personas y el inicio de sesión por LAN no sabe qué equipo es de quién.
- **Con esto:** hay una sola fuente de verdad; la persistencia la hace [`crud`](https://github.com/webtyp/crud) y las reglas del
  negocio (RUT único, solo se editan los campos listados) viven en los archivos `rule_*.go` ([arquitectura](docs/ARCHITECTURE.md)).

## Antes de comenzar

- Go instalado y acceso a los módulos `webtyp.com/*` (`GOPRIVATE=webtyp.com/*,github.com/veltylabs/*`).
- Una base de datos conectada como `*orm.DB`; las tablas las crea `migrate.Migrate` (usa `ddl.Sync`, así que agrega solo lo que falta).
- Un generador de ids, el id del inquilino (`TenantID`), un validador de RUT (`input.NormalizeRUT`) y un lector de equipos
  (`DeviceReader`, ver [interfaces.go](interfaces.go)): la app conecta `device_manager` por `Deps.Devices`.

## Guía paso a paso

### Paso 1: crear el módulo

**Objetivo:** conectar el registro de personal a una base en memoria.

```go
package main

import (
	staffmanager "github.com/veltylabs/staff_manager"
	"webtyp.com/fmt"
	"webtyp.com/input"
	"webtyp.com/orm"
	"webtyp.com/storage/mem"
	"webtyp.com/unixid"
)

// noDevices reemplaza a device_manager: ningún equipo está registrado, así que ninguna IP es de confianza.
type noDevices struct{}

func (noDevices) FindByIP(string) (string, bool) { return "", false }

func main() {
	ids, err := unixid.NewUnixID()
	if err != nil {
		fmt.Println(err)
		return
	}
	mod, err := staffmanager.New(orm.New(mem.New()), staffmanager.Deps{
		IDs: ids, TenantID: "acme", ValidateRUT: input.NormalizeRUT, Devices: noDevices{},
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(mod.ModelName())
}
```

**Resultado esperado:** imprime `staff_manager`.

### Paso 2: operar el registro

**Objetivo:** crear una persona y comprobar que un RUT repetido se rechaza.

Las operaciones son las de [`crud`](https://github.com/webtyp/crud), pedidas con `router.Operation{Module, Entity, Action}`
(por ejemplo `staff_member.create`). Los ejemplos completos, con sus resultados, están en
[tests/test_staff_duplicate_test.go](tests/test_staff_duplicate_test.go) y [tests/test_staff_update_test.go](tests/test_staff_update_test.go).

**Resultado esperado:** la persona queda guardada; una segunda con el mismo RUT recibe un error en el campo del RUT.

### Paso 3: ver el módulo en pantalla

**Objetivo:** abrir la demo con datos en memoria. La demo está en [seed/demo.go](seed/demo.go) y se abre con `webtyp` en la raíz del módulo.

**Resultado esperado:** una lista del personal con búsqueda.

## Puertos que satisface este módulo

| Método | Contrato | Para qué |
|---|---|---|
| `IsTrustedIP(userID, ip)` | `auth.TrustedIPStore` | ¿la IP es de un equipo asignado a esa persona? (inicio de sesión por LAN) |
| `StaffExists(tenantID, staffID)` | `StaffReader` | ¿existe esa persona en el inquilino? `(false, nil)` si no existe |
| `StaffCard(tenantID, staffID)` | `StaffReader` | nombre, especialidad y si está activa, sin datos personales |
| `StaffByUser(tenantID, userID)` | `DoctorReader` | el profesional de una cuenta (lo usa `clinical_encounter`) |

## Operaciones: sin y con `crud`

| Antes (manejadores a mano) | Después (declaración + reglas) |
|---|---|
| `list_staff` / `get_staff` | `staff_member.read` |
| `upsert_staff` | `staff_member.create` / `staff_member.update` |
| `delete_staff` | `staff_member.delete` |

## Mapa de documentación

| Orden | Documento | Qué pregunta responde y cuándo leerlo |
|---|---|---|
| 1 | [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) | ¿Qué entidades hay y quién decide qué? Léelo antes de tocar una regla. |
| 2 | [docs/TESTS.md](docs/TESTS.md) | ¿Qué comportamiento está probado? Léelo antes de cambiar una regla. |
