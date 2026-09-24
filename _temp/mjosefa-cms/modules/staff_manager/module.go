package staff_manager

// ID is this module's RBAC resource — the same string its own ops declare
// via .Requires("staff_manager", ...) and the domain library's own
// ModelName(). No Label: this module has no Browser (see server.go) — its
// screen lives inside "personal" instead — so it has an identity to grant
// permission against, but nothing of its own to put in a nav item.
const ID = "staff_manager"
