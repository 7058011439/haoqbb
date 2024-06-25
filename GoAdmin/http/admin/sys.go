package admin

import (
	"github.com/7058011439/haoqbb/GoAdmin/http/common"
	"github.com/7058011439/haoqbb/Http"
	"github.com/7058011439/haoqbb/Net"
	"github.com/7058011439/haoqbb/System"
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"os"
	"runtime"
	"time"
)

const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
)

type apiSys struct {
	Http.Api
}

func init() {
	a := common.ServerAdmin.RegeditGroup("/api/sys", &apiSys{}, common.CheckAdminToken).(*apiSys)
	a.RegeditApi(Http.TypeGet, "/server-monitor", a.serverMonitor)
}

func (a *apiSys) serverMonitor(c *gin.Context) {
	ret := Http.NewResult(c)

	osData := map[string]interface{}{}
	osData["goOs"] = runtime.GOOS
	osData["arch"] = runtime.GOARCH
	osData["mem"] = runtime.MemProfileRate
	osData["compiler"] = runtime.Compiler
	osData["version"] = runtime.Version()
	osData["numGoroutine"] = runtime.NumGoroutine()
	osData["ip"] = Net.GetInputBoundIP()
	osData["hostName"], _ = os.Hostname()
	osData["time"] = time.Now().Format("2006-01-02 15:04:05")

	cpuData := map[string]interface{}{}
	cpuData["cpuInfo"], _ = cpu.Info()
	cpuData["percent"] = System.GetCpuPercent()
	cpuData["cpuNum"], _ = cpu.Counts(false)

	memData := map[string]interface{}{}
	m, _ := mem.VirtualMemory()
	memData["total"] = m.Total / MB
	memData["used"] = m.Used / MB
	memData["percent"] = m.UsedPercent

	swapData := map[string]interface{}{}
	swapData["used"] = m.SwapTotal - m.SwapFree
	swapData["total"] = m.SwapTotal

	diskData := map[string]interface{}{}
	d, _ := disk.Usage("/")
	diskData["total"] = float64(d.Total) / GB
	diskData["used"] = float64(d.Used) / GB
	diskData["percent"] = float64(d.Used) / float64(d.Total) * 100

	netData := map[string]interface{}{}
	oldStatus, _ := net.IOCounters(false)
	currStatus, _ := net.IOCounters(false)
	netData["in"] = float64(currStatus[0].BytesRecv - oldStatus[0].BytesRecv)
	netData["out"] = float64(currStatus[0].BytesSent - oldStatus[0].BytesSent)

	bootTime, _ := host.BootTime()

	ret.Success(common.ResponseSuccess, map[string]interface{}{
		"os":       osData,
		"cpu":      cpuData,
		"mem":      memData,
		"swap":     swapData,
		"disk":     diskData,
		"net":      netData,
		"location": "haoqbb",
		"bootTime": time.Since(time.Unix(int64(bootTime), 0)).Hours(),
	})
}
