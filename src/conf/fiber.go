package conf

import (
	"errors"
	"io/fs"
	"net/http"
	"os"

	"crdx.org/lighthouse/pkg/env"
	"crdx.org/lighthouse/pkg/util"
	"crdx.org/lighthouse/views/helpers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

func GetFiberConfig(views fs.FS) fiber.Config {
	// The embedded views fs.FS expects all calls to Open to be prefixed with "views/", but it's
	// neater if we don't need to specify it in our calls to c.Render(), so wrap the fs.FS in an
	// instance of PrefixFS which will wrap calls to Open and add the prefix for us.
	prefixedFS := util.PrefixFS{FS: views, Prefix: "views"}
	viewsEngine := html.NewFileSystem(http.FS(&prefixedFS), ".go.html")

	viewsEngine.AddFuncMap(helpers.GetFuncMap())

	// https://docs.gofiber.io
	config := fiber.Config{
		Views:                   viewsEngine,
		ViewsLayout:             "layouts/main",
		Immutable:               true, // https://docs.gofiber.io/#zero-allocation
		EnableTrustedProxyCheck: true,
		TrustedProxies:          []string{"127.0.0.1"},
		ProxyHeader:             "X-Forwarded-For",
	}

	if env.Production() {
		config.ErrorHandler = func(c *fiber.Ctx, err error) error {
			if e := new(fiber.Error); errors.As(err, &e) {
				return c.SendStatus(e.Code)
			} else {
				return c.SendStatus(fiber.StatusInternalServerError)
			}
		}
	}

	return config
}

func GetTestFiberConfig() fiber.Config {
	views := html.New(os.Getenv("VIEWS_DIR"), ".go.html")
	views.AddFuncMap(helpers.GetFuncMap())

	// https://docs.gofiber.io
	return fiber.Config{
		Views:       views,
		ViewsLayout: "layouts/main",
		Immutable:   true, // https://docs.gofiber.io/#zero-allocation
	}
}
