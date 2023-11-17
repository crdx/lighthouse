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
}
