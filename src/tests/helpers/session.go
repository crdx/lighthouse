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

type Session struct {
	cookies []*http.Cookie
	app     *fiber.App
}

func NewSession(app *fiber.App) *Session {
	return &Session{
		app: app,
	}
}

func (self *Session) Get(path string) (*http.Response, string) {
	req := httptest.NewRequest(http.MethodGet, path, nil)
	return self.send(req)
}

func (self *Session) PostForm(path string, formData map[string]string) (*http.Response, string) {
	return self.sendForm(http.MethodPost, path, formData)
}

func (self *Session) PostJSON(path string, jsonData any) (*http.Response, string) {
	return self.sendJSON(http.MethodPost, path, jsonData)
}

func (self *Session) sendJSON(method, path string, jsonData any) (*http.Response, string) {
	data := lo.Must(json.Marshal(jsonData))
	payload := bytes.NewBuffer(data)

	req := httptest.NewRequest(method, path, payload)
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	return self.send(req)
}

func (self *Session) sendForm(method, path string, formData map[string]string) (*http.Response, string) {
	data := url.Values{}
	for key, value := range formData {
		data.Set(key, value)
	}
	payload := strings.NewReader(data.Encode())

	req := httptest.NewRequest(method, path, payload)
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationForm)

	return self.send(req)
}

func (self *Session) send(req *http.Request) (*http.Response, string) {
	for _, cookie := range self.cookies {
		req.AddCookie(cookie)
	}

	res := lo.Must(self.app.Test(req))
	self.cookies = res.Cookies()

	return res, string(lo.Must(io.ReadAll(res.Body)))
}
