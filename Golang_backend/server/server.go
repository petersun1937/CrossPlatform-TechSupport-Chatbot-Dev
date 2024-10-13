package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	config "crossplatform_chatbot/configs"
	"crossplatform_chatbot/handlers"

	"github.com/gin-gonic/gin"
)

type Server struct {
	svrcfg  config.ServerConfig
	handler *handlers.Handler
	// srv     *service.Service
	// conf   *config.Config
	router *gin.Engine
}

func New(svrcfg config.ServerConfig, handler *handlers.Handler) *Server {
	return &Server{
		svrcfg: svrcfg,
		// srv:     srv,
		handler: handler,
		// conf:   conf,
	}
}

/*
func (s *Server) Start() error {

		fmt.Println("Initializing server routes")
		router := gin.Default()

		// Initialize routes
		InitRoutes(router, s.conf, s.srv)
		//InitRoutes(router, s.db)

		// Run the routes
		fmt.Println("Starting the server on port", s.svrcfg.Port)
		err := router.Run("0.0.0.0:" + strconv.Itoa(s.svrcfg.Port)) // Binding to 0.0.0.0
		//return router.Run("0.0.0.0:" + strconv.Itoa(s.svrcfg.Port)) // Binding to 0.0.0.0
		if err != nil {
			return fmt.Errorf("failed to start server: %w", err)
		}
		fmt.Println("Server started and running...")
		return nil
	}
*/
// func (s *Server) Start(app *App) error {

// 	fmt.Println("Initializing server routes")
// 	//router := gin.Default()

// 	// Initialize routes
// 	app.InitRoutes(app.Router, s.conf, s.srv)

// 	// Run the routes
// 	fmt.Println("Starting the server on port", s.svrcfg.Port)
// 	err := app.Router.Run("0.0.0.0:" + strconv.Itoa(s.svrcfg.Port)) // Binding to 0.0.0.0
// 	//err := router.Run("0.0.0.0:8080") // Binding to 0.0.0.0
// 	if err != nil {
// 		return fmt.Errorf("failed to start server: %w", err)
// 	}
// 	fmt.Println("Server started and running...")
// 	return nil
// }

func (s *Server) V2Start() {
	// Run Gin Server
	fmt.Println("Initializing server routes")
	//router := gin.Default()

	// Initialize routes
	s.V2InitRoutes(s.handler)

	// Run the routes
	fmt.Println("Starting the server on port", s.svrcfg.Port)
	err := s.router.Run("0.0.0.0:" + strconv.Itoa(s.svrcfg.Port)) // Binding to 0.0.0.0
	//err := router.Run("0.0.0.0:8080") // Binding to 0.0.0.0
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
	fmt.Println("Server started and running...")

	// Init Service
	if err := s.handler.Service.Init(); err != nil {
		log.Fatal("Failed to initialize service:", err)
	}
}

// Gracefully shuts down the HTTP server.
func Shutdown(ctx context.Context) error {
	server := &http.Server{
		Addr: ":8080", //strconv.Itoa(s.svrcfg.Port),
	}

	// Perform any pre-shutdown tasks, like closing database connections, flushing logs, etc.
	fmt.Println("Performing pre-shutdown tasks...")

	// Attempt to gracefully shut down the server with a timeout context.
	err := server.Shutdown(ctx)
	if err != nil {
		fmt.Printf("Server Shutdown Failed:%+v", err)
		return err
	}

	// Optionally, log or perform any post-shutdown tasks.
	fmt.Println("Server gracefully stopped")
	return nil
}
