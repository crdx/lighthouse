package notificationController_test

import (
	"testing"

	"crdx.org/lighthouse/middleware/auth"
	"crdx.org/lighthouse/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {
	session := helpers.Init(auth.StateAdmin)

	res := session.Get("/notifications")
	assert.Equal(t, 200, res.StatusCode)
	assert.Contains(t, res.Body, "subject-8f3fdfea-f39c-427f-b8f5-0155119975ff")
	assert.Contains(t, res.Body, "body-be67c77f-595c-4c91-8d24-99c829de1bbe")
}

func TestListBadPageNumber(t *testing.T) {
	session := helpers.Init(auth.StateAdmin)

	res := session.Get("/notifications/?p=100")
	assert.Equal(t, 404, res.StatusCode)

	res = session.Get("/notifications/?p=0")
	assert.Equal(t, 404, res.StatusCode)
}
