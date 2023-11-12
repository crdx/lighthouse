package audit_test

import (
	"testing"

	"crdx.org/lighthouse/middleware/auth"
	"crdx.org/lighthouse/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {
	session := helpers.Init(auth.StateAdmin)

	res := session.Get("/admin/audit")
	assert.Equal(t, 200, res.StatusCode)

	assert.Contains(t, res.Body, "Edited device device1-625a5fa0-9b63-46d8-b4fa-578f92dca041")
	assert.Contains(t, res.Body, "root")

	assert.NotContains(t, res.Body, "Edited device device2-64774746-5937-412c-9aa4-f262d990cc7d")
	assert.NotContains(t, res.Body, "anon")
}

func TestUserCannotList(t *testing.T) {
	session := helpers.Init(auth.StateUser)

	res := session.Get("/admin/audit")
	assert.Equal(t, 404, res.StatusCode)
}
