---
PLAN: "feat: specialty and role on a staff member, plus the StaffExists port"
EXECUTOR: jules
REVIEWER: none
STATUS: review
SESSION: 17733230379642182941
PR: https://github.com/veltylabs/staff_manager/pull/1
---

> Este plan se distribuye a través del flujo de trabajo CodeJob. Ver habilidad: **agents-workflow**.

# Plan — `specialty`, `role`, y `StaffExists`

Eres un agente sin **contexto previo** y tienes **únicamente este repositorio**
(`github.com/veltylabs/staff_manager`). Todo lo que necesitas está en línea.

## 1. El problema

Dos brechas, ambas descubiertas al conectar este módulo junto a
`github.com/veltylabs/appointment_booking` en una aplicación de clínica.

### 1.1 Nada puede consultar si un miembro del personal existe

`appointment_booking` rechaza crear una reserva para un profesional que no puede
verificar. Declara la pregunta como un puerto estrecho que no importa:

```go
// appointment_booking/service.go
// StaffReader verifica que un miembro del staff existe y pertenece al tenant.
type StaffReader interface {
	StaffExists(tenantId, staffId string) (bool, error)
}
```

y su propio README documenta la conexión esperada como
`Staff: staffmodule.New(db)  // implementa StaffReader`, es decir, el propio módulo
de personal satisface el puerto estructuralmente, sin adaptador.

`*Module` no posee ese método actualmente. Cada aplicación que desee reservar una
cita debe escribir el mismo adaptador de seis líneas sobre `GetStaff`,
comparando `ErrNotFound` manualmente. Ese adaptador es código de aplicación haciendo
el trabajo del módulo, y se escribirá una vez por aplicación.

### 1.2 Un miembro del personal no tiene especialidad ni rol

`StaffMemberModel` hoy en día contiene `id`, `tenant_id`, `user_id`, `rut`, `name`,
`is_active`, `updated_at`. No hay dónde registrar que alguien es radiólogo.

La consecuencia ya es visible en el ecosistema: el módulo hermano
`github.com/veltylabs/clinical_encounter` escribe una columna
`doctor_specialty_snapshot` en cada registro médico, y no hay fuente para ese valor;
se escribe a mano en el punto de uso, por registro, sin nada que mantenga unificadas
dos escrituras de la misma especialidad. Una pantalla de reservas que desee ofrecer
"quién realiza ecografías" no tiene por qué agrupar.

## 2. Puerta de diseño

### Arte previo

- **Django** `Model.objects.filter(...).exists()` y **Rails**
  `Model.exists?(id)` exponen la existencia como una consulta booleana de primer orden,
  separada de obtener la fila, precisamente porque los llamantes que solo necesitan la
  respuesta no deben pagar ni procesar un registro completo.
- **Ent** (Go) genera `Query().Exist(ctx)` por entidad por la misma razón.
- Los tres retornan únicamente un booleano más un error. Coincidimos exactamente con esa forma.
  Diferimos en una cosa: el nuestro está delimitado por tenant (`tenantId` es el primer
  argumento), porque en este ecosistema ninguna lectura es global; cada método existente
  en este módulo ya recibe el tenant primero.

Para los dos campos nuevos no hay arte previo interesante del cual diferir: son
columnas simples en el propio registro del módulo.

### Prueba del nombre para principiantes

- "¿Existe este personal en este tenant?" → `StaffExists(tenantId, staffId)`.
  El puerto que declara `appointment_booking` ya utiliza este nombre y firma exactos;
  elegir uno diferente significaría que cada consumidor debe escribir un adaptador, que es el
  defecto que se está corrigiendo.
- "La especialidad del miembro del personal" → `member.Specialty`.
- "El rol del miembro del personal" → `member.Role`.

### Libro mayor de complejidad

| | cambio |
|---|---|
| Conceptos | `+0` — la existencia y la especialidad ya son conceptos en este dominio; ninguno introduce uno nuevo |
| Archivos | `+1` (`tests/staff_test.go`); `model.go`, `module.go`, `migrate/migrate.go` son editados |
| Líneas en sitios de llamada | `−6 por aplicación consumidora` (el adaptador `StaffReader` escrito a mano desaparece), `+1` aquí |
| Formas de hacerlo | `−1` — "¿es real este profesional?" tenía cero respuestas soportadas y una copiada; ahora tiene una |
| Neto | negativo |

