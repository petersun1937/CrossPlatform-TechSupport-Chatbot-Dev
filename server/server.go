package server

import (
	"fmt"
	"strconv"

	config "Tg_chatbot/configs"
	"Tg_chatbot/service"

	"github.com/gin-gonic/gin"
)

type Server struct {
	svrcfg config.ServerConfig
	srv    *service.Service
	conf   *config.Config
}

func New(svrcfg config.ServerConfig, srv *service.Service, conf *config.Config) *Server {
	return &Server{
		svrcfg: svrcfg,
		srv:    srv,
		conf:   conf,
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
func (s *Server) Start() error {

	fmt.Println("Initializing server routes")
	router := gin.Default()

	// Initialize routes
	InitRoutes(router, s.conf, s.srv)

	// Run the routes
	fmt.Println("Starting the server on port", s.svrcfg.Port)
	err := router.Run("0.0.0.0:" + strconv.Itoa(s.svrcfg.Port)) // Binding to 0.0.0.0
	//err := router.Run("0.0.0.0:8080") // Binding to 0.0.0.0
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}
	fmt.Println("Server started and running...")
	return nil
}
