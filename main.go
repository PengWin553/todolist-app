package main

// Import necessary packages
import (
	"fmt" // Used for printing messages to the console
	"log" // Used for logging errors

	"github.com/gofiber/fiber/v2" // The Fiber web framework for building APIs
)

// Define the structure of a Todo item
type Todo struct {
	ID        int    `json:"id"`        // Unique identifier for each Todo item (will be auto-assigned)
	Completed bool   `json:"completed"` // Indicates if the Todo is completed (true/false)
	Body      string `json:"body"`      // The content or description of the Todo
}

func main() {
	// Print a startup message to the console
	fmt.Println("Hello, World. Cat win")

	// Create a new Fiber application instance
	app := fiber.New()

	// A slice (dynamic array) to store Todo items in memory
	todos := []Todo{}

	// Define an HTTP POST route to create a new Todo
	app.Post("/api/todos", func(c *fiber.Ctx) error {
		// Create a new instance of the Todo struct
		todo := &Todo{}

		// Parse the JSON body from the incoming HTTP request into the Todo struct
		// c.BodyParser() reads the request body and maps it to the `todo` object
		if err := c.BodyParser(todo); err != nil {
			// If parsing fails, return a 400 (Bad Request) error with a helpful message
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
		}

		// Validate that the "body" field is not empty
		// If it's empty, return a 400 error with an appropriate message
		if todo.Body == "" {
			return c.Status(400).JSON(fiber.Map{"error": "Todo body is required"})
		}

		// Assign a unique ID to the new Todo
		// The ID is set to the current length of the `todos` slice + 1
		todo.ID = len(todos) + 1

		// Add the new Todo item to the `todos` slice
		todos = append(todos, *todo)

		// Respond to the client with a 201 (Created) status and the new Todo item in JSON format
		return c.Status(201).JSON(todo)
	})

	// Start the web server on port 4000
	// log.Fatal ensures the application stops if the server fails to start
	log.Fatal(app.Listen(":4000"))
}
