package System

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/StackExchange/wmi"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"golang.org/x/sys/windows"
	"net"
	"sort"
	"strings"
	"time"
	"unsafe"
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

type cpuInfo struct {
	CPU        int32  `json:"cpu"`
	VendorID   string `json:"vendorId"`
	PhysicalID string `json:"physicalId"`
}

type win32_Processor struct {
	Manufacturer string
	ProcessorID  *string
}

// GetMachineInfo  获取机器码
func GetMachineInfo() map[string]string {
	var ids []string
	if guid, err := getMachineGuid(); err != nil {
		panic(err.Error())
	} else {
		ids = append(ids, guid)
	}
	if cpus, err := getCPUInfo(); err != nil && len(cpus) > 0 {
		panic(err.Error())
	} else {
		ids = append(ids, cpus[0].VendorID+cpus[0].PhysicalID)
	}
	res := make(map[string]string)
	res, err := getMACAddress()
	if err != nil {
		panic(err.Error())
	} else {
		ids = append(ids, res["mac"])
	}
	sort.Strings(ids)
	idsStr := strings.Join(ids, ";")
	machineCode := getMd5String(idsStr, true, true)
	res["machineCode"] = machineCode
	return res
}

func GetMachineId() string {
	machineInfo := GetMachineInfo()
	machineCode := machineInfo["machineCode"]
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

// getCPUInfo 获取CPU信息
func getCPUInfo() ([]cpuInfo, error) {
	var ret []cpuInfo
	var dst []win32_Processor
	q := wmi.CreateQuery(&dst, "")
	if err := wmiQuery(q, &dst); err != nil {
		return ret, err
	}

	var procID string
	for i, l := range dst {
		procID = ""
		if l.ProcessorID != nil {
			procID = *l.ProcessorID
		}

		cpu := cpuInfo{
			CPU:        int32(i),
			VendorID:   l.Manufacturer,
			PhysicalID: procID,
		}
		ret = append(ret, cpu)
	}

	return ret, nil
}

// wmiQuery WMIQueryWithContext - wraps wmi.Query with a timed-out context to avoid hanging
func wmiQuery(query string, dst interface{}, connectServerArgs ...interface{}) error {
	ctx := context.Background()
	if _, ok := ctx.Deadline(); !ok {
		ctxTimeout, cancel := context.WithTimeout(ctx, 3000000000) //超时时间3s
		defer cancel()
		ctx = ctxTimeout
	}

	errChan := make(chan error, 1)
	go func() {
		errChan <- wmi.Query(query, dst, connectServerArgs...)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errChan:
		return err
	}
}

// getMachineGuid 获取MachineGuid
func getMachineGuid() (string, error) {
	// there has been reports of issues on 32bit using golang.org/x/sys/windows/registry, see https://github.com/shirou/gopsutil/pull/312#issuecomment-277422612
	// for rationale of using windows.RegOpenKeyEx/RegQueryValueEx instead of registry.OpenKey/GetStringValue
	var h windows.Handle
	err := windows.RegOpenKeyEx(windows.HKEY_LOCAL_MACHINE, windows.StringToUTF16Ptr(`SOFTWARE\Microsoft\Cryptography`), 0, windows.KEY_READ|windows.KEY_WOW64_64KEY, &h)
	if err != nil {
		return "", err
	}
	defer windows.RegCloseKey(h)

	const windowsRegBufLen = 74 // len(`{`) + len(`abcdefgh-1234-456789012-123345456671` * 2) + len(`}`) // 2 == bytes/UTF16
	const uuidLen = 36

	var regBuf [windowsRegBufLen]uint16
	bufLen := uint32(windowsRegBufLen)
	var valType uint32
	err = windows.RegQueryValueEx(h, windows.StringToUTF16Ptr(`MachineGuid`), nil, &valType, (*byte)(unsafe.Pointer(&regBuf[0])), &bufLen)
	if err != nil {
		return "", err
	}

	hostID := windows.UTF16ToString(regBuf[:])
	hostIDLen := len(hostID)
	if hostIDLen != uuidLen {
		return "", fmt.Errorf("HostID incorrect: %q\n", hostID)
	}

	return hostID, nil
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
