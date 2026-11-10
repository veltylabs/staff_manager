//go:build !browser
package module
import (
	"webtyp.com/appkit"
	"webtyp.com/input"
	"webtyp.com/model"
	"webtyp.com/view"
	staffmanager "github.com/veltylabs/staff_manager"
)
var server = appkit.ModuleServer{Build: build}
var browser view.ModuleBrowser
func build(d appkit.Deps) (model.ModuleNaming, error) {
	devices, err := appkit.Find[staffmanager.DeviceReader](d)
	if err != nil { return nil, err }
	registrar, err := appkit.Find[staffmanager.DeviceRegistrar](d)
	if err != nil { return nil, err }
	return staffmanager.New(d.DB, staffmanager.Deps{
		IDs: d.IDs, Publisher: d.Events, TenantID: d.TenantID,
		ValidateRUT: input.NormalizeRUT, Devices: devices, Registrar: registrar,
	})
}
