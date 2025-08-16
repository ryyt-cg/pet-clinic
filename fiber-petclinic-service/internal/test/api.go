package test

import "github.com/gofiber/fiber/v2"

// SetupRouter config the fiber engine for testing purpose
func SetupRouter() *fiber.App {
	r := fiber.New()
	return r
}
