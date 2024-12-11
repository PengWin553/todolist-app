// Package declaration for a standalone executable
package main

// Import necessary packages
import (
	"context" // Provides support for managing request-scoped values, cancellation signals, and deadlines
	"fmt"     // Implements formatted I/O functions for printing
	"log"     // Provides logging functionality
	"os"      // Provides a platform-independent interface to operating system functionality

	"github.com/joho/godotenv"                  // Loads environment variables from a .env file
	"go.mongodb.org/mongo-driver/mongo"         // MongoDB Go driver for database operations
	"go.mongodb.org/mongo-driver/mongo/options" // Provides configuration options for MongoDB client
)

// Todo struct defines the structure of a todo item in the application
// This represents how a todo will be stored in the MongoDB database
type Todo struct {
	// ID is the unique identifier for each Todo item
	// `json:"id"` defines how the field is serialized to JSON
	// `bson:"_id"` specifies how the field is stored in MongoDB (MongoDB uses '_id' as the default ID field)
	ID int `json:"id" bson:"_id"`

	// Completed indicates whether the todo item has been finished
	// Boolean field that can be true (completed) or false (not completed)
	Completed bool `json:"completed"`

	// Body contains the text description of the todo item
	Body string `json:"body"`
}

// Declare a global variable to hold the MongoDB collection
// This will be used to interact with the specific collection in the database
var collecton *mongo.Collection // Note: there's a typo here, it should be 'collection'

// main is the entry point of the Go application
func main() {
	// Print a greeting (slight typo in the original - should be "Hello" instead of "yellow")
	fmt.Println("yellow, world")

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

	// Verify the database connection by pinging the server
	// This ensures that the connection is active and working
	err = client.Ping(context.Background(), nil)
	if err != nil {
		// If ping fails, terminate the program
		log.Fatal(err)
	}

	// Print a success message when connection is established
	fmt.Println("Connected to MONGODB Atlas")
}
