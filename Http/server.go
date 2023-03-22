package Http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ServerMode = string

const (
	ServerModeRelease ServerMode = "release"
	ServerModeDebug   ServerMode = "debug"
	ServerModeTest    ServerMode = "test"
)

type Server struct {
	*gin.Engine
}

func (h *Server) Start(port int) {
	go h.Run(fmt.Sprintf(":%v", port))
}

func NewHttpServer(mode ServerMode) *Server {
	gin.SetMode(mode)
	en := gin.New()
	if mode != ServerModeRelease {
		en.Use(gin.Logger())
	}
	en.Use(gin.Recovery(), options)
	return &Server{Engine: en}
}

func options(c *gin.Context) {
	method := c.Request.Method
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
	c.Header("Access-Control-Allow-Credentials", "true")
	if method == "OPTIONS" {
		c.AbortWithStatus(http.StatusNoContent)
	}
	c.Next()
}
