package main

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/samber/lo"
)

func main() {
	const migrationDir = "migrations"
	var migrationNames []string

	for _, entry := range lo.Must(os.ReadDir(migrationDir)) {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".go") {
			contents := lo.Must(os.ReadFile(path.Join(migrationDir, entry.Name())))

			matches := regexp.MustCompile(`func (.*?)\(id string\) \*db\.Migration {`).FindSubmatch(contents)

			// We can assume that the only function defined is the migration. If there are more,
			// then that's a problem.
			if len(matches) > 2 {
				panic(fmt.Errorf("too many migrations matched for file %s/%s", migrationDir, entry.Name()))
			}

			if len(matches) == 0 || string(matches[1]) == "" {
				panic(fmt.Errorf("no migrations matched for file %s/%s", migrationDir, entry.Name()))
			}

			migrationNames = append(migrationNames, string(matches[1]))
		}
	}

	var code strings.Builder

	code.WriteString("package conf\n\n")
	code.WriteString("import (\n")
	code.WriteString("\t\"crdx.org/db\"\n")
	code.WriteString("\t\"crdx.org/lighthouse/migrations\"\n")
	code.WriteString(")\n\n")
	code.WriteString("//  GENERATED CODE — DO NOT EDIT \n\n")
	code.WriteString("var dbMigrations = []*db.Migration{\n")

	for _, name := range migrationNames {
		code.WriteString(fmt.Sprintf("\tmigrations.%s(\"%s\"),\n", name, name))
	}

	code.WriteString("}\n")

	lo.Must0(os.WriteFile("conf/migrations.go", []byte(code.String()), 0o644))
}
