# Pruebas

| Test | Regla |
|---|---|
| [TestStaff_UpdateOnlyListedFields](../tests/test_staff_update_test.go) | actualizar `["phone"]` no cambia el nombre |
| [TestStaff_DuplicateRutRejected](../tests/test_staff_duplicate_test.go) | dos personas con el mismo RUT: la segunda se rechaza |
| [TestStaffDevice_NotOffered](../tests/test_staff_device_test.go) | `staff_device` no se ofrece como operación: `unknown operation` |
| [TestTenant_ReadOtherTenantIsNotFound](../tests/test_tenant_read_test.go) | una persona del inquilino a leída por el inquilino b es `model.ErrorNotFound` |
| [TestReaders_RegisterFirstAdmin](../tests/test_readers_test.go) | `RegisterFirstAdmin` crea al administrador con las mismas reglas de RUT |
| [TestTrust_IsTrustedIPDecisions](../tests/trust_test.go) | `IsTrustedIP` responde solo con el equipo asignado a una persona activa |
| [TestReaders_StaffByUser](../tests/trust_test.go) | `StaffByUser` encuentra al profesional de una cuenta y responde `found=false` si no hay |
| [TestReaders_StaffCardAndExists](../tests/trust_test.go) | `StaffCard` y `StaffExists` respetan el inquilino y no devuelven datos personales |
| [TestReaders_ListStaffSearchesByName](../tests/trust_test.go) | `ListStaff` busca por nombre |
| [TestMount_MatchesDefinition](../tests/mount_test.go) | lo que monta `MountOperations` coincide con `Definition()` |
| [TestCreate_PublishesChange](../tests/mount_test.go) | crear una persona publica un `*crud.Change` en `crud.Topic(ModelName, "staff_member", model.Create)` |
| [TestSeed_Load](../tests/seed_test.go) | `seed.Load` carga las personas de la demo |
| [TestUI_Screens](../tests/ui_view_test.go) | la pantalla del personal se construye |
| [TestUI_ViewReloadReadsStaff](../tests/ui_view_test.go) | recargar la vista pide `staff_member.read` |
| [TestDemo_Valid](../tests/demo_test.go) | `seed.Demo()` (wasm) es válida y entrega al menos una pantalla |
