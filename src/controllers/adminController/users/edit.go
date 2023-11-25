package users

import (
	"fmt"

	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/m/repo/auditLogR"
	"crdx.org/lighthouse/middleware/util"
	"crdx.org/lighthouse/pkg/flash"
	"crdx.org/lighthouse/pkg/globals"
	"crdx.org/lighthouse/pkg/transform"
	"crdx.org/lighthouse/pkg/util/reflectutil"
	"crdx.org/lighthouse/pkg/validate"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
)

type EditForm struct {
	Password        string `form:"password" validate:"omitempty,min=4" transform:"trim"`
	ConfirmPassword string `form:"confirm_password" transform:"trim"`
	Role            string `form:"role" validate:"required,role"`
}

func ViewEdit(c *fiber.Ctx) error {
	user := util.Param[m.User](c)

	return c.Render("admin/index", fiber.Map{
		"tab":     "users",
		"mode":    "edit",
		"user":    user,
		"fields":  validate.Fields[EditForm](),
		"globals": globals.Get(c),
	})
}

func Edit(c *fiber.Ctx) error {
	user := util.Param[m.User](c)

	form := new(EditForm)
	lo.Must0(c.BodyParser(form))

	// Current user can't edit their own admin access so the disabled form field doesn't come
	// through in the request. Set it here to make the form object valid.
	if globals.IsCurrentUser(c, user) {
		form.Role = fmt.Sprint(user.Role)
	}

	transform.Struct(form)

	validatorMap := validate.ValidatorMap{
		"ConfirmPassword": validate.ConfirmPassword(form.Password),
	}

	if fields, err := validate.Struct(form, validatorMap); err != nil {
		flash.Failure(c, "Unable to save user")

		return c.Render("admin/index", fiber.Map{
			"tab":     "users",
			"mode":    "edit",
			"user":    user,
			"err":     err,
			"fields":  fields,
			"globals": globals.Get(c),
		})
	}

	values := reflectutil.StructToMap(form, "form")

	transform.PasswordFields(values)

	// Admins can't demote themselves.
	if globals.IsCurrentUser(c, user) {
		delete(values, "role")
	}

	user.Update(values)

	auditLogR.Add(c, "Edited user %s", user.Fresh().AuditName())
	flash.Success(c, "User saved")
	return c.Redirect("/admin/users")
}
