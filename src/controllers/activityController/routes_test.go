package activityController_test

import (
	"testing"

	"crdx.org/lighthouse/middleware/auth"
	"crdx.org/lighthouse/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func setup() *helpers.Session {
	return helpers.Init(auth.StateAdmin)
}

func TestList(t *testing.T) {
	session := setup()
	res, body := session.Get("/activity")

	assert.Equal(t, 200, res.StatusCode)
	assert.Contains(t, body, "device1")
	assert.Contains(t, body, "device2")
	assert.NotContains(t, body, "device3")
	assert.Contains(t, body, "online")
	assert.Contains(t, body, "offline")
}
