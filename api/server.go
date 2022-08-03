package api

import (
	"hamza72x/bankify/config"
	db "hamza72x/bankify/db/sqlc"
	"os"

	"github.com/gin-gonic/gin"
)

type Server struct {
	cfg    config.Config
	router *gin.Engine
	store  *db.Store
}

func New(cfg config.Config, store *db.Store) (*Server, error) {

	s := &Server{
		cfg:   cfg,
		store: store,
	}

	s.router = s.setupRoutes()

	return s, nil
}

func (s *Server) Start() error {
	return s.router.Run(":" + getPort())
}

func getPort() string {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "3001"
	}
	return port
}

// TODO:- format error according to the error type
func errResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
