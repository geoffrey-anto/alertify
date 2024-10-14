package server

import (
	"github.com/gofiber/fiber/v2"

	"server/internal/database"

	"github.com/gofiber/fiber/v2/middleware/cors"
)

type FiberServer struct {
	*fiber.App

	db database.Service
}

func New() *FiberServer {
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "server",
			AppName:      "server",
		}),

		db: database.New(),
	}

	// Initialize default config
	server.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins:     "http://localhost:3000",
	}))

	return server
}
