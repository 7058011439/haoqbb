package common

import (
	"bytes"
	"github.com/7058011439/haoqbb/GoAdmin/db/admin"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"time"
)

// 定义一个自定义的ResponseWriter
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// 重写Write方法来捕获响应内容
func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.body.Write(b)
	return rw.ResponseWriter.Write(b)
}

func insertOperateLog(c *gin.Context, path string) {
	permissionMutex.RLock()
	if _, ok := logApi[path]; !ok {
		permissionMutex.RUnlock()
		return
	}
	permissionMutex.RUnlock()
	start := time.Now()
	rw := &responseWriter{
		ResponseWriter: c.Writer,
		body:           bytes.NewBufferString(""),
	}
	c.Writer = rw
	param, _ := ioutil.ReadAll(c.Request.Body)
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(param))
	c.Next()

	if len(param) < 1 {
		param = []byte(c.Request.URL.RawQuery)
	}
	ret := rw.body.String()
	/*
		if len(param) > 250 {
			param = param[:250]
			param = append(param, "..."...)
		}
		if len(ret) > 250 {
			ret = ret[:250]
			ret += "..."
		}
	*/

	admin.InsertForumData(&admin.OperateLog{
		UserName:    GetString(c, TokenKeyAdminUserName),
		Method:      c.Request.Method,
		Url:         c.Request.URL.RequestURI(),
		FullPath:    c.FullPath(),
		Ip:          c.RemoteIP(),
		Param:       string(param),
		LatencyTime: time.Since(start).Seconds(),
		Result:      ret,
	})
}
