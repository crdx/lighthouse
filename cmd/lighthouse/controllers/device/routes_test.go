package device_test

import (
	"testing"

	"crdx.org/lighthouse/cmd/lighthouse/tests/helpers"
	"crdx.org/lighthouse/db"
	"crdx.org/lighthouse/pkg/constants"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	helpers.TestMain(m)
}

func TestList(t *testing.T) {
	helpers.Start(t)
	session := helpers.NewSession(constants.RoleAdmin)

	res := session.Get("/")
	assert.Equal(t, 200, res.StatusCode)
	assert.Contains(t, res.Body, "AA:AA:AA:AA:AA:AA")
	assert.Contains(t, res.Body, "127.0.0.1")
	assert.Contains(t, res.Body, "device1-625a5fa0-9b63-46d8-b4fa-578f92dca041")
}

func TestListSort(t *testing.T) {
	helpers.Start(t)
	session := helpers.NewSession(constants.RoleAdmin)

	res := session.Get("/?sc=seen&sd=asc")
	assert.Equal(t, 200, res.StatusCode)
}

func TestListBadSort(t *testing.T) {
	helpers.Start(t)
	session := helpers.NewSession(constants.RoleAdmin)

	res := session.Get("/?sd=foo")
	assert.Equal(t, 400, res.StatusCode)

	res = session.Get("/?sc=foo")
	assert.Equal(t, 400, res.StatusCode)
}

func TestView(t *testing.T) {
	helpers.Start(t)
	session := helpers.NewSession(constants.RoleAdmin)

	res := session.Get("/device/1")
	assert.Equal(t, 200, res.StatusCode)
	assert.Contains(t, res.Body, "AA:AA:AA:AA:AA:AA")
	assert.Contains(t, res.Body, "127.0.0.1")
	assert.Contains(t, res.Body, "adapter1-1d6d5f93-e5bf-4651-ae9f-662cf01aad25")
	assert.Contains(t, res.Body, "Vendor 1")
}

func TestViewEdit(t *testing.T) {
	helpers.Start(t)
	session := helpers.NewSession(constants.RoleAdmin)

	res := session.Get("/device/1/edit")
	assert.Equal(t, 200, res.StatusCode)
	assert.Contains(t, res.Body, "device1-625a5fa0-9b63-46d8-b4fa-578f92dca041")
}

func TestViewerCannotViewEdit(t *testing.T) {
	helpers.Start(t)
	session := helpers.NewSession(constants.RoleViewer)

	res := session.Get("/device/1/edit")
	assert.Equal(t, 404, res.StatusCode)
	assert.NotContains(t, res.Body, "device1-625a5fa0-9b63-46d8-b4fa-578f92dca041")
}

func TestEdit(t *testing.T) {
	helpers.Start(t)
	session := helpers.NewSession(constants.RoleEditor)

	name := uuid.NewString()
	notes := uuid.NewString()

	res := session.PostForm("/device/1/edit", map[string]string{
		"name":         name,
		"notes":        notes,
		"icon":         "solid:vials",
		"grace_period": "5 mins",
		"limit":        "1 hour",
	})

	assert.Equal(t, 303, res.StatusCode)

	res = session.Get("/device/1")
	assert.Contains(t, res.Body, name)
	assert.Contains(t, res.Body, notes)
	assert.Contains(t, res.Body, "fa-solid fa-vials")
}

func TestEditOrigin(t *testing.T) {
	helpers.Start(t)
	session := helpers.NewSession(constants.RoleEditor)

	name := uuid.NewString()
	notes := uuid.NewString()

	res := session.PostForm("/device/3/edit", map[string]string{
		"name":  name,
		"notes": notes,
		"icon":  "solid:vials",
	})

	assert.Equal(t, 303, res.StatusCode)

	res = session.Get("/device/3")
	assert.Contains(t, res.Body, name)
	assert.Contains(t, res.Body, notes)
	assert.Contains(t, res.Body, "fa-solid fa-vials")
}

