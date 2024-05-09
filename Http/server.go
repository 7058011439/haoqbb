package Http

import (
	"fmt"
	"github.com/7058011439/haoqbb/Log"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
	"time"
)

type ServerMode = string

const (
	ServerModeRelease ServerMode = "release"
	ServerModeDebug   ServerMode = "debug"
	ServerModeTest    ServerMode = "test"

	uristatus      = "/status"
	uriresetstatus = "/resetstatus"
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

func NewHttpServer(mode ServerMode, log bool) *Server {
	gin.SetMode(mode)
	en := gin.New()
	if log {
		en.Use(gin.Logger())
	}
	en.Use(gin.Recovery(), statusMiddleware, options)
	en.POST(uristatus, status)
	en.POST(uriresetstatus, resetStatus)
	return &Server{Engine: en}
}

type requestStatus struct {
	requestTimes int64   // 总请求次数
	costTime     float64 // 总处理时间
	maxTime      float64 // 最长耗时
	mutex        sync.Mutex
}

func (s *requestStatus) updateRequestData(t float64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.requestTimes++
	s.costTime += t
	if t > s.maxTime {
		s.maxTime = t
	}
}

func (s *requestStatus) toMap() map[string]interface{} {
	averageTime := 0.0
	if s.requestTimes > 0 {
		averageTime = s.costTime / float64(s.requestTimes)
	}
	return map[string]interface{}{
		"request times    ": s.requestTimes,
		"total cost time  ": fmt.Sprintf("%.3fs", s.costTime),
		"max cost time    ": fmt.Sprintf("%.3fs", s.maxTime),
		"average cost time": fmt.Sprintf("%.3fs", averageTime),
	}
}

var allStatus sync.Map

func statusMiddleware(c *gin.Context) {
	start := time.Now()
	c.Next()
	elapsed := time.Since(start)
	if data, ok := allStatus.Load(c.FullPath()); ok {
		data.(*requestStatus).updateRequestData(elapsed.Seconds())
	} else {
		s := &requestStatus{}
		s.updateRequestData(elapsed.Seconds())
		allStatus.Store(c.FullPath(), s)
	}
}

func status(c *gin.Context) {
	request := struct {
		URI string `json:"uri" form:"uri"`
	}{}
	if err := c.ShouldBind(&request); err == nil {
		if request.URI == "" {
			retData := map[string]interface{}{}
			allStatus.Range(func(key, value interface{}) bool {
				if key.(string) != uristatus && key.(string) != uriresetstatus {
					retData[key.(string)] = value.(*requestStatus).toMap()
				}
				return true
			})
			c.JSON(200, retData)
		} else {
			if data, ok := allStatus.Load(request.URI); ok {
				c.JSON(200, data.(*requestStatus).toMap())
			} else {
				c.JSON(200, (&requestStatus{}).toMap())
			}
		}
	} else {
		c.JSON(http.StatusOK, &WebResult{
			Msg:  err.Error(),
			Code: http.StatusInternalServerError,
		})
		Log.ErrorLog(err.Error())
	}
}

func resetStatus(c *gin.Context) {
	allStatus = sync.Map{}
	c.JSON(200, "ok")
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
