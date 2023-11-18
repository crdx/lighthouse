package audit_test

import (
	"testing"

	"crdx.org/lighthouse/constants"
	"crdx.org/lighthouse/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	helpers.TestMain(m)
}

func TestList(t *testing.T) {
	defer helpers.Start()()
	session := helpers.NewSession(constants.RoleAdmin)

	res := session.Get("/admin/audit")
	assert.Equal(t, 200, res.StatusCode)

	assert.Contains(t, res.Body, "Edited device device1-625a5fa0-9b63-46d8-b4fa-578f92dca041")
	assert.Contains(t, res.Body, "root")

	assert.NotContains(t, res.Body, "Edited device device2-64774746-5937-412c-9aa4-f262d990cc7d")
	assert.NotContains(t, res.Body, "anon")
}

func TestViewerCannotList(t *testing.T) {
	defer helpers.Start()()
	session := helpers.NewSession(constants.RoleViewer)

	res := session.Get("/admin/audit")
	assert.Equal(t, 404, res.StatusCode)
}
