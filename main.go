package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	fmt.Println("Hello, World. Cat win")

	app := fiber.New()

	log.Fatal(app.Listen(":4000"))
}
