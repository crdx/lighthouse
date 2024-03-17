package flash

import (
	"encoding/gob"
	"fmt"

	"crdx.org/session/v2"
	"github.com/gofiber/fiber/v2"
)

const (
	SuccessClass = "success"
	FailureClass = "danger"
)

type Message struct {
	Class   string
	Content string
}

func init() {
	gob.Register(Message{})
}

func Success(c *fiber.Ctx, message string, args ...any) {
	add(c, SuccessClass, message, args...)
}

func Failure(c *fiber.Ctx, message string, args ...any) {
	add(c, FailureClass, message, args...)
}

func add(c *fiber.Ctx, class string, message string, args ...any) {
	session.Set(c, "globals.flash", get(class, fmt.Sprintf(message, args...)))
}

func get(class string, message string) Message {
	return Message{
		Class:   class,
		Content: message,
	}
}
