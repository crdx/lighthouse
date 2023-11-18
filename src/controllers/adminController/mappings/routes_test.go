package mappings_test

import (
	"fmt"
	"testing"

	"crdx.org/db"
	"crdx.org/lighthouse/constants"
	"crdx.org/lighthouse/m"
	"crdx.org/lighthouse/tests/helpers"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	helpers.TestMain(m)
}

func TestList(t *testing.T) {
	defer helpers.Start()()
	session := helpers.NewSession(constants.RoleAdmin)

	res := session.Get("/admin/mappings")
	assert.Equal(t, 200, res.StatusCode)
}

func TestSaveSources(t *testing.T) {
	defer helpers.Start()()
	session := helpers.NewSession(constants.RoleAdmin)

	res := session.PostForm("/admin/mappings", map[string]string{
		"source_mac_addresses": "AA:AA:AA:AA:AA:AA",
	})

	assert.Equal(t, 302, res.StatusCode)

	res = session.Get("/admin/mappings")
	assert.Equal(t, 200, res.StatusCode)
	assert.Contains(t, res.Body, "AA:AA:AA:AA:AA:AA")
}

func TestSaveSourcesWithErrors(t *testing.T) {
	defer helpers.Start()()
	session := helpers.NewSession(constants.RoleAdmin)

	res := session.PostForm("/admin/mappings", map[string]string{
		"source_mac_addresses": "FF:00:00",
	})

	assert.Equal(t, 200, res.StatusCode)
	assert.Contains(t, res.Body, "must be a valid list of MAC addresses")
}

func TestAddMapping(t *testing.T) {
	defer helpers.Start()()
	session := helpers.NewSession(constants.RoleAdmin)

	label := uuid.NewString()[:20]

	res := session.PostForm("/admin/mappings/add", map[string]string{
		"label":       label,
		"mac_address": "FF:FF:FF:FF:FF:FF",
		"ip_address":  "127.0.0.1",
	})

	assert.Equal(t, 302, res.StatusCode)

	res = session.Get("/admin/mappings")
	assert.Equal(t, 200, res.StatusCode)
	assert.Contains(t, res.Body, label)
	assert.Contains(t, res.Body, "127.0.0.1")
	assert.Contains(t, res.Body, "FF:FF:FF:FF:FF:FF")
}

func TestCannotAddMappingWithDuplicateIP(t *testing.T) {
	defer helpers.Start()()
	session := helpers.NewSession(constants.RoleAdmin)

	label := uuid.NewString()[:20]

	res := session.PostForm("/admin/mappings/add", map[string]string{
		"label":       label,
		"mac_address": "AA:AA:AA:AA:AA:AA",
		"ip_address":  "127.0.0.1",
	})

	assert.Equal(t, 302, res.StatusCode)

	res = session.PostForm("/admin/mappings/add", map[string]string{
		"label":       label,
		"mac_address": "BB:BB:BB:BB:BB:BB",
		"ip_address":  "127.0.0.1",
	})

	assert.Equal(t, 200, res.StatusCode)
	assert.Contains(t, res.Body, "must be a unique IP address")
}

func TestDeleteMapping(t *testing.T) {
	defer helpers.Start()()
	session := helpers.NewSession(constants.RoleAdmin)

	label := uuid.NewString()[:20]

	res := session.PostForm("/admin/mappings/add", map[string]string{
		"label":       label,
		"mac_address": "FF:FF:FF:FF:FF:FF",
		"ip_address":  "127.0.0.1",
	})

	assert.Equal(t, 302, res.StatusCode)
	mapping, _ := db.B[m.Mapping]().First()

	res = session.PostForm(fmt.Sprintf("/admin/mappings/%d/delete", mapping.ID), map[string]string{})
	assert.Equal(t, 302, res.StatusCode)

	res = session.Get("/admin/mappings")
	assert.Equal(t, 200, res.StatusCode)
	assert.NotContains(t, res.Body, label)
	assert.NotContains(t, res.Body, "127.0.0.1")
	assert.NotContains(t, res.Body, "FF:FF:FF:FF:FF:FF")
}
