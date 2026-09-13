package migrate

import (
	"webtyp.com/ddl"

	staffmanager "github.com/veltylabs/staff_manager"
)

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