func TestEditWithErrors(t *testing.T) {
	helpers.Start(t)
	session := helpers.NewSession(constants.RoleEditor)

	name := uuid.NewString()
	notes := uuid.NewString()
	icon := uuid.NewString()

	res := session.PostForm("/device/1/edit", map[string]string{
		"name":         name,
		"notes":        notes,
		"icon":         icon,
		"grace_period": "",
		"limit":        "",
	})

	assert.Equal(t, 200, res.StatusCode)
	assert.Contains(t, res.Body, "required field")

	res = session.Get("/device/1")
	assert.NotContains(t, res.Body, notes)
	assert.NotContains(t, res.Body, icon)
}

func TestViewerCannotEdit(t *testing.T) {
	helpers.Start(t)
	session := helpers.NewSession(constants.RoleViewer)

	res := session.PostForm("/device/1/edit", map[string]string{
		"name": uuid.NewString(),
	})

	assert.Equal(t, 404, res.StatusCode)
}

func TestMerge1(t *testing.T) {
	helpers.Start(t)
	session := helpers.NewSession(constants.RoleAdmin)

	res := session.PostForm("/device/1/merge", map[string]string{
		"device_id": "2",
	})

	assert.Equal(t, 303, res.StatusCode)

	res = session.Get("/device/1")
	assert.Contains(t, res.Body, "2023-10-01")
	assert.Contains(t, res.Body, "adapter1-1d6d5f93-e5bf-4651-ae9f-662cf01aad25")
	assert.Contains(t, res.Body, "adapter2-c71739fd-d6f2-44e8-966f-fc5cdf2eec59")

	device := lo.Must(db.FindDevice(1))

	assert.Len(t, device.Adapters(), 2)
	assert.NotNil(t, device.DeletedAt)

	_, found := db.FindDevice(2)
	assert.False(t, found)
}

func TestMerge2(t *testing.T) {
	helpers.Start(t)
	session := helpers.NewSession(constants.RoleAdmin)

	res := session.PostForm("/device/2/merge", map[string]string{
		"device_id": "1",
	})

	assert.Equal(t, 303, res.StatusCode)

	res = session.Get("/device/1")
	assert.Contains(t, res.Body, "2023-10-01")
	assert.Contains(t, res.Body, "adapter1-1d6d5f93-e5bf-4651-ae9f-662cf01aad25")
	assert.Contains(t, res.Body, "adapter2-c71739fd-d6f2-44e8-966f-fc5cdf2eec59")

	device := lo.Must(db.FindDevice(1))

	assert.Len(t, device.Adapters(), 2)
	assert.NotNil(t, device.DeletedAt)

	_, found := db.FindDevice(2)
	assert.False(t, found)
}

func TestMergeBadDevice(t *testing.T) {
	helpers.Start(t)
	session := helpers.NewSession(constants.RoleAdmin)

	res := session.PostForm("/device/1/merge", map[string]string{
		"device_id": "100",
	})

	assert.Equal(t, 400, res.StatusCode)
}

func TestViewerCannotMerge(t *testing.T) {
	helpers.Start(t)
	session := helpers.NewSession(constants.RoleViewer)

	res := session.PostForm("/device/1/merge", map[string]string{
		"device_id": "2",
	})

	assert.Equal(t, 404, res.StatusCode)
}

func TestDelete(t *testing.T) {
	helpers.Start(t)
	session := helpers.NewSession(constants.RoleAdmin)

	res := session.Get("/device/1")
	assert.Equal(t, 200, res.StatusCode)

	res = session.PostForm("/device/1/delete", nil)
	assert.Equal(t, 303, res.StatusCode)

	res = session.Get("/device/1")
	assert.Equal(t, 404, res.StatusCode)
}

func TestViewerCannotDelete(t *testing.T) {
	helpers.Start(t)
	session := helpers.NewSession(constants.RoleViewer)

	res := session.PostForm("/device/1/delete", nil)
	assert.Equal(t, 404, res.StatusCode)
}
