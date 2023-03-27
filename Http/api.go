package Http

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

const (
	TypePost = "post"
	TypeGet  = "get"
)

type IApi interface {
	setApi(api *Api)
}

type Api struct {
	*gin.RouterGroup
}

func (a *Api) RegeditApi(reType string, uri string, fun func(c *gin.Context), middleware ...gin.HandlerFunc) error {
	middleware = append(middleware, fun)
	switch reType {
	case TypePost:
		a.RouterGroup.POST(uri, middleware...)
	case TypeGet:
		a.RouterGroup.GET(uri, middleware...)
	default:
		return fmt.Errorf("unknown type error")
	}
	return nil
}

func (a *Api) setApi(api *Api) {
	a.RouterGroup = api.RouterGroup
}
