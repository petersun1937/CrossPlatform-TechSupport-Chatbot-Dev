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

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (app *App) InitRoutes(r *gin.Engine, conf *config.Config, srv *service.Service) {
	// Set up logging to a file (bot.log)
	file, err := os.OpenFile("bot.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	gin.DefaultWriter = file

	// Enable CORS
	app.Router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Or use "*" to allow all origins
		AllowMethods:     []string{"POST", "GET", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Middleware to log requests
	app.Router.Use(gin.Logger())
	app.Router.Use(gin.Recovery())

	// Define routes
	app.Router.POST("/webhook/line", func(c *gin.Context) {
		handlers.HandleLineWebhook(c, app.LineBot)
	})
	app.Router.POST("/webhook/telegram", func(c *gin.Context) {
		handlers.HandleTelegramWebhook(c, app.TgBot)
	})
	app.Router.GET("/messenger/webhook", handlers.VerifyMessengerWebhook) // For webhook verification
	app.Router.POST("/messenger/webhook", func(c *gin.Context) {
		handlers.HandleMessengerWebhook(c, app.FbBot)
	})
	app.Router.POST("/api/message", func(c *gin.Context) {
		handlers.HandlerGeneralBot(c, app.GeneralBot) // Pass the generalBot instance here
	})

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

func (app *App) RunRoutes(conf *config.Config, svc *service.Service) {

	cfg := config.ServerConfig{
		Host:    os.Getenv("HOST"),
		Port:    8080, // Default port, can be overridden
		Timeout: 30 * time.Second,
		MaxConn: 100,
	}

	//if p := os.Getenv("APP_PORT"); p != "" {
	if p := app.Config.AppPort; p != "" {
		pInt, err := strconv.Atoi(p)
		if err == nil {
			cfg.Port = pInt
		}
	}

	//db := svc.GetDB() // Get the initialized DB instance

	//fmt.Println("Starting the server on port " + strconv.Itoa(cfg.Port))
	svr := New(cfg, svc, conf)
	//svr := New(cfg)
	if err := svr.Start(app); err != nil {
		log.Fatal(err)
	}

}
