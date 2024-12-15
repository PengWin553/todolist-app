// Package declaration for a standalone executable
package main

// Import necessary packages
import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors" // Add CORS middleware
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Todo struct remains the same
type Todo struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Completed bool               `json:"completed"`
	Body      string             `json:"body"`
}

// Global variable for MongoDB collection
var collection *mongo.Collection

func main() {
	fmt.Println("Hello, world!!!")

	// Load environment variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	// Retrieve MongoDB connection URI
	MONGODB_URI := os.Getenv("MONGODB_URI")

	// Create client options
	clientOptions := options.Client().ApplyURI(MONGODB_URI)

	// Establish database connection
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Ensure database connection is closed
	defer client.Disconnect(context.Background())

	// Verify the database connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MONGODB Atlas")

	// Set up the collection
	collection = client.Database("todolist_app_db").Collection("todos")

	// Create a new Fiber web application instance
	app := fiber.New()

	// Add CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173", // Your frontend URL
		AllowMethods:     "GET,POST,PATCH,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept",
		AllowCredentials: true,
	}))

	// Define API routes
	app.Get("/api/todos", getTodos)
	app.Post("/api/todos", createTodo)
	app.Patch("/api/todos/:id", updateTodo)
	app.Delete("/api/todos/:id", deleteTodo)

	// Determine the port to run the server on
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	// Start the web server
	fmt.Printf("Server starting on port %s\n", port)
	log.Fatal(app.Listen("0.0.0.0:" + port))
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
