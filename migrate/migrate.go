package migrate

import (
	"webtyp.com/ddl"

	staffmanager "github.com/veltylabs/staff_manager"
)

func Migrate(conn ddl.Execer, ddlCompiler ddl.Compiler) error {
	// Sync, not CreateTable: CreateTable compiles to CREATE TABLE IF NOT
	// EXISTS and is a no-op against a table that already exists, so a column
	// added to a model would never reach a deployed database. Sync creates the
	// table when it is absent and adds the missing columns when it is not.
	return ddl.New(conn, ddlCompiler).Sync(
		&staffmanager.StaffMember{},
		&staffmanager.StaffDevice{},
	)
}
