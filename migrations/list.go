package migrations

import (
	"embed"
	"strings"

	"crdx.org/lighthouse/db"
)

//go:embed *.sql
var fs embed.FS

func List() []*db.Migration {
	files, err := fs.ReadDir(".")
	if err != nil {
		panic(err)
	}

	var migrations []*db.Migration
	for _, file := range files {
		sql, err := fs.ReadFile(file.Name())
		if err != nil {
			panic(err)
		}

		migrations = append(migrations, &db.Migration{
			Name: strings.TrimSuffix(file.Name(), ".sql"),
			SQL:  string(sql),
		})
	}

	return migrations
}
