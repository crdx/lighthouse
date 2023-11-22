package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"
	"time"
	"unicode"

	"github.com/samber/lo"
)

func main() {
	log.SetFlags(0)
	if len(os.Args) != 2 {
		log.Fatal("Usage: mkmigration <name>")
	}

	name := os.Args[1]
	if !isProperCase(name) {
		log.Fatal("\033[31mname must be ProperCase without any punctuation, beginning with a letter\033[m")
	}

	fileName := fmt.Sprintf("%d_%s.go", time.Now().Unix(), snakeCase(name))
	filePath := "src/migrations/" + fileName
	tpl := lo.Must(template.New("migration.template").ParseFiles("src/cmd/mkmigration/migration.template"))
	file := lo.Must(os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0o644))

	if err := tpl.Execute(file, name); err != nil {
		log.Fatalf("\033[31m✗ %s: %s\033[m\n", fileName, err)
	} else {
		fmt.Printf("\033[32m✓ %s\033[m\n", filePath)
	}
}

func isProperCase(s string) bool {
	if s == "" || !unicode.IsLetter(rune(s[0])) || !unicode.IsUpper(rune(s[0])) {
		return false
	}

	var prevIsUpper bool
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return false
		}

		if unicode.IsUpper(r) {
			if prevIsUpper {
				return false
			}
			prevIsUpper = true
		} else {
			prevIsUpper = false
		}
	}

	return true
}

func snakeCase(str string) string {
	var buf bytes.Buffer
	for i, rune := range str {
		if unicode.IsUpper(rune) {
			if i > 0 {
				buf.WriteByte('_')
			}
			buf.WriteRune(unicode.ToLower(rune))
		} else {
			buf.WriteRune(rune)
		}
	}
	return buf.String()
}
