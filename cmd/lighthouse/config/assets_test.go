package config_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"

	"crdx.org/lighthouse/cmd/lighthouse/config"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const yearOne = "Mon, 01 Jan 0001 00:00:00 GMT"

func newAssetsApp(content string) *fiber.App {
	app := fiber.New()
	app.Use(config.NewCompress())
	config.InitAssets(app, fstest.MapFS{
		"assets/main.css":   {Data: []byte(content)},
		"assets/font.woff2": {Data: []byte(content)},
	})
	return app
}

func getAsset(t *testing.T, app *fiber.App, headers map[string]string) (*http.Response, string) {
	t.Helper()
	return getPath(t, app, "/assets/main.css", headers)
}

func getPath(t *testing.T, app *fiber.App, path string, headers map[string]string) (*http.Response, string) {
	t.Helper()

	req := httptest.NewRequest(http.MethodGet, path, nil)
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	res, err := app.Test(req)
	require.NoError(t, err)

	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	return res, string(body)
}

func TestAssetsRevalidate(t *testing.T) {
	res, body := getAsset(t, newAssetsApp("body { color: red }"), nil)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "body { color: red }", body)
	assert.Equal(t, "no-cache", res.Header.Get(fiber.HeaderCacheControl))
	assert.Empty(t, res.Header.Get(fiber.HeaderLastModified))
	assert.NotEmpty(t, res.Header.Get(fiber.HeaderETag))
}

func TestAssetsUnchangedIsNotModified(t *testing.T) {
	app := newAssetsApp("body { color: red }")
	first, _ := getAsset(t, app, nil)

	res, body := getAsset(t, app, map[string]string{
		fiber.HeaderIfNoneMatch: first.Header.Get(fiber.HeaderETag),
	})

	assert.Equal(t, http.StatusNotModified, res.StatusCode)
	assert.Empty(t, body)
	assert.Equal(t, "no-cache", res.Header.Get(fiber.HeaderCacheControl))
}

func TestAssetsChangedIsSentInFull(t *testing.T) {
	first, _ := getAsset(t, newAssetsApp("body { color: red }"), nil)

	res, body := getAsset(t, newAssetsApp("body { color: blue }"), map[string]string{
		fiber.HeaderIfNoneMatch:     first.Header.Get(fiber.HeaderETag),
		fiber.HeaderIfModifiedSince: yearOne,
	})

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "body { color: blue }", body)
}

func TestAssetsCompressedKeepETag(t *testing.T) {
	app := newAssetsApp("body { color: red }")
	plain, _ := getAsset(t, app, nil)

	res, _ := getAsset(t, app, map[string]string{
		fiber.HeaderAcceptEncoding: "gzip",
		fiber.HeaderIfNoneMatch:    plain.Header.Get(fiber.HeaderETag),
	})

	assert.Equal(t, http.StatusNotModified, res.StatusCode)
}

func TestCompressSkipsCompressedFormats(t *testing.T) {
	app := newAssetsApp(strings.Repeat("body { color: red }\n", 100))
	headers := map[string]string{fiber.HeaderAcceptEncoding: "br"}

	css, _ := getPath(t, app, "/assets/main.css", headers)
	font, _ := getPath(t, app, "/assets/font.woff2", headers)

	assert.Equal(t, "br", css.Header.Get(fiber.HeaderContentEncoding))
	assert.Empty(t, font.Header.Get(fiber.HeaderContentEncoding))
}

func TestVersionedAssetIsCachedForGood(t *testing.T) {
	app := newAssetsApp("body { color: red }")

	url := config.AssetURL("main.css")
	assert.Regexp(t, `^/assets/v[0-9a-f]{12}/main\.css$`, url)

	res, body := getPath(t, app, url, nil)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "body { color: red }", body)
	assert.Equal(t, "public, max-age=31536000, immutable", res.Header.Get(fiber.HeaderCacheControl))
}

func TestOtherVersionIsRevalidated(t *testing.T) {
	app := newAssetsApp("body { color: red }")

	res, body := getPath(t, app, "/assets/v000000000000/main.css", nil)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "body { color: red }", body)
	assert.Equal(t, "no-cache", res.Header.Get(fiber.HeaderCacheControl))
}

func TestVersionChangesWithContent(t *testing.T) {
	newAssetsApp("body { color: red }")
	red := config.AssetURL("main.css")

	newAssetsApp("body { color: blue }")
	assert.NotEqual(t, red, config.AssetURL("main.css"))
}
