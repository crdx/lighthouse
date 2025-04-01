package adapter

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
	Name   string `form:"name" validate:"max=100"`
	Vendor string `form:"vendor" validate:"max=100"`
}

func ViewEdit(c *fiber.Ctx) error {
	adapter := parseparam.Get[db.Adapter](c)

	return c.Render("adapter/edit", fiber.Map{
		"adapter": adapter,
		"device":  adapter.Device(),
		"fields":  validate.Fields[EditForm](),
		"globals": globals.Get(c),
	})
}

func Edit(c *fiber.Ctx) error {
	adapter := parseparam.Get[db.Adapter](c)

	form := new(EditForm)
	lo.Must0(c.BodyParser(form))
	transform.Struct(form)

	if fields, err := validate.Struct(form); err != nil {
		flash.Failure(c, "Unable to save adapter")

		return c.Render("adapter/edit", fiber.Map{
			"adapter": adapter,
			"device":  adapter.Device(),
			"err":     err,
			"fields":  fields,
			"globals": globals.Get(c),
		})
	}

	db.UpdateAdapter(db.UpdateAdapterParams{
		ID:     adapter.ID,
		Name:   form.Name,
		Vendor: form.Vendor,
	})

	adapter.Reload()

	auditLogR.Add(c, "Edited adapter %s", adapter.AuditName())
	flash.Success(c, "Adapter saved")
	return c.Redirect(fmt.Sprintf("/device/%d", adapter.DeviceID))
}
