package server

import (
	"Tg_chatbot/database"
	"Tg_chatbot/handlers"
	"Tg_chatbot/middleware"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitRoutes(r *gin.Engine, db *gorm.DB) {
	// Set up logging to a file
	file, err := os.OpenFile("bot.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	gin.DefaultWriter = file

	// Middleware to log requests
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Define routes
	r.POST("/webhook", handlers.HandleLineWebhook)
	//r.POST("/webhook/telegram", handlers.HandleTelegramWebhook)
	//r.POST("/login", handlers.Login)

	// Protected routes
	authorized := r.Group("/api")
	authorized.Use(middleware.JWTMiddleware())
	{
		authorized.POST("/message", handlers.HandleCustomMessage)
		// Add other protected routes here
	}

	fmt.Println("Server started")
	r.Run(":8080")
}

func RunRoutes() {

	cfg := Config{
		Host:    os.Getenv("HOST"),
		Port:    8080, // Default port, can be overridden
		Timeout: 30 * time.Second,
		MaxConn: 100,
	}

	if p := os.Getenv("APP_PORT"); p != "" {

		pInt, err := strconv.Atoi(p)
		if err == nil {
			cfg.Port = pInt
		}
	}

	db := database.DB // Get the initialized DB instance

	fmt.Println("Starting the server on port " + strconv.Itoa(cfg.Port))
	svr := New(cfg, db)
	//svr := New(cfg)
	if err := svr.Start(); err != nil {
		log.Fatal(err)
	}

}
