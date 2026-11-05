Las fechas que ve una persona salen de la zona de la instalación (`Deps.Zone` / `BrowserDeps.Zone`); un día guardado como clave usa las funciones UTC.

## Integración con otros módulos

Los contratos (interfaces) que este módulo requiere de otros módulos para funcionar (como `DeviceReader` o `DeviceRegistrar`) se definen exclusivamente en el archivo `interfaces.go`.

### Contratos que satisface

- `IsTrustedIP`: Satisface `auth.TrustedIPStore`.
- `StaffExists`: Satisface `StaffReader.StaffExists`.
- `StaffCard`: Satisface `StaffReader.StaffCard`.
- `StaffByUser`: Satisface `DoctorReader.StaffByUser` de `clinical_encounter`.
