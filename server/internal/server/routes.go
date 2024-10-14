package server

import (
	"encoding/json"
	"fmt"
	"server/internal/models"
	"server/internal/utils"
	"time"

	"github.com/gofiber/fiber/v2"
)

func (s *FiberServer) RegisterFiberRoutes() {
	s.App.Get("/", s.HelloWorldHandler)

	s.App.Get("/health", s.healthHandler)

	api := s.App.Group("/api")

	v1 := api.Group("/v1")

	v1.Post("/login", s.LoginHandler)

	v1.Post("/register", s.RegisterUserHandler)

}

func (s *FiberServer) HelloWorldHandler(c *fiber.Ctx) error {
	resp := fiber.Map{
		"message": "Hello World",
	}

	return c.JSON(resp)
}

func (s *FiberServer) healthHandler(c *fiber.Ctx) error {
	return c.JSON(s.db.Health())
}

func (s *FiberServer) LoginHandler(c *fiber.Ctx) error {
	// get email and password from request body

	user := new(models.UserLogin)

	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	// check if user exists in database
	userFromDatabase, err := s.db.GetUser(user.Email)

	fmt.Printf("User from database: %v\n", err)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// check if password is correct
	fmt.Printf("User from database: %v\n", userFromDatabase)
	fmt.Printf("User from request: %v\n", user)

	if !utils.CheckPasswordHash(user.Pass, userFromDatabase.Pass) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid password",
		})
	}

	userJson, err := json.Marshal(fiber.Map{
		"email": user.Email,
		"fname": userFromDatabase.Fname,
		"lname": userFromDatabase.Lname,
		"id":    userFromDatabase.ID,
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	token, err := utils.CreateToken(string(userJson))

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		Secure:   true,
		HTTPOnly: true,
		Expires:  time.Now().Add(time.Hour * 24),
	})

	c.Status(fiber.StatusOK)
	return c.JSON(fiber.Map{
		"email": user.Email,
	})
}

func (s *FiberServer) RegisterUserHandler(c *fiber.Ctx) error {
	user := new(models.User)

	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	// hash password
	hashedPassword, err := utils.HashPassword(user.Pass)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	err = s.db.SaveUser(user.Email, hashedPassword, user.Fname, user.Lname)

	if err != nil {
		fmt.Printf("Error saving user: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "User could not be saved! User already exists",
		})
	}

	c.Status(fiber.StatusCreated)

	return c.JSON(fiber.Map{
		"email": user.Email,
	})
}
