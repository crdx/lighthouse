package helpers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"crdx.org/lighthouse/cmd/lighthouse/config"
	"crdx.org/lighthouse/pkg/constants"
	"crdx.org/lighthouse/pkg/middleware/auth"
	"github.com/gofiber/fiber/v3"
	"github.com/samber/lo"
)

type Session struct {
	cookies []*http.Cookie
	app     *fiber.App
}

type Response struct {
	*http.Response

	Body string
}

// NewSession returns a new session with the requested auth state.
func NewSession(role int64, handlers ...func(c fiber.Ctx) error) *Session {
	app := fiber.New(config.GetTestFiberConfig())

	app.Use(config.NewTestSessionMiddleware(dbConfig))

	if role == constants.RoleNone {
		app.Use(auth.New())
	} else {
		app.Use(auth.AutoLogin(role))
	}

	for _, handler := range handlers {
		app.Use(handler)
	}

	config.InitRoutes(app)

	return &Session{
		app: app,
	}
}

func (self *Session) Get(path string) *Response {
	req := httptest.NewRequest(http.MethodGet, path, nil)
	return self.send(req)
}

func (self *Session) PostForm(path string, formData map[string]string) *Response {
	return self.sendForm(http.MethodPost, path, formData)
}

func (self *Session) PostJSON(path string, jsonData any) *Response {
	return self.sendJSON(http.MethodPost, path, jsonData)
}

func (self *Session) sendJSON(method, path string, jsonData any) *Response {
	data := lo.Must(json.Marshal(jsonData))
	payload := bytes.NewBuffer(data)

	req := httptest.NewRequest(method, path, payload)
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	return self.send(req)
}

func (self *Session) sendForm(method, path string, formData map[string]string) *Response {
	data := url.Values{}
	for key, value := range formData {
		data.Set(key, value)
	}
	payload := strings.NewReader(data.Encode())

	req := httptest.NewRequest(method, path, payload)
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationForm)

	return self.send(req)
}

func (self *Session) send(req *http.Request) *Response {
	for _, cookie := range self.cookies {
		req.AddCookie(cookie)
	}

	res := lo.Must(self.app.Test(req))
	self.cookies = res.Cookies()

	return &Response{
		Response: res,
		Body:     string(lo.Must(io.ReadAll(res.Body))),
	}
}
