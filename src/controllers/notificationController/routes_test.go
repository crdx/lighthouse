package notificationController_test

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
	res := session.Get("/notifications")

	assert.Equal(t, 200, res.StatusCode)
	assert.Contains(t, res.Body, "a thing has happened")
	assert.Contains(t, res.Body, "here are more details about the thing that happened")
}
