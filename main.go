// Package declaration for a standalone executable
package main

// Import necessary packages
import (
	"context" // Provides support for managing request-scoped values, cancellation signals, and deadlines
	"fmt"     // Implements formatted I/O functions for printing
	"log"     // Provides logging functionality
	"os"      // Provides a platform-independent interface to operating system functionality

	"github.com/gofiber/fiber/v2" // High-performance web framework for Go
	"github.com/joho/godotenv"    // Loads environment variables from a .env file
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"         // MongoDB Go driver for database operations
	"go.mongodb.org/mongo-driver/mongo/options" // Provides configuration options for MongoDB client
)

// Todo struct defines the structure of a todo item in the application
// This represents how a todo will be stored in the MongoDB database
type Todo struct {
	// ID is the unique identifier for each Todo item
	// `json:"id"` defines how the field is serialized to JSON
	// `bson:"_id"` specifies how the field is stored in MongoDB (MongoDB uses '_id' as the default ID field)
	// omitempty means the field will be omitted from the output if it's empty
	ID primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`

	// Completed indicates whether the todo item has been finished
	// Boolean field that can be true (completed) or false (not completed)
	Completed bool `json:"completed"`

	// Body contains the text description of the todo item
	Body string `json:"body"`
}

// Declare a global variable to hold the MongoDB collection
// This will be used to interact with the specific collection in the database
var collection *mongo.Collection

// main is the entry point of the Go application
// It sets up the database connection and configures the web server
func main() {
	// Print a greeting to confirm the application is starting
	fmt.Println("Hello, world!")

	// Load environment variables from .env file
	// This allows storing sensitive information like database credentials outside the code
	err := godotenv.Load(".env")
	if err != nil {
		// If .env file can't be loaded, terminate the program with an error message
		log.Fatal("Error loading .env file:", err)
	}

	// Retrieve MongoDB connection URI from environment variables
	// This keeps sensitive connection information out of the source code
	MONGODB_URI := os.Getenv("MONGODB_URI")

	// Create client options using the MongoDB URI
	// This configures how the application will connect to the MongoDB database
	clientOptions := options.Client().ApplyURI(MONGODB_URI)

	// Establish a connection to the MongoDB database
	// context.Background() provides a default context
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		// If connection fails, terminate the program
		log.Fatal(err)
	}

	// Ensure the database connection is closed when the application exits
	defer client.Disconnect(context.Background())

	// Verify the database connection by pinging the server
	// This ensures that the connection is active and working
	err = client.Ping(context.Background(), nil)
	if err != nil {
		// If ping fails, terminate the program
		log.Fatal(err)
	}

	// Print a success message when connection is established
	fmt.Println("Connected to MONGODB Atlas")

	// Set up the specific collection in the database that will be used for todos
	collection = client.Database("todolist_app_db").Collection("todos")

	// Create a new Fiber web application instance
	app := fiber.New()

	// Define API routes for CRUD operations on todos
	app.Get("/api/todos", getTodos)          // Retrieve all todos
	app.Post("/api/todos", createTodo)       // Create a new todo
	app.Patch("/api/todos/:id", updateTodo)  // Update an existing todo
	app.Delete("/api/todos/:id", deleteTodo) // Delete a todo

	// Determine the port to run the server on
	// Use PORT from environment variables, or default to 5000
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	// Start the web server listening on all network interfaces
	app.Listen("0.0.0.0:" + port)
}

// getTodos retrieves all todo items from the database
func getTodos(c *fiber.Ctx) error {
	// Slice to store all retrieved todos
	var todos []Todo

	// Find all documents in the collection
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return err
	}

	// Ensure the cursor is closed after we're done with it
	defer cursor.Close(context.Background())

	// Iterate through all documents
	for cursor.Next(context.Background()) {
		var todo Todo
		// Decode each document into a Todo struct
		if err := cursor.Decode(&todo); err != nil {
			return err
		}

		// Add the decoded todo to the slice
		todos = append(todos, todo)
	}

	// Return the list of todos as a JSON response
	return c.JSON(todos)
}

// createTodo handles the creation of a new todo item
func createTodo(c *fiber.Ctx) error {
	// Create a new Todo struct to parse the request body into
	todo := new(Todo)

	// Parse the request body into the todo struct
	if err := c.BodyParser(todo); err != nil {
		return err
	}

	// Validate that the todo body is not empty
	if todo.Body == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Todo body cannot be empty"})
	}

	// Insert the new todo into the database
	insertResult, err := collection.InsertOne(context.Background(), todo)
	if err != nil {
		return err
	}

	// Set the ID of the todo to the newly inserted document's ID
	todo.ID = insertResult.InsertedID.(primitive.ObjectID)

	// Return the created todo with a 201 Created status
	return c.Status(201).JSON(todo)
}

// updateTodo handles marking a todo as completed
func updateTodo(c *fiber.Ctx) error {
	// Get the todo ID from the URL parameter
	id := c.Params("id")

	// Convert the string ID to a MongoDB ObjectID
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid todo ID"})
	}

	// Create a filter to find the specific todo
	filter := bson.M{"_id": objectId}
	// Create an update operation to set completed to true
	update := bson.M{"$set": bson.M{"completed": true}}

	// Perform the update in the database
	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	// Return a success response
	return c.Status(200).JSON(fiber.Map{"success": true})
}

// deleteTodo handles the deletion of a todo item
func deleteTodo(c *fiber.Ctx) error {
	// Get the todo ID from the URL parameter
	id := c.Params("id")

	// Convert the string ID to a MongoDB ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid todo ID"})
	}

	// Create a filter to find the specific todo
	filter := bson.M{"_id": objectID}

	// Delete the todo from the database
	_, err = collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}

	// Return a success response
	return c.Status(200).JSON(fiber.Map{"success": true})
}