### Dónde pertenece

Aquí. `StaffExists` es una pregunta sobre la propia tabla de este módulo, y los dos
campos son atributos del propio registro de este módulo. Poner la especialidad en
cualquier otro lugar —una tabla de búsqueda en el catálogo, una columna en el registro
clínico— es lo que produjo `doctor_specialty_snapshot` sin fuente en primer lugar.

### Qué elimina

Nada en este repositorio; esto es aditivo. Lo que permite eliminar es el adaptador
`StaffReader` por aplicación, en esas aplicaciones.

## 3. Decisiones ya tomadas — no reevaluar

1. **`specialty` es texto libre (`input.Text()`), no un slug y no un enum.**
   Una columna de tipo slug necesitaría una tabla slug→etiqueta, y la única instancia existente
   de eso en este ecosistema es un `switch` en español codificado en una demo —una aplicación
   codificando datos que la base de datos podría proveer. Un consumidor que desee agrupar
   profesionales por especialidad agrupa por los valores distintos que lee desde `list_staff`.
   Este módulo no procesa ningún lenguaje humano propio.
2. **`role` también es texto libre**, y **no** es un permiso. La autorización en
   este ecosistema es `webtyp.com/rbac`, claveada por usuario, recurso y acción.
   `role` aquí es un título de trabajo descriptivo mostrado a una persona ("Médico",
   "TENS", "Administrativo"). No lo conecte a nada que otorgue acceso.
3. **Ambos campos son opcionales** (sin `NotNull`). Las filas existentes no tienen ninguno, y
   una columna requerida rompería todos los registros actuales.
4. **`Item()` no se modifica.** `StaffMember.Item()` actualmente retorna
   `Description: m.Rut`, y una pantalla de personal existente renderiza eso. Cambiarlo
   para mostrar la especialidad es una decisión de interfaz para la aplicación consumidora, no
   de este plan.

## 4. Etapas

### Etapa 1 — los dos campos, en `model.go`

Agregar ambos a `StaffMemberModel.Fields`, **inmediatamente después de `name`** y antes
de `is_active`:

```go
{Name: "specialty", Type: input.Text(), Permitted: model.Permitted{Maximum: 120}},
{Name: "role", Type: input.Text(), Permitted: model.Permitted{Maximum: 60}},
```

Escribir un comentario sobre ellos registrando las decisiones 3.1 y 3.2 de este plan —
que los valores son texto libre a propósito, y que `role` es descriptivo y
nunca una entrada de autorización. Ese razonamiento debe perdurar más allá de este archivo de plan.

Regenerar `model_orm.go` con `ormc`. **Nunca editar `model_orm.go` a mano.** La
estructura generada obtiene `Specialty` y `Role` (mayúsculas directas de los nombres de columna).

### Etapa 2 — transmitir los campos a través de `UpsertStaff`, en `module.go`

`UpsertStaff` tiene **dos** rutas de actualización, y ambas copian campos campo por campo
en un registro `existing`. Un nuevo campo agregado a solo una de ellas desaparece silenciosamente
en la otra ruta.

**Ruta A — actualización por id** (`if member.Id != ""`, después de que `ReadOne` tiene éxito):

```go
existing.Name = member.Name
existing.Rut = member.Rut
existing.Specialty = member.Specialty   // AGREGAR
existing.Role = member.Role             // AGREGAR
existing.IsActive = member.IsActive
existing.UserId = member.UserId
existing.UpdatedAt = member.UpdatedAt
```

**Ruta B — actualización por tenant+rut** (el bloque después de que
`Where("tenant_id").Eq(...).Where("rut").Eq(...)` tiene éxito):

```go
existing.Name = member.Name
existing.Specialty = member.Specialty   // AGREGAR
existing.Role = member.Role             // AGREGAR
existing.IsActive = member.IsActive
existing.UpdatedAt = member.UpdatedAt
```

La ruta de creación no necesita cambios — escribe `member` completo.

