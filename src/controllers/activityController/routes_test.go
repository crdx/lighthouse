package activityController_test

import (
	"testing"

	"crdx.org/lighthouse/constants"
	"crdx.org/lighthouse/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {
	session := helpers.Init(constants.RoleAdmin)

	res := session.Get("/activity")
	assert.Equal(t, 200, res.StatusCode)
	assert.Contains(t, res.Body, "device1-625a5fa0-9b63-46d8-b4fa-578f92dca041")
	assert.Contains(t, res.Body, "device2-64774746-5937-412c-9aa4-f262d990cc7d")
	assert.NotContains(t, res.Body, "device3-5acf7b73-b02c-4fe5-a63e-869f8bfc329e")
	assert.Contains(t, res.Body, "online")
	assert.Contains(t, res.Body, "offline")
}

func TestListDevice(t *testing.T) {
	session := helpers.Init(constants.RoleAdmin)

	res := session.Get("/activity/?device_id=2")
	assert.Equal(t, 200, res.StatusCode)
	assert.NotContains(t, res.Body, "device1-625a5fa0-9b63-46d8-b4fa-578f92dca041")
	assert.Contains(t, res.Body, "device2-64774746-5937-412c-9aa4-f262d990cc7d")
	assert.NotContains(t, res.Body, "device3-5acf7b73-b02c-4fe5-a63e-869f8bfc329e")
}

func TestListBadDevice(t *testing.T) {
	session := helpers.Init(constants.RoleAdmin)

	res := session.Get("/activity/?device_id=100")
	assert.Equal(t, 404, res.StatusCode)
}

func TestListBadPageNumber(t *testing.T) {
	session := helpers.Init(constants.RoleAdmin)

	res := session.Get("/activity/?p=100")
	assert.Equal(t, 404, res.StatusCode)

	res = session.Get("/activity/?p=0")
	assert.Equal(t, 404, res.StatusCode)
}
