package adapterController_test

import (
	"testing"

	"crdx.org/lighthouse/controllers/adapterController"
	"crdx.org/lighthouse/controllers/deviceController"
	"crdx.org/lighthouse/middleware/auth"
	"crdx.org/lighthouse/tests/helpers"
	"crdx.org/lighthouse/util/stringutil"
	"github.com/stretchr/testify/assert"
)

func setup() *helpers.Session {
	helpers.Init()
	app := helpers.App(auth.StateAdmin)
	adapterController.InitRoutes(app)
	deviceController.InitRoutes(app)
	return helpers.NewSession(app)
}

func TestEdit(t *testing.T) {
	session := setup()

	nameUUID := stringutil.UUID()
	vendorUUID := stringutil.UUID()

	res, _ := session.PostForm("/adapter/1/edit", map[string]string{
		"name":   nameUUID,
		"vendor": vendorUUID,
	})

	assert.Equal(t, 302, res.StatusCode)

	_, body := session.Get("/device/1")

	assert.Contains(t, body, nameUUID)
	assert.Contains(t, body, vendorUUID)
}

func TestDelete(t *testing.T) {
	session := setup()

	_, body := session.Get("/device/1/")
	assert.Contains(t, body, "adapter1")

	res, _ := session.PostForm("/adapter/1/delete", nil)
	assert.Equal(t, 302, res.StatusCode)

	_, body = session.Get("/device/1")
	assert.NotContains(t, body, "adapter1")

	res, _ = session.Get("/adapter/1/edit")
	assert.Equal(t, 404, res.StatusCode)
}
