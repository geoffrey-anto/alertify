package server

import (
	"encoding/json"
	"server/internal/utils"

	"github.com/gofiber/fiber/v2"
)

func (s *FiberServer) AuthMiddleware(c *fiber.Ctx) error {
	cookie := c.Cookies("token")
	authHeader := c.Get("Authorization")

	if cookie == "" && authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	var token string
	var err error

	if cookie != "" {
		err = utils.VerifyToken(cookie)
	} else {
		// Extract token from Authorization Bearer header
		// Assuming the header format is "Bearer <token>"
		token = authHeader[len("Bearer "):]
		err = nil
	}

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	err = utils.VerifyToken(token)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	type UserPayload struct {
		Email string `json:"email"`
		FName string `json:"fname"`
		LName string `json:"lname"`
		ID    int    `json:"id"`
	}

	user := new(UserPayload)

	payload, err := utils.GetPayload(token)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	err = json.Unmarshal([]byte(payload), &user)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	c.Locals("user_id", user.ID)
	c.Locals("user_email", user.Email)
	c.Locals("user_fname", user.FName)
	c.Locals("user_lname", user.LName)

	return c.Next()
}
