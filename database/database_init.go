package database

import (
	"fmt"

	config "Tg_chatbot/configs"
	"Tg_chatbot/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//var DB *gorm.DB

// Database2 interface defines the methods for database operations
type Database2 interface {
	Init() error
	GetDB() *gorm.DB
}

// database2 struct holds the connection details and the gorm DB instance
type database2 struct {
	conf *config.Config
	//user string
	//pwd  string
	db *gorm.DB
}

// NewDatabase2 creates a new instance of database2 with the provided config
func NewDatabase2(config *config.Config) Database2 {
	return &database2{
		conf: config,
		//user: config.GetDBUser(),
		//pwd:  config.GetDBPwd(),
	}
}

// Init initializes the database connection and performs migrations
func (db2 *database2) Init() error {
	//dbstr := fmt.Sprintf("host=localhost user=%s password=%s dbname=chatbot port=5432 sslmode=disable", db2.user, db2.pwd)
	dbstr := db2.conf.GetDBString()

	db, err := gorm.Open(postgres.Open(dbstr), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto migrate the User schema
	if err := db.AutoMigrate(&models.User{}); err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	// Assign the initialized database connection to the db field
	db2.db = db
	fmt.Println("Database connected!")

	return nil
}

// GetDB returns the gorm DB instance
func (db2 *database2) GetDB() *gorm.DB {
	return db2.db
}

// Initialize the database connection (Postgres)
/*func InitPostgresDB(dbstr string) {
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
}*/

//var Client *mongo.Client
//var ItemCollection *mongo.Collection

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
