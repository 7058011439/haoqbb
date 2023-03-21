package System

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"net"
	"strings"
	"time"
)

func GetCpuPercent() float64 {
	percent, _ := cpu.Percent(time.Microsecond, false)
	return percent[0]
}

func GetMemPercent() float64 {
	memInfo, _ := mem.VirtualMemory()
	return memInfo.UsedPercent
}

func GetDiskPercent() float64 {
	parts, _ := disk.Partitions(true)
	diskInfo, _ := disk.Usage(parts[0].Mountpoint)
	return diskInfo.UsedPercent
}

// GetMachineInfo  获取机器码
func GetMachineInfo() map[string]string {
	res := make(map[string]string)
	res, err := getMACAddress()
	if err != nil {
		panic(err.Error())
	}
	if cpus, err := cpu.Info(); err != nil && len(cpus) > 0 {
		panic(err.Error())
	} else {
		res["cpu"] = cpus[0].VendorID + "_" + cpus[0].PhysicalID
	}

	return res
}

func GetMachineId() string {
	machineInfo := GetMachineInfo()
	byteInfo, _ := json.Marshal(machineInfo)
	machineCode := getMd5String(string(byteInfo), true, true)
	ret := ""
	for i := 0; i < len(machineCode); {
		ret += machineCode[i : i+4]
		i += 4
		if i < len(machineCode) {
			ret += "-"
		}
	}

	return ret
}

// getMACAddress 获取网卡MAC地址
func getMACAddress() (map[string]string, error) {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		panic(err.Error())
	}
	mac, macErr := "", errors.New("无法获取到正确的MAC地址")
	for i := 0; i < len(netInterfaces); i++ {
		//fmt.Println(netInterfaces[i])
		if (netInterfaces[i].Flags&net.FlagUp) != 0 && (netInterfaces[i].Flags&net.FlagLoopback) == 0 {
			adds, _ := netInterfaces[i].Addrs()
			for _, address := range adds {
				ipnet, ok := address.(*net.IPNet)
				if ok && ipnet.IP.IsGlobalUnicast() {
					// 如果IP是全局单拨地址，则返回MAC地址
					mac = netInterfaces[i].HardwareAddr.String()
					result := map[string]string{
						"mac":  mac,
						"ipv4": ipnet.IP.To4().String(),
						"ipv6": ipnet.IP.To16().String(),
					}
					//fmt.Println(result)
					return result, nil
				}
			}
		}
	}
	result := map[string]string{
		"mac":  mac,
		"ipv4": "",
		"ipv6": "",
	}
	return result, macErr
}

// getMd5String 生成32位md5字串
func getMd5String(s string, upper bool, half bool) string {
	h := md5.New()
	h.Write([]byte(s))
	result := hex.EncodeToString(h.Sum(nil))
	if upper == true {
		result = strings.ToUpper(result)
	}
	if half == true {
		result = result[8:24]
	}
	return result
}
