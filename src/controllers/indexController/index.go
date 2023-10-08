package indexController

import (
	"fmt"
	"slices"

	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/pkg/flash"
	"crdx.org/lighthouse/pkg/validate"
	"crdx.org/lighthouse/repos/deviceR"
	"crdx.org/lighthouse/tpl"
	"crdx.org/lighthouse/util/reflectutil"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"golang.org/x/exp/maps"
)

func Get(c *fiber.Ctx) error {
	columns := map[string]tpl.SortableColumnConfig{
		"name":   {Label: "Name", DefaultSortDirection: "asc"},
		"ip":     {Label: "IP Address", DefaultSortDirection: "asc"},
		"vendor": {Label: "Vendor", DefaultSortDirection: "asc"},
		"mac":    {Label: "MAC Address", DefaultSortDirection: "asc"},
		"seen":   {Label: "Last Seen", DefaultSortDirection: "desc"},
	}

	currentSortColumn := c.Query("sc", "seen")
	currentSortDirection := c.Query("sd", "desc")

	if !slices.Contains([]string{"asc", "desc"}, currentSortDirection) {
		return c.SendStatus(400)
	}

	if !slices.Contains(maps.Keys(columns), currentSortColumn) {
		return c.SendStatus(400)
	}

	return c.Render("index", fiber.Map{
		"devices": deviceR.GetListView(currentSortColumn, currentSortDirection),
		"columns": tpl.AddSortMetadata(currentSortColumn, currentSortDirection, columns),
		flash.Key: c.Locals(flash.Key),
	})
}

func ViewDevice(c *fiber.Ctx) error {
	device, found := getDevice(c)
	if !found {
		return c.SendStatus(404)
	}

	return c.Render("view", fiber.Map{
		"device":   device,
		"adapters": device.Adapters(),
		flash.Key:  c.Locals(flash.Key),
	})
}

func DeleteDevice(c *fiber.Ctx) error {
	device, found := getDevice(c)
	if !found {
		return c.SendStatus(404)
	}

	device.Delete()
	flash.AddSuccess(c, "Device deleted")
	return c.Redirect("/")
}

func ViewEditDevice(c *fiber.Ctx) error {
	device, found := getDevice(c)
	if !found {
		return c.SendStatus(404)
	}

	return c.Render("edit", fiber.Map{
		"device":  device,
		flash.Key: c.Locals(flash.Key),
	})
}

func EditDevice(c *fiber.Ctx) error {
	device, found := getDevice(c)
	if !found {
		return c.SendStatus(404)
	}

	type Form struct {
		Name        string `form:"name" validate:"required,max=100"`
		Icon        string `form:"icon" validate:"required,max=100"`
		Notes       string `form:"notes" validate:"max=4096"`
		GracePeriod uint   `form:"grace_period" validate:"required,gte=1,lte=60"`
	}

	form := new(Form)
	if err := c.BodyParser(form); err != nil {
		return err
	}

	if fields, err := validate.Struct(form); err {
		return c.Render("edit", fiber.Map{
			"device":  device,
			"err":     err,
			"fields":  fields,
			flash.Key: flash.GetFailure("Unable to save device"),
		})
	}

	device.Update(reflectutil.StructToMap(form, "form"))

	flash.AddSuccess(c, "Device saved")
	return c.Redirect(fmt.Sprintf("/device/%d", device.ID))
}

func getDevice(c *fiber.Ctx) (*m.Device, bool) {
	return m.ForDevice(uint(lo.Must(c.ParamsInt("id")))).First()
}
