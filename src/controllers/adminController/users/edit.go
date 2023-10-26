package users

import (
	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/pkg/flash"
	"crdx.org/lighthouse/pkg/globals"
	"crdx.org/lighthouse/pkg/transform"
	"crdx.org/lighthouse/pkg/validate"
	"crdx.org/lighthouse/util/reflectutil"
	"crdx.org/lighthouse/util/stringutil"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
)

type EditForm struct {
	Username string `form:"username" validate:"required,max=20" transform:"trim"`
	Password string `form:"password" transform:"trim"`
	Admin    bool   `form:"admin"`
}

func ViewEdit(c *fiber.Ctx) error {
	user := c.Locals("user").(*m.User)

	return c.Render("admin/index", fiber.Map{
		"tab":     "users",
		"mode":    "edit",
		"user":    user,
		"fields":  validate.Fields[EditForm](),
		"globals": globals.Get(c),
	})
}

func Edit(c *fiber.Ctx) error {
	user := c.Locals("user").(*m.User)

	form := new(EditForm)
	lo.Must0(c.BodyParser(form))

	// Current user can't edit their own admin access so the disabled form field doesn't come
	// through in the request. Set it here to make the form object valid.
	if globals.IsCurrentUser(c, user) {
		form.Admin = user.Admin
	}

	transform.Struct(form)

	if fields, err := validate.Struct(form); err {
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

	password := form.Password

	values := reflectutil.StructToMap(form, "form")

	if password != "" {
		values["password_hash"] = stringutil.Hash(password)
	}
	delete(values, "password")

	// Current user can't edit their own admin access.
	if globals.IsCurrentUser(c, user) {
		delete(values, "admin")
	}

	user.Update(values)

	flash.Success(c, "User saved")
	return c.Redirect("/admin/users")
}
