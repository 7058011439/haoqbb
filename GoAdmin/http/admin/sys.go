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

type sysRetOs struct {
	GoOs         string `json:"goOs"`         // 系统
	Arch         string `json:"arch"`         // 系统架构
	Mem          int    `json:"mem"`          // 内存分配采样频率
	Compiler     string `json:"compiler"`     // 未知
	Version      string `json:"version"`      // go版本(未使用)
	NumGoroutine int    `json:"numGoroutine"` // 当前活跃的goroutine数量
	Ip           string `json:"ip"`           // Ip(内网)地址
	HostName     string `json:"hostName"`     // 主机名
	Time         string `json:"time"`         // 当前时间(YYYY-MM-DD hh:mm:ss)
}

type sysRetCpu struct {
	CpuInfo []cpu.InfoStat `json:"cpuInfo"` // cpu信息
	Percent float64        `json:"percent"` // cpu使用率(N%)
	CpuNum  int            `json:"cpuNum"`  // cpu核心数
}

type sysRetMem struct {
	Total   uint64  `json:"total"`   // 总共内存
	Used    uint64  `json:"used"`    // 已使用内存
	Percent float64 `json:"percent"` // 内存使用率(N%)
}

type sysRetSwap struct {
	Total uint64 `json:"total"` // 总共内存
	Used  uint64 `json:"used"`  // 已使用内存
}

type sysRetDisk struct {
	Total   float64 `json:"total"`   // 总共硬盘
	Used    float64 `json:"used"`    // 已使用硬盘
	Percent float64 `json:"percent"` // 内存使用率(N%)
}

type sysRetNet struct {
	In  float64 `json:"in"`  // 下载带宽
	Out float64 `json:"out"` // 上传带宽
}

type sysRet struct {
	OS   sysRetOs   `json:"os"`   // 系统信息
	Cpu  sysRetCpu  `json:"cpu"`  // cpu信息
	Mem  sysRetMem  `json:"mem"`  // 内存信息
	Swap sysRetSwap `json:"swap"` // 交换内存
	Disk sysRetDisk `json:"disk"` // 硬盘信息
	Net  sysRetNet  `json:"net"`  // 网络信息

	Location string  `json:"location"`
	BootTime float64 `json:"bootTime"`
}

// @Tags     好奇宝宝后台-系统工具
// @Summary  服务监控
// @Param    token  header    string  true  "token"
// @Success  200    {object}  Http.WebResult{data=sysRet}
// @Router   /api/sys/server-monitor [get]
func (a *apiSys) serverMonitor(c *gin.Context) {
	ret := Http.NewResult(c)

	var retData sysRet

	retData.OS.GoOs = runtime.GOOS
	retData.OS.Arch = runtime.GOARCH
	retData.OS.Mem = runtime.MemProfileRate
	retData.OS.Compiler = runtime.Compiler
	retData.OS.Version = runtime.Version()
	retData.OS.NumGoroutine = runtime.NumGoroutine()
	retData.OS.Ip = Net.GetInputBoundIP()
	retData.OS.HostName, _ = os.Hostname()
	retData.OS.Time = time.Now().Format("2006-01-02 15:04:05")

	retData.Cpu.CpuInfo, _ = cpu.Info()
	retData.Cpu.Percent = System.GetCpuPercent()
	retData.Cpu.CpuNum, _ = cpu.Counts(false)

	m, _ := mem.VirtualMemory()
	retData.Mem.Total = m.Total / MB
	retData.Mem.Used = m.Used / MB
	retData.Mem.Percent = m.UsedPercent

	retData.Swap.Used = m.SwapTotal - m.SwapFree
	retData.Swap.Total = m.SwapTotal

	d, _ := disk.Usage("/")
	retData.Disk.Total = float64(d.Total) / GB
	retData.Disk.Used = float64(d.Used) / GB
	retData.Disk.Percent = float64(d.Used) / float64(d.Total) * 100

	oldStatus, _ := net.IOCounters(false)
	currStatus, _ := net.IOCounters(false)
	retData.Net.In = float64(currStatus[0].BytesRecv - oldStatus[0].BytesRecv)
	retData.Net.Out = float64(currStatus[0].BytesSent - oldStatus[0].BytesSent)

	bootTime, _ := host.BootTime()
	retData.BootTime = time.Since(time.Unix(int64(bootTime), 0)).Hours()
	retData.Location = "haoqbb"

	ret.Success(common.ResponseSuccess, retData)
}
