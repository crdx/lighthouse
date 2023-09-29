package helpers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
)

func send(app *fiber.App, req *http.Request) (*http.Response, string) {
	res := lo.Must(app.Test(req))
	return res, string(lo.Must(io.ReadAll(res.Body)))
}

func Get(app *fiber.App, path string) (*http.Response, string) {
	req := httptest.NewRequest(http.MethodGet, path, nil)

	return send(app, req)
}

func SendForm(app *fiber.App, method, path string, formData map[string]string) (*http.Response, string) {
	data := url.Values{}
	for key, value := range formData {
		data.Set(key, value)
	}
	payload := strings.NewReader(data.Encode())

	req := httptest.NewRequest(method, path, payload)
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationForm)

	return send(app, req)
}

func SendJSON(app *fiber.App, method, path string, jsonData any) (*http.Response, string) {
	data := lo.Must(json.Marshal(jsonData))
	payload := bytes.NewBuffer(data)

	req := httptest.NewRequest(method, path, payload)
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	return send(app, req)
}

func PostForm(app *fiber.App, path string, formData map[string]string) (*http.Response, string) {
	return SendForm(app, http.MethodPost, path, formData)
}

func PostJSON(app *fiber.App, path string, jsonData any) (*http.Response, string) {
	return SendJSON(app, http.MethodPost, path, jsonData)
}
