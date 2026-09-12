package migrate

import (
	"webtyp.com/ddl"

	staffmanager "github.com/veltylabs/staff_manager"
)

func Migrate(conn ddl.Execer, ddlCompiler ddl.Compiler) error {
	d := ddl.New(conn, ddlCompiler)
	if err := d.CreateTable(&staffmanager.StaffMember{}); err != nil {
		return err
	}
	if err := d.CreateTable(&staffmanager.StaffDevice{}); err != nil {
		return err
	}
	return nil
}
