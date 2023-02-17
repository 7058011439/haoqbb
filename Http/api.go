package Http

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

const (
	TypePost = "post"
	TypeGet  = "get"
)

type Api struct {
	BaseUrl string
	server  *Server
	group   *gin.RouterGroup
}

func (a *Api) Regedit(reType string, uri string, fun func(c *gin.Context), middleware ...gin.HandlerFunc) error {
	if a.server != nil {
		middleware = append(middleware, fun)
		switch reType {

		case TypePost:
			a.group.POST(uri, middleware...)
		case TypeGet:
			a.group.GET(uri, middleware...)
		default:
			return errors.Errorf("unknown type error")
		}
	}
	return errors.Errorf("engine is nil")
}

func NewApi(server *Server, baseUrl string, middleware ...gin.HandlerFunc) *Api {
	ret := &Api{
		BaseUrl: baseUrl,
		server:  server,
		group:   server.Group(baseUrl, middleware...),
	}
	return ret
}
