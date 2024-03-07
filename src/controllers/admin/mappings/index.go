package mappings

import (
	"errors"

	"crdx.org/db"
	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/m/repo/auditLogR"
	"crdx.org/lighthouse/m/repo/settingR"
	"crdx.org/lighthouse/middleware/util"
	"crdx.org/lighthouse/pkg/flash"
	"crdx.org/lighthouse/pkg/globals"
	"crdx.org/lighthouse/pkg/transform"
	"crdx.org/lighthouse/pkg/util/reflectutil"
	"crdx.org/lighthouse/pkg/validate"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
)

type SourceForm struct {
	MACAddresses string `form:"source_mac_addresses"  validate:"omitempty,mac_address_list" transform:"trim,upper"`
}

type MappingForm struct {
	Label      string `form:"label" validate:"max=20" transform:"trim"`
	MACAddress string `form:"mac_address" validate:"required,mac_address" transform:"trim,upper"`
	IPAddress  string `form:"ip_address" validate:"required,ip_address" transform:"trim"`
}

func View(c *fiber.Ctx) error {
	return c.Render("admin/index", fiber.Map{
		"tab": "mappings",
		"source": fiber.Map{
			"fields": validate.Fields[SourceForm](),
			"values": sources(),
		},
		"mapping": fiber.Map{
			"fields": validate.Fields[MappingForm](),
			"values": db.B[m.Mapping]().Find(),
		},
		"globals": globals.Get(c),
	})
}

func AddMapping(c *fiber.Ctx) error {
	form := new(MappingForm)
	lo.Must0(c.BodyParser(form))
	transform.Struct(form)

	mappings := db.B[m.Mapping]().Find()

	validators := validate.ValidatorMap{
		"IPAddress": func(value string) error {
			for _, mapping := range mappings {
				if mapping.IPAddress == value {
					return errors.New("must be a unique IP address")
				}
			}
			return nil
		},
	}

	if fields, err := validate.Struct(form, validators); err != nil {
		flash.Failure(c, "Unable to save mapping")

		return c.Render("admin/index", fiber.Map{
			"tab": "mappings",
			"source": fiber.Map{
				"fields": validate.Fields[SourceForm](),
				"values": sources(),
			},
			"mapping": fiber.Map{
				"err":    err,
				"fields": fields,
				"values": mappings,
			},
			"globals": globals.Get(c),
		})
	}

	mapping := db.CreateFromMap[m.Mapping](reflectutil.StructToMap(form, "form"))

	auditLogR.Add(c, "Added mapping %s", mapping.AuditName())
	flash.Success(c, "Mapping added")
	return c.Redirect("/admin/mappings")
}

func DeleteMapping(c *fiber.Ctx) error {
	mapping := util.Param[m.Mapping](c)
	mapping.Delete()

	auditLogR.Add(c, "Deleted mapping %s", mapping.AuditName())
	flash.Success(c, "Mapping deleted")
	return c.Redirect("/admin/mappings")
}

func EditSources(c *fiber.Ctx) error {
	form := new(SourceForm)
	lo.Must0(c.BodyParser(form))
	transform.Struct(form)

	if fields, err := validate.Struct(form); err != nil {
		flash.Failure(c, "Unable to save source MAC addresses")

		return c.Render("admin/index", fiber.Map{
			"tab": "mappings",
			"source": fiber.Map{
				"err":    err,
				"fields": fields,
				"values": sources(),
			},
			"mapping": fiber.Map{
				"fields": fields,
				"values": db.B[m.Mapping]().Find(),
			},
			"globals": globals.Get(c),
		})
	}

	for name, value := range reflectutil.StructToMap(form, "form") {
		settingR.Set(name, value)
	}

	auditLogR.Add(c, "Saved source MAC addresses")
	flash.Success(c, "Source MAC addresses saved")
	return c.Redirect("/admin/mappings")
}

func sources() SourceForm {
	return reflectutil.MapToStruct[SourceForm](map[string]string{
		"source_mac_addresses": settingR.SourceMACAddresses(),
	}, "form")
}
