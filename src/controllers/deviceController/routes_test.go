package deviceController_test

import (
	"testing"

	"crdx.org/db"
	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/middleware/auth"
	"crdx.org/lighthouse/tests/helpers"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func setup() *helpers.Session {
	return helpers.Init(auth.StateAdmin)
}

func TestList(t *testing.T) {
	session := setup()
	res := session.Get("/")

	assert.Equal(t, 200, res.StatusCode)
	assert.Contains(t, res.Body, "AA:AA:AA:AA:AA:AA")
	assert.Contains(t, res.Body, "127.0.0.1")
	assert.Contains(t, res.Body, "device1")
}

func TestView(t *testing.T) {
	session := setup()
	res := session.Get("/device/1")
	assert.Equal(t, 200, res.StatusCode)
	assert.Contains(t, res.Body, "AA:AA:AA:AA:AA:AA")
	assert.Contains(t, res.Body, "127.0.0.1")
	assert.Contains(t, res.Body, "adapter1")
	assert.Contains(t, res.Body, "Corp 1")
}

func TestEdit(t *testing.T) {
	session := setup()

	nameUUID := uuid.NewString()
	notesUUID := uuid.NewString()
	iconUUID := uuid.NewString()

	res := session.PostForm("/device/1/edit", map[string]string{
		"name":         nameUUID,
		"notes":        notesUUID,
		"icon":         iconUUID,
		"grace_period": "6",
	})

	assert.Equal(t, 302, res.StatusCode)

	res = session.Get("/device/1")

	assert.Contains(t, res.Body, nameUUID)
	assert.Contains(t, res.Body, notesUUID)
	assert.Contains(t, res.Body, iconUUID)
}

func TestEditWithErrors(t *testing.T) {
	session := setup()

	nameUUID := uuid.NewString()
	notesUUID := uuid.NewString()
	iconUUID := uuid.NewString()

	res := session.PostForm("/device/1/edit", map[string]string{
		"name":         nameUUID,
		"notes":        notesUUID,
		"icon":         iconUUID,
		"grace_period": "",
	})

	assert.Equal(t, 200, res.StatusCode)
	assert.Contains(t, res.Body, "required field")

	res = session.Get("/device/1")

	assert.NotContains(t, res.Body, notesUUID)
	assert.NotContains(t, res.Body, iconUUID)
}

func TestMerge(t *testing.T) {
	session := setup()

	res := session.PostForm("/device/1/merge", map[string]string{
		"device_id": "2",
	})

	assert.Equal(t, 302, res.StatusCode)

	res = session.Get("/device/1")

	assert.Contains(t, res.Body, "2023-10-01")
	assert.Contains(t, res.Body, "adapter1")
	assert.Contains(t, res.Body, "adapter2")

	device := lo.Must(db.First[m.Device](1))

	assert.Len(t, device.Adapters(), 2)
	assert.NotNil(t, device.DeletedAt)

	_, found := db.First[m.Device](2)
	assert.False(t, found)
}

func TestDelete(t *testing.T) {
	session := setup()

	res := session.Get("/device/1")
	assert.Equal(t, 200, res.StatusCode)

	res = session.PostForm("/device/1/delete", nil)
	assert.Equal(t, 302, res.StatusCode)

	res = session.Get("/device/1")
	assert.Equal(t, 404, res.StatusCode)
}
