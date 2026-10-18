package config

import (
	"crypto/sha256"
	"encoding/hex"
	"io/fs"
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/etag"
	"github.com/gofiber/fiber/v3/middleware/static"
)

const assetsPrefix = "/assets"

func InitAssets(router fiber.Router, assets fs.FS) {
	assetsVersion = digestAssets(assets)
	router.Use(assetsPrefix, revalidateAssets)
	router.Use(assetsPrefix, etag.New(etag.Config{Weak: true}))
	router.Use(static.New("/assets", static.Config{FS: assets}))
}

var (
	assetsVersion  string
	versionPattern = regexp.MustCompile(`^` + assetsPrefix + `/v[0-9a-f]{12}/`)
)

func AssetURL(path string) string {
	if assetsVersion == "" {
		return assetsPrefix + "/" + path
	}
	return assetsPrefix + "/v" + assetsVersion + "/" + path
}

func digestAssets(assets fs.FS) string {
	hash := sha256.New()
	err := fs.WalkDir(assets, ".", func(path string, entry fs.DirEntry, err error) error {
		if err != nil || entry.IsDir() {
			return err
		}
		contents, err := fs.ReadFile(assets, path)
		if err != nil {
			return err
		}
		hash.Write([]byte(path))
		hash.Write(contents)
		return nil
	})
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(hash.Sum(nil))[:12]
}

func revalidateAssets(c fiber.Ctx) error {
	c.Request().Header.Del(fiber.HeaderIfModifiedSince)

	immutable := false
	if prefix := versionPattern.FindString(c.Path()); prefix != "" {
		immutable = assetsVersion != "" && prefix == assetsPrefix+"/v"+assetsVersion+"/"
		c.Path(assetsPrefix + "/" + strings.TrimPrefix(c.Path(), prefix))
	}

	if err := c.Next(); err != nil {
		return err
	}

	c.Response().Header.Del(fiber.HeaderLastModified)
	if immutable {
		c.Set(fiber.HeaderCacheControl, "public, max-age=31536000, immutable")
	} else {
		c.Set(fiber.HeaderCacheControl, "no-cache")
	}

	return nil
}
