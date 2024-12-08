package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	fmt.Println("Hello, World. Cat win")
	app := fiber.New()

	// first route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status((200)).JSON(fiber.Map{"msg": "hello world!"})
	})

	log.Fatal(app.Listen(":4000"))
}
