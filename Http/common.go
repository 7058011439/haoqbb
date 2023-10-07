package Http

import (
	"github.com/7058011439/haoqbb/Log"
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewResult(c *gin.Context) *WebResult {
	return &WebResult{
		c: c,
	}
}

func Bind(c *gin.Context, data interface{}) bool {
	if err := c.ShouldBind(data); err == nil {
		return true
	} else {
		c.JSON(http.StatusOK, &WebResult{
			Msg:  err.Error(),
			Code: http.StatusInternalServerError,
		})
		Log.ErrorLog(err.Error())
		return false
	}
}

type WebResult struct {
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data"`
	Code  int         `json:"code"`
	c     *gin.Context
	abort bool
}

func (w *WebResult) result(msg string, data interface{}, code int) {
	w.Msg = msg
	w.Code = code
	if w.Data == nil && data != nil {
		w.Data = data
	}
	w.c.JSON(http.StatusOK, w)
}

func (w *WebResult) Fail(msg string, data interface{}) {
	w.result(msg, data, http.StatusInternalServerError)
}

func (w *WebResult) Success(msg string, data interface{}) {
	w.result(msg, data, http.StatusOK)
}

func (w *WebResult) AddData(key string, value interface{}) {
	if w.Data == nil {
		w.Data = make(map[string]interface{}, 4)
	}
	w.Data.(map[string]interface{})[key] = value
}

func (w *WebResult) Abort(msg string, data interface{}) {
	w.abort = true
	w.Fail(msg, data)
	w.c.Abort()
}

func (w *WebResult) IsAbort() bool {
	return w.abort
}
