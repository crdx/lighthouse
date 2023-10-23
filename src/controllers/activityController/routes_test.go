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
	res := session.Get("/activity")

	assert.Equal(t, 200, res.StatusCode)
	assert.Contains(t, res.Body, "device1")
	assert.Contains(t, res.Body, "device2")
	assert.NotContains(t, res.Body, "device3")
	assert.Contains(t, res.Body, "online")
	assert.Contains(t, res.Body, "offline")
}
