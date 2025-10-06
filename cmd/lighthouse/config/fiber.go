package config

import (
	"errors"
	"io/fs"
	"net/http"
	"os"
	"strings"

	"crdx.org/lighthouse/pkg/env"
	"crdx.org/lighthouse/pkg/util/stringutil"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/template/html/v2"
	"github.com/samber/lo"
)

func GetFiberConfig(views fs.FS, subdir string) fiber.Config {
	viewsEngine := html.NewFileSystem(http.FS(lo.Must(fs.Sub(views, subdir))), ".go.html")

	viewsEngine.AddFuncMap(GetViewFuncMap())

	// https://docs.gofiber.io
	config := fiber.Config{
		Views:       viewsEngine,
		ViewsLayout: "layout/main",
		Immutable:   false, // https://docs.gofiber.io/#zero-allocation
	}

	if env.TrustedProxies() != "" {
		config.TrustProxy = true
		config.TrustProxyConfig = fiber.TrustProxyConfig{
			Proxies: lo.Map(strings.Split(env.TrustedProxies(), ","), stringutil.MapTrimSpace),
		}
		config.ProxyHeader = "X-Forwarded-For"
	}

	if env.Production() {
		config.ErrorHandler = func(c fiber.Ctx, err error) error {
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
	views.AddFuncMap(GetViewFuncMap())

	// https://docs.gofiber.io
	return fiber.Config{
		Views:       views,
		ViewsLayout: "layout/main",
		Immutable:   false, // https://docs.gofiber.io/#zero-allocation
	}
}
