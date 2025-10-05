package users

import (
	"errors"
	"strconv"

	"crdx.org/lighthouse/db"
	"crdx.org/lighthouse/db/repo/auditLogR"
	"crdx.org/lighthouse/pkg/flash"
	"crdx.org/lighthouse/pkg/globals"
	"crdx.org/lighthouse/pkg/transform"
	"crdx.org/lighthouse/pkg/util/stringutil"
	"crdx.org/lighthouse/pkg/validate"
	"github.com/gofiber/fiber/v3"
	"github.com/samber/lo"
)

type CreateForm struct {
	Username        string `form:"username" validate:"required,max=20"`
	Password        string `form:"password" validate:"required,min=4"`
	ConfirmPassword string `form:"confirm_password" validate:"required"`
	Role            string `form:"role" validate:"required,role"`
}

func ViewCreate(c fiber.Ctx) error {
	return c.Render("admin/index", fiber.Map{
		"tab":     "users",
		"mode":    "create",
		"fields":  validate.Fields[CreateForm](),
		"globals": globals.Get(c),
	})
}

func Create(c fiber.Ctx) error {
	form := new(CreateForm)
	lo.Must0(c.Bind().Body(form))
	transform.Struct(form)

	validatorMap := validate.ValidatorMap{
		"Username": func(value string) error {
			for _, user := range db.FindUsers() {
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

	user := db.User{
		Username: form.Username,
		Role:     lo.Must(strconv.ParseInt(form.Role, 10, 64)),
	}

	if form.Password != "" {
		user.PasswordHash = stringutil.Hash(form.Password)
	}

	db.CreateUser(&user)

	auditLogR.Add(c, "Created user %s", user.AuditName())
	flash.Success(c, "User created")
	return c.Redirect().To("/admin/users")
}
