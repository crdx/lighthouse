package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
	"unicode"
)

func main() {
	log.SetFlags(0)
	if len(os.Args) != 2 {
		os.Exit(1)
	}
	name := os.Args[1]

	fileName := fmt.Sprintf("%d_%s", time.Now().UTC().Unix(), snakeCase(name)) + ".sql"
	filePath := filepath.Join("migrations", fileName)

	_, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0o644)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("\033[32m✓ %s\033[m", fileName)
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
