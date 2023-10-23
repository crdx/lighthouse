package settings

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
	// General
	Watch bool `form:"watch"`

	// Mail
	EnableMail      bool   `form:"enable_mail"`
	MailFromHeader  string `form:"mail_from_header"  transform:"trim" validate:"required_if=EnableMail 1,omitempty,mailaddr"`
	MailFromAddress string `form:"mail_from_address" transform:"trim" validate:"required_if=EnableMail 1,omitempty,email"`
	MailToHeader    string `form:"mail_to_header"    transform:"trim" validate:"required_if=EnableMail 1,omitempty,mailaddr"`
	MailToAddress   string `form:"mail_to_address"   transform:"trim" validate:"required_if=EnableMail 1,omitempty,email"`
	SMTPHost        string `form:"smtp_host" transform:"trim" validate:"required_if=EnableMail 1"`
	SMTPPort        string `form:"smtp_port" transform:"trim" validate:"required_if=EnableMail 1,omitempty,number,excludesall=-."`
	SMTPUser        string `form:"smtp_user" transform:"trim" validate:"required_if=EnableMail 1"`
	SMTPPass        string `form:"smtp_pass" transform:"trim" validate:"required_if=EnableMail 1"`

	// System
	MACVendorsAPIKey string `form:"macvendors_api_key" transform:"trim" validate:"max=500"`
	Timezone         string `form:"timezone"           transform:"trim" validate:"required,timezone"`
}

func List(c *fiber.Ctx) error {
	return c.Render("admin/index", fiber.Map{
		"tab":      "settings",
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

		return c.Render("admin/index", fiber.Map{
			"tab":     "settings",
			"err":     err,
			"fields":  fields,
			"globals": globals.Get(c),
		})
	}

	for name, value := range reflectutil.StructToMap(form, "form") {
		settingR.Set(name, value)
	}

	flash.Success(c, "Settings saved")
	return c.Redirect("/admin/settings")
}
