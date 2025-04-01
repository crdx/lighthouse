package service

import (
	"fmt"

	"crdx.org/lighthouse/db"
	"crdx.org/lighthouse/db/repo/auditLogR"
	"crdx.org/lighthouse/pkg/flash"
	"crdx.org/lighthouse/pkg/globals"
	"crdx.org/lighthouse/pkg/middleware/parseparam"
	"crdx.org/lighthouse/pkg/transform"
	"crdx.org/lighthouse/pkg/validate"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
)

type EditForm struct {
	Name string `form:"name" validate:"max=100"`
}

func ViewEdit(c *fiber.Ctx) error {
	service := parseparam.Get[db.Service](c)

	return c.Render("service/edit", fiber.Map{
		"service": service,
		"device":  service.Device(),
		"fields":  validate.Fields[EditForm](),
		"globals": globals.Get(c),
	})
}

func Edit(c *fiber.Ctx) error {
	service := parseparam.Get[db.Service](c)

	form := new(EditForm)
	lo.Must0(c.BodyParser(form))
	transform.Struct(form)

	if fields, err := validate.Struct(form); err != nil {
		flash.Failure(c, "Unable to save service")

		return c.Render("service/edit", fiber.Map{
			"service": service,
			"device":  service.Device(),
			"err":     err,
			"fields":  fields,
			"globals": globals.Get(c),
		})
	}

	db.UpdateService(db.UpdateServiceParams{
		ID:   service.ID,
		Name: form.Name,
	})

	service.Reload()

	auditLogR.Add(c, "Edited service %s", service.AuditName())
	flash.Success(c, "Service saved")
	return c.Redirect(fmt.Sprintf("/device/%d", service.DeviceID))
}
