package Http

import (
	"fmt"
	"github.com/7058011439/haoqbb/Log"
	"github.com/7058011439/haoqbb/Timer"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"sync"
	"testing"
)

func TestGetHttpSync(t *testing.T) {
	fmt.Println(os.Getenv("HTTP_PROXY"))
	timer := Timer.NewCountAddUp(Timer.Millisecond, "http访问")
	group := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		group.Add(1)
		go func() {
			for index := 0; index < 100000; index++ {
				if _, err := GetHttpSync("http://39.101.212.11:8888/api/forum/card/card-info?id=2", NewHead(nil).Add("token", "xE.YHURaksbUuoJ.9qxlWDMhtChAtuO0HqyKUISr3s.GqvmfNlZt0dOpnoAS62oH")); err != nil {
					Log.ErrorLog("http 访问失败, err = %v", err.Error())
					break
				}
				timer.PrintCost(1000, "")
				//time.Sleep(time.Millisecond * 1)
			}
			group.Done()
		}()
	}

	group.Wait()
	//data, _ := GetHttpSync("https://fapi.binance.com/fapi/v1/ticker/24hr", nil)
	//fmt.Println(string(data))
}

func TestGetHttpSyncProxy(t *testing.T) {
	data, _ := GetHttpSync("https://fapi.binance.com/fapi/v1/ticker/24hr", nil)
	fmt.Println(string(data))
}

type testApi struct {
	Api
}

func (t *testApi) login(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "登录")
}

func (t *testApi) regedit(ctx *gin.Context) {
	ctx.String(http.StatusOK, "注册")
}

func TestNewHttpServer(t *testing.T) {
	server := NewHttpServer(ServerModeDebug)
	account := server.RegeditGroup("/account", &testApi{}).(*testApi)

	account.RegeditApi(TypeGet, "/login", account.login)
	account.RegeditApi(TypeGet, "/regedit", account.regedit)

	server.Start(1000)

	select {}
}
