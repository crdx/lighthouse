package settings

import (
	"crdx.org/lighthouse/m/repo/auditLogR"
	"crdx.org/lighthouse/m/repo/settingR"
	"crdx.org/lighthouse/pkg/flash"
	"crdx.org/lighthouse/pkg/globals"
	"crdx.org/lighthouse/pkg/transform"
	"crdx.org/lighthouse/pkg/util/reflectutil"
	"crdx.org/lighthouse/pkg/validate"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
)

type Form struct {
	// General
	Watch    bool `form:"watch"`
	WatchNew bool `form:"watch_new"`

	// Mail
	EnableMail      bool   `form:"enable_mail"`
	MailFromHeader  string `form:"mail_from_header"  validate:"required_if=EnableMail 1,omitempty,mailaddr" transform:"trim"`
	MailFromAddress string `form:"mail_from_address" validate:"required_if=EnableMail 1,omitempty,email" transform:"trim"`
	MailToHeader    string `form:"mail_to_header"    validate:"required_if=EnableMail 1,omitempty,mailaddr" transform:"trim"`
	MailToAddress   string `form:"mail_to_address"   validate:"required_if=EnableMail 1,omitempty,email" transform:"trim"`
	SMTPHost        string `form:"smtp_host"         validate:"required_if=EnableMail 1" transform:"trim"`
	SMTPPort        string `form:"smtp_port"         validate:"required_if=EnableMail 1,omitempty,number,excludesall=-." transform:"trim"`
	SMTPUser        string `form:"smtp_user"         validate:"required_if=EnableMail 1" transform:"trim"`
	SMTPPass        string `form:"smtp_pass"         validate:"required_if=EnableMail 1" transform:"trim"`

	// System
	MACVendorsAPIKey string `form:"macvendors_api_key" validate:"max=500" transform:"trim"`
	Timezone         string `form:"timezone"           validate:"required,timezone" transform:"trim"`

	// Scanning
	Passive      bool   `form:"passive"`
	ScanInterval string `form:"scan_interval" validate:"required,duration,dmin=1 min,dmax=30 mins" transform:"trim"`
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

	auditLogR.Add(c, "Saved settings")
	flash.Success(c, "Settings saved")
	return c.Redirect("/admin/settings")
}
