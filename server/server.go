package server

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Server struct {
	cfg Config
	db  *gorm.DB
}

func New(cfg Config, db *gorm.DB) *Server {
	return &Server{cfg: cfg,
		db: db}
}

/*func New(cfg Config) *Server {
	return &Server{cfg: cfg}
}*/

func (s *Server) Start() error {

	fmt.Println("Initializing server routes")
	router := gin.Default()

	// Initialize routes
	//InitRoutes(router)
	InitRoutes(router, s.db)

	fmt.Println("Starting the server on port", s.cfg.Port)
	//return router.Run(s.cfg.Host + ":" + string(s.cfg.Port))
	return router.Run("0.0.0.0:" + strconv.Itoa(s.cfg.Port)) // Binding to 0.0.0.0
}
