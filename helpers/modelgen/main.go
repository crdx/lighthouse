package main

import (
	"fmt"
	"io/ioutil"
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

			match := regexp.MustCompile("type (.*?) struct {").FindSubmatch(contents)

			// We can assume that the only struct defined is the model. If there are more, then
			// that's a problem.
			if len(match) > 2 {
				panic(fmt.Errorf("too many model names matched for file %s/%s", modelDir, entry.Name()))
			}

			if len(match) == 0 || string(match[1]) == "" {
				panic(fmt.Errorf("no model names matched for file %s/%s", modelDir, entry.Name()))
			}

			modelNames = append(modelNames, string(match[1]))
		}
	}

	var code strings.Builder

	code.WriteString("package conf\n\n")
	code.WriteString("import \"crdx.org/lighthouse/m\"\n\n")
	code.WriteString("//  GENERATED CODE — DO NOT EDIT \n\n")
	code.WriteString("var models = []any{\n")

	for _, model := range modelNames {
		code.WriteString(fmt.Sprintf("\t&m.%s{},\n", model))
	}

	code.WriteString("}\n")

	lo.Must0(ioutil.WriteFile("conf/models.go", []byte(code.String()), 0644))
}
