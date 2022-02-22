package main

import "github.com/gofiber/fiber/v2"

func setupRoutes(app *fiber.App) {

	// Routes with Handlers
	app.Get("/", home())

	// Start: Magic
	app.Post("/magic", parser())

}
