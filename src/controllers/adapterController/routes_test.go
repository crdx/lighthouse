package adapterController_test

import (
	"testing"

	"crdx.org/lighthouse/middleware/auth"
	"crdx.org/lighthouse/tests/helpers"
	"crdx.org/lighthouse/util/stringutil"
	"github.com/stretchr/testify/assert"
)

func setup() *helpers.Session {
	return helpers.Init(auth.StateAdmin)
}

func TestEdit(t *testing.T) {
	session := setup()

	nameUUID := stringutil.UUID()
	vendorUUID := stringutil.UUID()

	res := session.PostForm("/adapter/1/edit", map[string]string{
		"name":   nameUUID,
		"vendor": vendorUUID,
	})

	assert.Equal(t, 302, res.StatusCode)

	res = session.Get("/device/1")

	assert.Contains(t, res.Body, nameUUID)
	assert.Contains(t, res.Body, vendorUUID)
}

func TestDelete(t *testing.T) {
	session := setup()

	res := session.Get("/device/1/")
	assert.Contains(t, res.Body, "adapter1")

	res = session.PostForm("/adapter/1/delete", nil)
	assert.Equal(t, 302, res.StatusCode)

	res = session.Get("/device/1")
	assert.NotContains(t, res.Body, "adapter1")

	res = session.Get("/adapter/1/edit")
	assert.Equal(t, 404, res.StatusCode)
}
