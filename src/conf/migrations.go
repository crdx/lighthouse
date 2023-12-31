package conf

import (
	"crdx.org/db"
	"crdx.org/lighthouse/migrations"
)

//  GENERATED CODE — DO NOT EDIT 

var dbMigrations = []*db.Migration{
	migrations.AddDefaultScanInterval("AddDefaultScanInterval"),
	migrations.RenameLastSeenToLastSeenAt("RenameLastSeenToLastSeenAt"),
	migrations.ConvertIconClasses("ConvertIconClasses"),
	migrations.ConvertDurations("ConvertDurations"),
	migrations.MigrateToSimpleRoles("MigrateToSimpleRoles"),
	migrations.AddForeignKeyConstraints("AddForeignKeyConstraints"),
	migrations.AddCheckConstraints("AddCheckConstraints"),
	migrations.AddMappingTableConstraints("AddMappingTableConstraints"),
	migrations.AddAuditLogDeviceIdConstraint("AddAuditLogDeviceIdConstraint"),
	migrations.RenameSettings("RenameSettings"),
	migrations.AddServiceConstraints("AddServiceConstraints"),
	migrations.AddDeviceIpAddressLogConstraints("AddDeviceIpAddressLogConstraints"),
}
