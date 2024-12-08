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

	// GET ALL TODOS
	app.Get("/api/todos", func(c *fiber.Ctx) error {
		// Return all todos as JSON
		return c.Status(200).JSON(todos)
	})

	// CREATE A TODO
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

	// UPDATE A TODO
	// Define an HTTP PATCH route to update a Todo's "Completed" status
	app.Patch("/api/todos/:id", func(c *fiber.Ctx) error {
		// Extract the "id" parameter from the URL path
		// c.Params("id") retrieves the value of the `:id` placeholder in the route
		id := c.Params("id")

		// Iterate over the `todos` slice to find the Todo with a matching ID
		for i, todo := range todos {
			// Use `fmt.Sprint` to convert the integer `todo.ID` to a string for comparison
			if fmt.Sprint(todo.ID) == id {
				// If the ID matches, update the "Completed" status of the Todo to `true`
				todos[i].Completed = true

				// Respond to the client with the updated Todo and a 200 (OK) status
				return c.Status(200).JSON(todos[i])
			}
		}

		// If no Todo with the matching ID is found:
		// - Respond with a 404 (Not Found) status code
		// - Include a JSON error message for the client
		return c.Status(404).JSON(fiber.Map{"error": "Todo not found"})
	})

	// DELETE A TODO
	// Define an HTTP DELETE route to delete a Todo by its ID
	app.Delete("/api/todos/:id", func(c *fiber.Ctx) error {
		// Extract the "id" parameter from the URL path
		// c.Params("id") retrieves the value of the `:id` placeholder in the route
		id := c.Params("id")

		// Iterate over the `todos` slice to find the Todo with a matching ID
		for i, todo := range todos {
			// Use `fmt.Sprint` to convert the integer `todo.ID` to a string for comparison
			if fmt.Sprint(todo.ID) == id {
				// If the ID matches, remove the Todo from the `todos` slice
				// Use slicing to create a new slice excluding the matched item
				todos = append(todos[:i], todos[i+1:]...)

				// Respond to the client with a 200 (OK) status and a success message
				return c.Status(200).JSON(fiber.Map{"success": true})
			}
		}

		// If no Todo with the matching ID is found:
		// - Respond with a 404 (Not Found) status code
		// - Include a JSON error message for the client
		return c.Status(404).JSON(fiber.Map{"error": "Todo not found"})
	})

	// Start the web server on port 4000
	// log.Fatal ensures the application stops if the server fails to start
	log.Fatal(app.Listen(":4000"))
}
