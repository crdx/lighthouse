package flash

import (
	"encoding/gob"

	"crdx.org/session"
	"github.com/gofiber/fiber/v2"
)

const (
	Key          = "i_flash"
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

func New() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if flashMessage, found := session.GetOnce[Message](c, Key); found {
			c.Locals(Key, flashMessage)
		}

		return c.Next()
	}
}

func AddSuccess(c *fiber.Ctx, message string) {
	add(c, SuccessClass, message)
}

func AddFailure(c *fiber.Ctx, message string) {
	add(c, FailureClass, message)
}

func GetSuccess(message string) Message {
	return get(SuccessClass, message)
}

func GetFailure(message string) Message {
	return get(FailureClass, message)
}

func add(c *fiber.Ctx, class string, message string) {
	session.Set(c, Key, get(class, message))
}

func get(class string, message string) Message {
	return Message{
		Class:   class,
		Content: message,
	}
}
