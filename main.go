// Package declaration for a standalone executable
package main

// Import necessary packages
import (
	"context" // Provides support for managing request-scoped values, cancellation signals, and deadlines
	"fmt"     // Implements formatted I/O functions for printing
	"log"     // Provides logging functionality
	"os"      // Provides a platform-independent interface to operating system functionality

	"github.com/gofiber/fiber/v2" // Corrected import to use fiber/v2
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
	ID primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`

	// Completed indicates whether the todo item has been finished
	// Boolean field that can be true (completed) or false (not completed)
	Completed bool `json:"completed"`

	// Body contains the text description of the todo item
	Body string `json:"body"`
}

// Declare a global variable to hold the MongoDB collection
// This will be used to interact with the specific collection in the database
var collection *mongo.Collection // Corrected typo from 'collecton' to 'collection'

// main is the entry point of the Go application
func main() {
	// Print a greeting
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

	collection = client.Database("todolist_app_db").Collection("todos")

	app := fiber.New()

	app.Get("/api/todos", getTodos)
	app.Post("/api/todos", createTodo)
	app.Patch("/api/todos/:id", updateTodo)
	app.Delete("/api/todos/:id", deleteTodo)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	app.Listen("0.0.0.0:" + port) // Removed log.Fatal to prevent unreachable code
}

func getTodos(c *fiber.Ctx) error {
	var todos []Todo

	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return err
	}

	defer cursor.Close(context.Background()) // Added to properly close the cursor

	for cursor.Next(context.Background()) {
		var todo Todo
		if err := cursor.Decode(&todo); err != nil {
			return err
		}

		todos = append(todos, todo)
	}

	return c.JSON(todos)
}

func createTodo(c *fiber.Ctx) error {
	todo := new(Todo)
	// {id:0,completed:false,body:""}

	if err := c.BodyParser(todo); err != nil {
		return err
	}

	if todo.Body == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Todo body cannot be empty"})
	}

	insertResult, err := collection.InsertOne(context.Background(), todo)
	if err != nil {
		return err
	}

	todo.ID = insertResult.InsertedID.(primitive.ObjectID)

	return c.Status(201).JSON(todo)

	// c.BodyParser((todo))
}

func updateTodo(c *fiber.Ctx) error {
	id := c.Params("id")
	objectId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid todo ID"})
	}

	filter := bson.M{"_id": objectId}
	update := bson.M{"$set": bson.M{"completed": true}}

	_, err = collection.UpdateOne(context.Background(), filter, update)

	if err != nil {
		return err
	}

	return c.Status(200).JSON(fiber.Map{"success": true})

}

func deleteTodo(c *fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid todo ID"})
	}

	filter := bson.M{"_id": objectID}
	_, err = collection.DeleteOne(context.Background(), filter)

	if err != nil {
		return err
	}

	return c.Status(200).JSON(fiber.Map{"success": true})
}
