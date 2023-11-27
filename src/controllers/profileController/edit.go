package profileController

import (
	"crdx.org/lighthouse/m/repo/auditLogR"
	"crdx.org/lighthouse/middleware/auth"
	"crdx.org/lighthouse/pkg/flash"
	"crdx.org/lighthouse/pkg/globals"
	"crdx.org/lighthouse/pkg/transform"
	"crdx.org/lighthouse/pkg/util/stringutil"
	"crdx.org/lighthouse/pkg/validate"
	"crdx.org/session"
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

	// Expire any other sessions for this user ID (but not ours).
	auth.ExpireUserID(user.ID, session.GetID(c))

	auditLogR.Add(c, "Changed own password")
	flash.Success(c, "Password changed")
	return c.Redirect("/profile")
}
