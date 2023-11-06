package settings

import (
	"crdx.org/lighthouse/m/repo/settingR"
	"crdx.org/lighthouse/pkg/flash"
	"crdx.org/lighthouse/pkg/globals"
	"crdx.org/lighthouse/pkg/transform"
	"crdx.org/lighthouse/pkg/validate"
	"crdx.org/lighthouse/util/reflectutil"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
)

type Form struct {
	// General
	Watch    bool `form:"watch"`
	WatchNew bool `form:"watch_new"`

	// Mail
	EnableMail      bool   `form:"enable_mail"`
	MailFromHeader  string `form:"mail_from_header"  transform:"trim" validate:"required_if=EnableMail 1,omitempty,mailaddr"`
	MailFromAddress string `form:"mail_from_address" transform:"trim" validate:"required_if=EnableMail 1,omitempty,email"`
	MailToHeader    string `form:"mail_to_header"    transform:"trim" validate:"required_if=EnableMail 1,omitempty,mailaddr"`
	MailToAddress   string `form:"mail_to_address"   transform:"trim" validate:"required_if=EnableMail 1,omitempty,email"`
	SMTPHost        string `form:"smtp_host"         transform:"trim" validate:"required_if=EnableMail 1"`
	SMTPPort        string `form:"smtp_port"         transform:"trim" validate:"required_if=EnableMail 1,omitempty,number,excludesall=-."`
	SMTPUser        string `form:"smtp_user"         transform:"trim" validate:"required_if=EnableMail 1"`
	SMTPPass        string `form:"smtp_pass"         transform:"trim" validate:"required_if=EnableMail 1"`

	// System
	MACVendorsAPIKey string `form:"macvendors_api_key" transform:"trim" validate:"max=500"`
	Timezone         string `form:"timezone"           transform:"trim" validate:"required,timezone"`

	// Scanning
	Passive      bool   `form:"passive"`
	ScanInterval string `form:"scan_interval" transform:"trim" validate:"required,duration=1 min:30 mins"`
}

func List(c *fiber.Ctx) error {
	return c.Render("admin/index", fiber.Map{
		"tab":      "settings",
		"fields":   validate.Fields[Form](),
		"settings": reflectutil.MapToStruct[Form](settingR.Map(), "form"),
		"globals":  globals.Get(c),
	})
}

func Save(c *fiber.Ctx) error {
	form := new(Form)
	lo.Must0(c.BodyParser(form))

	transform.Struct(form)

	if fields, err := validate.Struct(form); err != nil {
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