**Aceptación:** `grep -n "existing.Specialty" module.go` → exactamente 2 coincidencias.
Igual para `existing.Role`.

### Etapa 3 — `StaffExists`, en `module.go`

Agregar inmediatamente después de `GetStaff`, dado que es la misma consulta reducida a un
booleano:

```go
// StaffExists reporta si un miembro del personal con este id pertenece a este
// tenant. Satisface el puerto estrecho StaffReader que un módulo de programación
// declara en su propio lado (StaffExists(tenantId, staffId) (bool, error)) —
// estructuralmente, sin adaptador y sin importación en ninguna dirección.
//
// Una fila ausente es (false, nil), NO un error: "este id no es nuestro" es
// la respuesta que el llamante solicitó. Solo un fallo real de almacenamiento retorna un
// error no nulo, para que un llamante nunca confunda una base de datos caída con un "no" limpio.
func (m *Module) StaffExists(tenantID, staffID string) (bool, error) {
	if tenantID == "" || staffID == "" {
		return false, nil
	}
	_, err := m.GetStaff(tenantID, staffID)
	if err != nil {
		if err == ErrNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
```

**Evitar errores.** No "optimizar" esto haciendo la consulta en línea y retornando
`err == nil`. `GetStaff` es el único lugar que mapea `orm.ErrNotFound` al
sentinela de dominio; una segunda copia de ese mapeo es cómo ambos divergen. Y no
reducir el error a un `bool` simple — un llamante que no puede distinguir "ausente" de
"base de datos inalcanzable" reservará una cita para un miembro del personal que nunca
verificó.

### Etapa 4 — `migrate/` debe AGREGAR COLUMNAS, no solo crear tablas

Esta es la etapa con mayor probabilidad de error, porque el código actual parece
terminado.

`migrate/migrate.go` hoy llama a `CreateTable`, que se compila en
`CREATE TABLE IF NOT EXISTS`. Contra una base de datos donde `staff_member` ya
existe —lo cual sucede en cada despliegue que ha funcionado— es una **no-operación**, por lo que
las dos nuevas columnas nunca aparecerían y cada escritura fallaría por columna desconocida.

`webtyp.com/ddl` ya tiene la operación correcta: `(*ddl.DB).Sync(models...)`
emite `CreateTable` y luego, para cada columna en el modelo que la tabla no
posee, un `OpAddColumn` aditivo — dentro de una transacción cuando el backend lo soporte.
Reemplazar el cuerpo:

```go
func Migrate(conn ddl.Execer, ddlCompiler ddl.Compiler) error {
	// Sync, no CreateTable: CreateTable se compila a CREATE TABLE IF NOT
	// EXISTS y es una no-operación contra una tabla que ya existe, por lo que una columna
	// agregada a un modelo nunca llegaría a una base de datos desplegada. Sync crea la
	// tabla cuando está ausente y agrega las columnas faltantes cuando no lo está.
	return ddl.New(conn, ddlCompiler).Sync(
		&staffmanager.StaffMember{},
		&staffmanager.StaffDevice{},
	)
}
```

`Sync` es únicamente aditivo —nunca elimina ni restringe una columna— por lo que es seguro de
ejecutar repetidamente y seguro de ejecutar contra producción.

**Evitar errores.** Los módulos hermanos (`device_manager`, `item_catalog`,
`clinical_encounter`) aún llaman a `CreateTable` en sus propios paquetes `migrate/`.
Tienes **únicamente este repositorio**; no intentes cambiarlos, y
no trates su código como el patrón a copiar aquí. La diferencia es
deliberada: este es el módulo agregando columnas a una tabla existente.

### Etapa 5 — pruebas, en `tests/`

`tests/` actualmente contiene `setup_test.go` y `trust_test.go`. Leer
`setup_test.go` primero y reutilizar su construcción de módulo y sus fakes — no
construir un segundo arnés.

Nuevo archivo `tests/staff_test.go`:

