package config

import (
	"path"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
)

var compressedExtensions = map[string]bool{
	".woff2": true,
	".woff":  true,
	".png":   true,
	".jpg":   true,
	".jpeg":  true,
	".gif":   true,
	".webp":  true,
	".avif":  true,
}

func NewCompress() fiber.Handler {
	return compress.New(compress.Config{
		Next: func(c fiber.Ctx) bool {
			return compressedExtensions[path.Ext(c.Path())]
		},
	})
}
