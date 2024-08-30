package server

import (
	"Tg_chatbot/handlers"
	"Tg_chatbot/service"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	config "Tg_chatbot/configs"

	"github.com/gin-gonic/gin"
)

type App struct {
	Config  *config.Config
	Service *service.Service
	Router  *gin.Engine
}

func NewApp(conf *config.Config, srv *service.Service) *App {
	return &App{
		Config:  conf,
		Service: srv,
		Router:  gin.Default(),
	}
}

func InitRoutes(r *gin.Engine, conf *config.Config, srv *service.Service) {
	// Set up logging to a file (bot.log)
	file, err := os.OpenFile("bot.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	gin.DefaultWriter = file

	/*lineBot, err := bot.NewLineBot(conf, srv)
	if err != nil {
		log.Fatal("Failed to initialize LINE bot:", err)
	}*/

	// Middleware to log requests
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Define routes
	r.POST("/webhook", handlers.HandleLineWebhook)
	/*r.POST("/webhook", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})*/
	/*r.POST("/webhook", func(c *gin.Context) {
		handlers.HandleLineWebhook(c, lineBot)
	})*/

	//r.POST("/webhook/telegram", handlers.HandleTelegramWebhook)
	//r.POST("/login", handlers.Login)

	// Protected routes
	/*authorized := r.Group("/api")
	authorized.Use(middleware.JWTMiddleware())
	{
		authorized.POST("/message", handlers.HandleCustomMessage)
		// Add other protected routes here
	}*/

	fmt.Println("Server routes initialized")
	//fmt.Println("Server started")
	//r.Run(":8080")
}

func RunRoutes(conf *config.Config, svc *service.Service) {

	cfg := config.ServerConfig{
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

	//db := svc.GetDB() // Get the initialized DB instance

	//fmt.Println("Starting the server on port " + strconv.Itoa(cfg.Port))
	svr := New(cfg, svc, conf)
	//svr := New(cfg)
	if err := svr.Start(); err != nil {
		log.Fatal(err)
	}

}
