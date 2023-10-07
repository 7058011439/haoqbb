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

func (s *Server) Start(port int) {
	go s.Run(fmt.Sprintf(":%v", port))
}

func (s *Server) RegeditGroup(baseUrl string, api IApi, middleware ...gin.HandlerFunc) IApi {
	ret := &Api{
		RouterGroup: s.Group(baseUrl, middleware...),
	}
	api.setApi(ret)
	return api
}

func (s *Server) RegeditApi(reType string, uri string, fun func(c *gin.Context), middleware ...gin.HandlerFunc) error {
	middleware = append(middleware, fun)
	switch reType {
	case TypePost:
		s.POST(uri, middleware...)
	case TypeGet:
		s.GET(uri, middleware...)
	default:
		return fmt.Errorf("unknown type error")
	}
	return nil
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
