package api

import (
	"github.com/gin-gonic/gin"
)

func (s *Server) setupRoutes() *gin.Engine {
	r := gin.Default()

	s.registerAccountRoutes(r)
	s.registerTransferRoutes(r)

	return r
}
