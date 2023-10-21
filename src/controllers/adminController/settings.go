package adminController

import (
	"crdx.org/lighthouse/m/repo/settingR"
	"crdx.org/lighthouse/pkg/flash"
	"crdx.org/lighthouse/pkg/globals"
	"crdx.org/lighthouse/pkg/transform"
	"crdx.org/lighthouse/pkg/validate"
	"crdx.org/lighthouse/util/reflectutil"
	"github.com/gofiber/fiber/v2"
)

type SettingsForm struct {
	Watch            bool   `form:"watch"`
	MACVendorsAPIKey string `form:"macvendors_api_key" validate:"max=500" transform:"trim"`
}

func Index(c *fiber.Ctx) error {
	return c.Render("admin/settings", fiber.Map{
		"fields":   validate.Fields[SettingsForm](),
		"settings": reflectutil.MapToStruct[SettingsForm](settingR.Map(), "form"),
		"globals":  globals.Get(c),
	})
}

func Save(c *fiber.Ctx) error {
	form := new(SettingsForm)
	if err := c.BodyParser(form); err != nil {
		return err
	}

	transform.Struct(form)

	if fields, err := validate.Struct(form); err {
		flash.Failure(c, "Unable to save settings")

		return c.Render("admin/settings", fiber.Map{
			"err":     err,
			"fields":  fields,
			"globals": globals.Get(c),
		})
	}

	for name, value := range reflectutil.StructToMap(form, "form") {
		settingR.Set(name, value)
	}

	flash.Success(c, "Settings saved")
	return c.Redirect("/admin")
}