| Prueba | Afirma |
|---|---|
| `TestStaffExists_True` | un miembro insertado/actualizado retorna `(true, nil)` |
| `TestStaffExists_UnknownID` | `(false, nil)` — **no** es un error |
| `TestStaffExists_WrongTenant` | un miembro del tenant A es `(false, nil)` para el tenant B |
| `TestStaffExists_EmptyArgs` | `("", "")` y `("t1", "")` ambos `(false, nil)` |
| `TestUpsertStaff_PersistsSpecialtyAndRole` | crear, luego `GetStaff`, ambos campos se persisten correctamente |
| `TestUpsertStaff_UpdateByIDKeepsSpecialty` | ruta de actualización A cambia la especialidad y la lectura muestra el nuevo valor |
| `TestUpsertStaff_UpdateByRutKeepsSpecialty` | ruta de actualización B — upsert **sin** `Id`, coincidiendo en tenant+rut — cambia la especialidad y la lectura muestra el nuevo valor |

Las últimas dos son las que detectan un campo llevado a través de solo una de las dos
rutas de actualización. Ambas deben existir.

Agregar una afirmación en tiempo de compilación, en `tests/staff_test.go`, que fije la forma
del puerto para que un cambio de firma futuro falle la compilación aquí en lugar de en un
consumidor:

```go
// El puerto que appointment_booking declara en su propio lado. Redeclarado localmente a
// propósito: este repositorio no debe depender de un módulo de programación para probar que
// satisface una interfaz estructural.
type staffReader interface {
	StaffExists(tenantId, staffId string) (bool, error)
}

var _ staffReader = (*staffmanager.Module)(nil)
```

Ejecutar `gotest ./...` — todo en verde, incluyendo `trust_test.go` sin modificar.

### Etapa 6 — documentación

`README.md` tiene cuatro líneas y una tabla de Ops. Extenderlo:

- Indicar qué es un miembro del personal: el registro delimitado por tenant que vincula un usuario (por
  RUT) a dispositivos, **y** almacena su especialidad y rol.
- Agregar una sección **"Puertos que satisface este módulo"** nombrando `IsTrustedIP`
  (`auth.TrustedIPStore`) y el nuevo `StaffExists`, cada uno con su pregunta de una línea
  y su firma.
- Agregar una sección de **Esquema** listando las columnas de `StaffMember`, marcando `specialty`
  y `role` como texto libre e indicando que `role` es descriptivo, nunca una
  entrada de autorización.
- Agregar una nota indicando que `migrate.Migrate` utiliza `Sync` y es aditivo.

No vincular ningún documento permanente a `docs/PLAN.md` — se elimina cuando
esto se fusiona.

## 5. Tabla de etapas

| # | Etapa | Archivos | Aceptación |
|---|---|---|---|
| 1 | Campos | `model.go`, `model_orm.go` (ormc) | `Specialty` y `Role` en la estructura generada |
| 2 | Rutas de Upsert | `module.go` | `grep -c "existing.Specialty" module.go` → 2 |
| 3 | `StaffExists` | `module.go` | fila ausente → `(false, nil)` |
| 4 | Migración aditiva | `migrate/migrate.go` | `grep -n "CreateTable" migrate/migrate.go` → vacío |
| 5 | Pruebas | `tests/staff_test.go` (nuevo) | 7 casos + la afirmación del puerto |
| 6 | Docs | `README.md` | puertos y esquema documentados |

## 6. Criterios de aceptación

- `grep -n "CreateTable" migrate/migrate.go` → **vacío**; `Sync` es la única
  llamada DDL.
- `grep -c "existing.Specialty" module.go` → `2`; igual para `existing.Role`.
- `var _ staffReader = (*staffmanager.Module)(nil)` compila en
  `tests/staff_test.go`.
- `StaffExists` retorna `(false, nil)` —nunca un error— para un id desconocido, un
  id de otro tenant, y argumentos vacíos.
- `model_orm.go` cambiado solo mediante regeneración por `ormc`.
- `StaffMember.Item()` no cambia.
- Sin lista de especialidades codificada a mano, tabla de slugs o `switch` de etiquetas en ninguna parte del
  repositorio: `grep -rniE "radiolog|traumatolog|kinesiolog" .` → **vacío**.
- `gotest ./...` en verde, `tests/trust_test.go` no modificado.
