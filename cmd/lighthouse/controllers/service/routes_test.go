package service_test

import (
	"strings"
	"testing"

	"crdx.org/lighthouse/cmd/lighthouse/tests/helpers"
	"crdx.org/lighthouse/pkg/constants"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	helpers.TestMain(m)
}

func TestViewEdit(t *testing.T) {
	helpers.Start(t)
	session := helpers.NewSession(constants.RoleAdmin)

	res := session.Get("/service/1/edit")
	assert.Equal(t, 200, res.StatusCode)
	assert.Contains(t, res.Body, "service1-f6d0b172-7e23-4d6c-a9bd-e456208c01fe")
}

func TestViewEditBadService(t *testing.T) {
	helpers.Start(t)
	session := helpers.NewSession(constants.RoleAdmin)

	res := session.Get("/service/100/edit")
	assert.Equal(t, 404, res.StatusCode)
}

func TestViewerCannotViewEdit(t *testing.T) {
	helpers.Start(t)
	session := helpers.NewSession(constants.RoleViewer)

	res := session.Get("/service/1/edit")
	assert.Equal(t, 404, res.StatusCode)
	assert.NotContains(t, res.Body, "service1-f6d0b172-7e23-4d6c-a9bd-e456208c01fe")
}

func TestEdit(t *testing.T) {
	helpers.Start(t)
	session := helpers.NewSession(constants.RoleEditor)

	name := uuid.NewString()

	res := session.PostForm("/service/1/edit", map[string]string{
		"name": name,
	})

	assert.Equal(t, 303, res.StatusCode)

	res = session.Get("/device/1")
	assert.Contains(t, res.Body, name)
}

func TestViewerCannotEdit(t *testing.T) {
	helpers.Start(t)
	session := helpers.NewSession(constants.RoleViewer)

	name := uuid.NewString()

	res := session.PostForm("/service/1/edit", map[string]string{
		"name": name,
	})

	assert.Equal(t, 404, res.StatusCode)
	assert.NotContains(t, res.Body, name)
}

func TestEditWithErrors(t *testing.T) {
	helpers.Start(t)
	session := helpers.NewSession(constants.RoleEditor)

	name := strings.Repeat(uuid.NewString(), 100)

	res := session.PostForm("/service/1/edit", map[string]string{
		"name": name,
	})

	assert.Equal(t, 200, res.StatusCode)
	assert.Contains(t, res.Body, name)
	assert.Contains(t, res.Body, "Unable to")
	assert.Contains(t, res.Body, "must be a maximum of")
	assert.Contains(t, res.Body, "characters in length")
}

func TestDelete(t *testing.T) {
	helpers.Start(t)
	session := helpers.NewSession(constants.RoleAdmin)

	res := session.Get("/device/1/")
	assert.Contains(t, res.Body, "service1-f6d0b172-7e23-4d6c-a9bd-e456208c01fe")

	res = session.PostForm("/service/1/delete", nil)
	assert.Equal(t, 303, res.StatusCode)

	res = session.Get("/device/1")
	assert.NotContains(t, res.Body, "service1-f6d0b172-7e23-4d6c-a9bd-e456208c01fe")

	res = session.Get("/service/1/edit")
	assert.Equal(t, 404, res.StatusCode)
}

func TestViewerCannotDelete(t *testing.T) {
	helpers.Start(t)
	session := helpers.NewSession(constants.RoleViewer)

	res := session.PostForm("/service/1/delete", nil)
	assert.Equal(t, 404, res.StatusCode)
}
