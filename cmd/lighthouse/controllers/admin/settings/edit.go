package settings

import (
	"crdx.org/lighthouse/db/repo/auditLogR"
	"crdx.org/lighthouse/db/repo/settingR"
	"crdx.org/lighthouse/pkg/flash"
	"crdx.org/lighthouse/pkg/globals"
	"crdx.org/lighthouse/pkg/transform"
	"crdx.org/lighthouse/pkg/util/reflectutil"
	"crdx.org/lighthouse/pkg/validate"
	"github.com/gofiber/fiber/v3"
	"github.com/samber/lo"
)

type EditForm struct {
	// Device Scanning
	EnableDeviceScan   bool   `form:"enable_device_scan"`
	DeviceScanInterval string `form:"device_scan_interval" validate:"required,duration,dmin=1 min,dmax=30 mins"`
	NotifyOnNewDevice  bool   `form:"notify_on_new_device"`
	WatchNew           bool   `form:"watch_new"`
	PingNew            bool   `form:"ping_new"`

	// Service Scanning
	EnableServiceScan   bool   `form:"enable_service_scan"`
	ServiceScanInterval string `form:"service_scan_interval" validate:"required,duration,dmin=1 hour,dmax=1 week"`
	NotifyOnNewService  bool   `form:"notify_on_new_service"`

	// Mail
	EnableMail      bool   `form:"enable_mail"`
	MailFromHeader  string `form:"mail_from_header"  validate:"required_if=EnableMail true,omitempty,mailaddr"`
	MailFromAddress string `form:"mail_from_address" validate:"required_if=EnableMail true,omitempty,email"`
	MailToHeader    string `form:"mail_to_header"    validate:"required_if=EnableMail true,omitempty,mailaddr"`
	MailToAddress   string `form:"mail_to_address"   validate:"required_if=EnableMail true,omitempty,email"`
	SMTPHost        string `form:"smtp_host"         validate:"required_if=EnableMail true"`
	SMTPPort        string `form:"smtp_port"         validate:"required_if=EnableMail true,omitempty,number,excludesall=-."`
	SMTPUser        string `form:"smtp_user"         validate:"required_if=EnableMail true"`
	SMTPPass        string `form:"smtp_pass"         validate:"required_if=EnableMail true"`

	// System
	MACVendorsAPIKey string `form:"macvendors_api_key" validate:"max=500"`
	Timezone         string `form:"timezone"           validate:"required,timezone"`
}

func List(c fiber.Ctx) error {
	return c.Render("admin/index", fiber.Map{
		"tab":      "settings",
		"fields":   validate.Fields[EditForm](),
		"settings": reflectutil.MapToStruct[EditForm](settingR.Map(), "form"),
		"globals":  globals.Get(c),
	})
}

func Save(c fiber.Ctx) error {
	form := new(EditForm)
	lo.Must0(c.Bind().Body(form))
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
	return c.Redirect().To("/admin/settings")
}
