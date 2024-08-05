package database

import (
	"fmt"

	"Tg_chatbot/models"

	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Define database interface (for DI)
type Database interface {
	Create(value interface{}) error
	Where(query interface{}, args ...interface{}) Database
	First(out interface{}, where ...interface{}) error
	Save(value interface{}) error
	Model(value interface{}) Database
	Take(out interface{}, where ...interface{}) error
	Delete(value interface{}, where ...interface{}) error
	Find(out interface{}, where ...interface{}) error
	Updates(values interface{}) error
}

// Initialize the database connection (Postgres)
func InitPostgresDB(dbstr string) {
	if DB != nil {
		return // Already initialized
	}

	// Specify the db details (via .env file)
	// use input
	//dbstr := os.Getenv("DATABASE_URL")

	// Connect to the database with gorm (postgres)
	db, err := gorm.Open(postgres.Open(dbstr), &gorm.Config{})
	if err != nil {
		panic("Error: Failed to connect to database!")
	}

	// Auto migrate the User and Item schemas
	// Write in two use for array
	if err := db.AutoMigrate(&models.User{}); err != nil {
		panic("Error: Migration failed!")
	}

	// Assign the database connection to a global variable
	DB = db
	fmt.Println("Database connected!")
}

var Client *mongo.Client
var ItemCollection *mongo.Collection

// Initialize the database connection (MongoDB)
/*
func InitMongoDB() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	Client = client
	ItemCollection = client.Database("testDB").Collection("items")
	log.Println("Connected to MongoDB!")

	return client, nil
}*/
