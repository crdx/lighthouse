package profileController

import (
	"crdx.org/lighthouse/m/repo/auditLogR"
	"crdx.org/lighthouse/pkg/flash"
	"crdx.org/lighthouse/pkg/globals"
	"crdx.org/lighthouse/pkg/transform"
	"crdx.org/lighthouse/pkg/validate"
	"crdx.org/lighthouse/util/stringutil"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
)

type EditForm struct {
	CurrentPassword    string `form:"current_password" validate:"required" transform:"trim"`
	NewPassword        string `form:"new_password" validate:"required,min=4" transform:"trim"`
	ConfirmNewPassword string `form:"confirm_new_password" validate:"required" transform:"trim"`
}

func Edit(c *fiber.Ctx) error {
	form := new(EditForm)
	lo.Must0(c.BodyParser(form))
	transform.Struct(form)

	user := globals.Get(c).User

	validatorMap := validate.ValidatorMap{
		"CurrentPassword":    validate.CurrentPassword(user.PasswordHash),
		"ConfirmNewPassword": validate.ConfirmPassword(form.NewPassword),
	}

	if fields, err := validate.Struct(form, validatorMap); err != nil {
		flash.Failure(c, "Unable to change password")

		return c.Render("profile/view", fiber.Map{
			"fields":  fields,
			"err":     err,
			"globals": globals.Get(c),
		})
	}

	user.Update("password_hash", stringutil.Hash(form.NewPassword))

	auditLogR.Add(c, "Changed password")
	flash.Success(c, "Password changed")
	return c.Redirect("/profile")
}
