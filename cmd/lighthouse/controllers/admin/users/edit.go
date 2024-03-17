package users

import (
	"fmt"
	"strconv"

	"crdx.org/lighthouse/db"
	"crdx.org/lighthouse/db/repo/auditLogR"
	"crdx.org/lighthouse/pkg/flash"
	"crdx.org/lighthouse/pkg/globals"
	"crdx.org/lighthouse/pkg/middleware/parseparam"
	"crdx.org/lighthouse/pkg/transform"
	"crdx.org/lighthouse/pkg/util/stringutil"
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
	user := parseparam.Get[db.User](c)

	return c.Render("admin/index", fiber.Map{
		"tab":     "users",
		"mode":    "edit",
		"user":    user,
		"fields":  validate.Fields[EditForm](),
		"globals": globals.Get(c),
	})
}

func Edit(c *fiber.Ctx) error {
	user := parseparam.Get[db.User](c)

	form := new(EditForm)
	lo.Must0(c.BodyParser(form))

	// Current user can't edit their own admin access so the disabled form field doesn't come
	// through in the request. Set it here to make the form object valid.
	if globals.IsCurrentUser(c, user.ID) {
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

	params := db.UpdateUserParams{
		ID:           user.ID,
		PasswordHash: user.PasswordHash,
		Role:         user.Role,
	}

	if !globals.IsCurrentUser(c, user.ID) {
		params.Role = lo.Must(strconv.ParseInt(form.Role, 10, 64))
	}

	if form.Password != "" {
		params.PasswordHash = stringutil.Hash(form.Password)
	}

	db.UpdateUser(params)
	user.Reload()

	auditLogR.Add(c, "Edited user %s", user.AuditName())
	flash.Success(c, "User saved")
	return c.Redirect("/admin/users")
}
