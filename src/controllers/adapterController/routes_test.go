package adapterController_test

import (
	"strings"
	"testing"

	"crdx.org/lighthouse/m/repo/userR"
	"crdx.org/lighthouse/tests/helpers"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestViewEdit(t *testing.T) {
	session := helpers.Init(userR.RoleAdmin)

	res := session.Get("/adapter/1/edit")
	assert.Equal(t, 200, res.StatusCode)
	assert.Contains(t, res.Body, "adapter1-1d6d5f93-e5bf-4651-ae9f-662cf01aad25")
}

func TestViewEditBadDevice(t *testing.T) {
	session := helpers.Init(userR.RoleAdmin)

	res := session.Get("/adapter/100/edit")
	assert.Equal(t, 404, res.StatusCode)
}

func TestViewerCannotViewEdit(t *testing.T) {
	session := helpers.Init(userR.RoleViewer)

	res := session.Get("/adapter/1/edit")
	assert.Equal(t, 404, res.StatusCode)
	assert.NotContains(t, res.Body, "adapter1-1d6d5f93-e5bf-4651-ae9f-662cf01aad25")
}

func TestEdit(t *testing.T) {
	session := helpers.Init(userR.RoleAdmin)

	name := uuid.NewString()
	vendor := uuid.NewString()

	res := session.PostForm("/adapter/1/edit", map[string]string{
		"name":   name,
		"vendor": vendor,
	})

	assert.Equal(t, 302, res.StatusCode)

	res = session.Get("/device/1")
	assert.Contains(t, res.Body, name)
	assert.Contains(t, res.Body, vendor)
}

func TestViewerCannotEdit(t *testing.T) {
	session := helpers.Init(userR.RoleViewer)

	name := uuid.NewString()
	vendor := uuid.NewString()

	res := session.PostForm("/adapter/1/edit", map[string]string{
		"name":   name,
		"vendor": vendor,
	})

	assert.Equal(t, 404, res.StatusCode)
	assert.NotContains(t, res.Body, name)
	assert.NotContains(t, res.Body, vendor)
}

func TestEditWithErrors(t *testing.T) {
	session := helpers.Init(userR.RoleAdmin)

	name := strings.Repeat(uuid.NewString(), 100)
	vendor := uuid.NewString()

	res := session.PostForm("/adapter/1/edit", map[string]string{
		"name":   name,
		"vendor": vendor,
	})

	assert.Equal(t, 200, res.StatusCode)
	assert.Contains(t, res.Body, name)
	assert.Contains(t, res.Body, vendor)
	assert.Contains(t, res.Body, "Unable to")
	assert.Contains(t, res.Body, "must be a maximum of")
	assert.Contains(t, res.Body, "characters in length")
}

func TestDelete(t *testing.T) {
	session := helpers.Init(userR.RoleAdmin)

	res := session.Get("/device/1/")
	assert.Contains(t, res.Body, "adapter1-1d6d5f93-e5bf-4651-ae9f-662cf01aad25")

	res = session.PostForm("/adapter/1/delete", nil)
	assert.Equal(t, 302, res.StatusCode)

	res = session.Get("/device/1")
	assert.NotContains(t, res.Body, "adapter1-1d6d5f93-e5bf-4651-ae9f-662cf01aad25")

	res = session.Get("/adapter/1/edit")
	assert.Equal(t, 404, res.StatusCode)
}

func TestViewerCannotDelete(t *testing.T) {
	session := helpers.Init(userR.RoleViewer)

	res := session.PostForm("/adapter/1/delete", nil)
	assert.Equal(t, 404, res.StatusCode)
}
