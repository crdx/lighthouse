package users

import (
	"errors"

	"crdx.org/db"
	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/m/repo/auditLogR"
	"crdx.org/lighthouse/pkg/flash"
	"crdx.org/lighthouse/pkg/globals"
	"crdx.org/lighthouse/pkg/transform"
	"crdx.org/lighthouse/pkg/util/reflectutil"
	"crdx.org/lighthouse/pkg/validate"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
)

type CreateForm struct {
	Username        string `form:"username" validate:"required,max=20" transform:"trim"`
	Password        string `form:"password" validate:"required,min=4" transform:"trim"`
	ConfirmPassword string `form:"confirm_password" validate:"required" transform:"trim"`
	Role            string `form:"role" validate:"required,role"`
}

func ViewCreate(c *fiber.Ctx) error {
	return c.Render("admin/index", fiber.Map{
		"tab":     "users",
		"mode":    "create",
		"fields":  validate.Fields[CreateForm](),
		"globals": globals.Get(c),
	})
}

func Create(c *fiber.Ctx) error {
	form := new(CreateForm)
	lo.Must0(c.BodyParser(form))
	transform.Struct(form)

	validatorMap := validate.ValidatorMap{
		"Username": func(value string) error {
			for _, user := range db.B[m.User]().Find() {
				if user.Username == value {
					return errors.New("must be an available username")
				}
			}
			return nil
		},
		"ConfirmPassword": validate.ConfirmPassword(form.Password),
	}

	if fields, err := validate.Struct(form, validatorMap); err != nil {
		flash.Failure(c, "Unable to create user")

		return c.Render("admin/index", fiber.Map{
			"tab":     "users",
			"mode":    "create",
			"err":     err,
			"fields":  fields,
			"globals": globals.Get(c),
		})
	}

	values := reflectutil.StructToMap(form, "form")
	transform.PasswordFields(values)

	user := db.CreateFromMap[m.User](values)

	auditLogR.Add(c, "Created user %s", user.AuditName())
	flash.Success(c, "User created")
	return c.Redirect("/admin/users")
}
