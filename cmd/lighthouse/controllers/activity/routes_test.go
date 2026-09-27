package activity_test

import (
	"strings"
	"testing"

	"crdx.org/lighthouse/cmd/lighthouse/tests/helpers"
	"crdx.org/lighthouse/pkg/constants"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	helpers.TestMain(m)
}

func TestList(t *testing.T) {
	helpers.Start(t)
	session := helpers.NewSession(constants.RoleAdmin)

	res := session.Get("/activity")
	assert.Equal(t, 200, res.StatusCode)
	assert.Contains(t, res.Body, "device1-625a5fa0-9b63-46d8-b4fa-578f92dca041")
	assert.Contains(t, res.Body, "device2-64774746-5937-412c-9aa4-f262d990cc7d")
	assert.NotContains(t, res.Body, "device3-5acf7b73-b02c-4fe5-a63e-869f8bfc329e")
	assert.Contains(t, res.Body, "online")
	assert.Contains(t, res.Body, "offline")
}

func TestListSort(t *testing.T) {
	helpers.Start(t)
	session := helpers.NewSession(constants.RoleAdmin)

	res := session.Get("/activity/?sc=name&sd=asc")
	assert.Equal(t, 200, res.StatusCode)
	assert.Less(
		t,
		strings.Index(res.Body, "device1-625a5fa0-9b63-46d8-b4fa-578f92dca041"),
		strings.Index(res.Body, "device2-64774746-5937-412c-9aa4-f262d990cc7d"),
	)

	res = session.Get("/activity/?sc=name&sd=desc")
	assert.Equal(t, 200, res.StatusCode)
	assert.Less(
		t,
		strings.Index(res.Body, "device2-64774746-5937-412c-9aa4-f262d990cc7d"),
		strings.Index(res.Body, "device1-625a5fa0-9b63-46d8-b4fa-578f92dca041"),
	)
}

func TestListBadSort(t *testing.T) {
	helpers.Start(t)
	session := helpers.NewSession(constants.RoleAdmin)

	res := session.Get("/activity/?sd=foo")
	assert.Equal(t, 400, res.StatusCode)

	res = session.Get("/activity/?sc=foo")
	assert.Equal(t, 400, res.StatusCode)
}

func TestListDevice(t *testing.T) {
	helpers.Start(t)
	session := helpers.NewSession(constants.RoleAdmin)

	res := session.Get("/activity/?device_id=2")
	assert.Equal(t, 200, res.StatusCode)
	assert.NotContains(t, res.Body, "device1-625a5fa0-9b63-46d8-b4fa-578f92dca041")
	assert.Contains(t, res.Body, "device2-64774746-5937-412c-9aa4-f262d990cc7d")
	assert.NotContains(t, res.Body, "device3-5acf7b73-b02c-4fe5-a63e-869f8bfc329e")
	assert.Contains(t, res.Body, "?sc=name&sd=asc&device_id=2")
}

func TestListBadDevice(t *testing.T) {
	helpers.Start(t)
	session := helpers.NewSession(constants.RoleAdmin)

	res := session.Get("/activity/?device_id=100")
	assert.Equal(t, 404, res.StatusCode)
}

func TestListBadPageNumber(t *testing.T) {
	helpers.Start(t)
	session := helpers.NewSession(constants.RoleAdmin)

	res := session.Get("/activity/?p=100")
	assert.Equal(t, 404, res.StatusCode)

	res = session.Get("/activity/?p=0")
	assert.Equal(t, 404, res.StatusCode)
}
