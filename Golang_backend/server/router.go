package server

import (
	"crossplatform_chatbot/handlers"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (s *Server) V2InitRoutes(handler *handlers.Handler) {
	// Set up logging to a file (bot.log)
	file, err := os.OpenFile("bot.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	gin.DefaultWriter = file

	// Enable CORS
	s.router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Or use "*" to allow all origins
		AllowMethods:     []string{"POST", "GET", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Middleware to log requests
	s.router.Use(gin.Logger())
	s.router.Use(gin.Recovery())

	// Define routes
	s.router.POST("/webhook/line", handler.V2HandleLineWebhook)
	// s.router.POST("/webhook/line", func(c *gin.Context) {
	// 	handlers.HandleLineWebhook(c, app.Bots["line"].(bot.LineBot))
	// })
	// s.router.POST("/webhook/telegram", func(c *gin.Context) {
	// 	handlers.HandleTelegramWebhook(c, app.Bots["tg"].(bot.TgBot))
	// })
	// s.router.GET("/messenger/webhook", handlers.VerifyMessengerWebhook) // For webhook verification
	// s.router.POST("/messenger/webhook", func(c *gin.Context) {
	// 	handlers.HandleMessengerWebhook(c, app.Bots["fb"].(bot.FbBot))
	// })
	// // Instagram webhook for verification and message handling
	// s.router.GET("/instagram/webhook", handlers.VerifyInstagramWebhook) // For webhook verification
	// s.router.POST("/instagram/webhook", func(c *gin.Context) {
	// 	handlers.HandleInstagramWebhook(c, app.Bots["ig"].(bot.IgBot))
	// })
	// s.router.POST("/api/message", func(c *gin.Context) {
	// 	handlers.HandlerGeneralBot(c, app.Bots["general"].(bot.GeneralBot)) // Pass the generalBot instance here
	// })
	// s.router.POST("/api/document/upload", func(c *gin.Context) {
	// 	handlers.HandlerDocumentUpload(c, app.Bots["general"].(bot.GeneralBot)) // Bind the document upload handler to the route
	// })

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

// func (app *App) InitRoutes(r *gin.Engine, conf *config.Config, srv *service.Service) {
// 	// Set up logging to a file (bot.log)
// 	file, err := os.OpenFile("bot.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	gin.DefaultWriter = file

// 	// Enable CORS
// 	app.Router.Use(cors.New(cors.Config{
// 		AllowOrigins:     []string{"http://localhost:3000"}, // Or use "*" to allow all origins
// 		AllowMethods:     []string{"POST", "GET", "OPTIONS"},
// 		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
// 		ExposeHeaders:    []string{"Content-Length"},
// 		AllowCredentials: true,
// 		MaxAge:           12 * time.Hour,
// 	}))

// 	// Middleware to log requests
// 	app.Router.Use(gin.Logger())
// 	app.Router.Use(gin.Recovery())

// 	// Define routes
// 	app.Router.POST("/webhook/line", handlers.V2HandleLineWebhook)
// 	// app.Router.POST("/webhook/line", func(c *gin.Context) {
// 	// 	handlers.HandleLineWebhook(c, app.Bots["line"].(bot.LineBot))
// 	// })
// 	app.Router.POST("/webhook/telegram", func(c *gin.Context) {
// 		handlers.HandleTelegramWebhook(c, app.Bots["tg"].(bot.TgBot))
// 	})
// 	app.Router.GET("/messenger/webhook", handlers.VerifyMessengerWebhook) // For webhook verification
// 	app.Router.POST("/messenger/webhook", func(c *gin.Context) {
// 		handlers.HandleMessengerWebhook(c, app.Bots["fb"].(bot.FbBot))
// 	})
// 	// Instagram webhook for verification and message handling
// 	app.Router.GET("/instagram/webhook", handlers.VerifyInstagramWebhook) // For webhook verification
// 	app.Router.POST("/instagram/webhook", func(c *gin.Context) {
// 		handlers.HandleInstagramWebhook(c, app.Bots["ig"].(bot.IgBot))
// 	})
// 	app.Router.POST("/api/message", func(c *gin.Context) {
// 		handlers.HandlerGeneralBot(c, app.Bots["general"].(bot.GeneralBot)) // Pass the generalBot instance here
// 	})
// 	app.Router.POST("/api/document/upload", func(c *gin.Context) {
// 		handlers.HandlerDocumentUpload(c, app.Bots["general"].(bot.GeneralBot)) // Bind the document upload handler to the route
// 	})

// 	//r.POST("/login", handlers.Login)

// 	// Protected routes
// 	/*authorized := r.Group("/api")
// 	authorized.Use(middleware.JWTMiddleware())
// 	{
// 		authorized.POST("/message", handlers.HandleCustomMessage)
// 		// Add other protected routes here
// 	}*/

// 	fmt.Println("Server routes initialized")
// 	//fmt.Println("Server started")
// 	//r.Run(":8080")
// }

// func (app *App) RunRoutes(conf *config.Config, svc *service.Service) {

// 	cfg := config.ServerConfig{
// 		Host:    os.Getenv("HOST"),
// 		Port:    8080, // Default port, can be overridden
// 		Timeout: 30 * time.Second,
// 		MaxConn: 100,
// 	}

// 	//if p := os.Getenv("APP_PORT"); p != "" {
// 	if p := app.Config.AppPort; p != "" {
// 		pInt, err := strconv.Atoi(p)
// 		if err == nil {
// 			cfg.Port = pInt
// 		}
// 	}

// 	//db := svc.GetDB() // Get the initialized DB instance

// 	//fmt.Println("Starting the server on port " + strconv.Itoa(cfg.Port))
// 	svr := New(cfg, svc, conf)
// 	//svr := New(cfg)
// 	if err := svr.Start(app); err != nil {
// 		log.Fatal(err)
// 	}

// }
