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
	const modelDir = "m"
	var modelNames []string

	for _, entry := range lo.Must(os.ReadDir(modelDir)) {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".go") {
			contents := lo.Must(os.ReadFile(path.Join(modelDir, entry.Name())))

			matches := regexp.MustCompile("type (.*?) struct {").FindSubmatch(contents)

			// We can assume that the only struct defined is the model. If there are more, then
			// that's a problem.
			if len(matches) > 2 {
				panic(fmt.Errorf("too many model names matched for file %s/%s", modelDir, entry.Name()))
			}

			if len(matches) == 0 || string(matches[1]) == "" {
				panic(fmt.Errorf("no model names matched for file %s/%s", modelDir, entry.Name()))
			}

			modelNames = append(modelNames, string(matches[1]))
		}
	}

	var code strings.Builder

	code.WriteString("package conf\n\n")
	code.WriteString("import (\n")
	code.WriteString("\t\"crdx.org/db\"\n")
	code.WriteString("\t\"crdx.org/lighthouse/m\"\n")
	code.WriteString(")\n\n")
	code.WriteString("//  GENERATED CODE — DO NOT EDIT \n\n")
	code.WriteString("var models = []db.Model{\n")

	for _, model := range modelNames {
		code.WriteString(fmt.Sprintf("\t&m.%s{},\n", model))
	}

	code.WriteString("}\n")

	lo.Must0(os.WriteFile("conf/models.go", []byte(code.String()), 0o644))
}
