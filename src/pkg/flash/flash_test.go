package flash_test

import (
	"fmt"
	"testing"

	"crdx.org/lighthouse/middleware/auth"
	"crdx.org/lighthouse/pkg/flash"
	"crdx.org/lighthouse/tests/helpers"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestFlash(t *testing.T) {
	testCases := []struct {
		success bool
		message string
		class   string
	}{
		{true, uuid.NewString(), flash.SuccessClass},
		{false, uuid.NewString(), flash.FailureClass},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("%s,%v", testCase.message, testCase.success), func(t *testing.T) {
			session := helpers.Init(auth.StateAdmin, func(c *fiber.Ctx) error {
				if testCase.success {
					flash.Success(c, testCase.message)
				} else {
					flash.Failure(c, testCase.message)
				}
				return c.Next()
			})

			res := session.Get("/")
			assert.Equal(t, 200, res.StatusCode)
			assert.Contains(t, res.Body, testCase.message)
			assert.Contains(t, res.Body, testCase.class)
		})
	}
}
