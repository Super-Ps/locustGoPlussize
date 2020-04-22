package monitor

import (
	"encoding/json"
	"github.com/shirou/gopsutil/host"
	"runtime"
	"time"
)


type OsInfo struct {
	HostID	 			string			`json:"host_id"`
	ProcessPid 			int				`json:"process_pid"`
	ProcessPpid 		int				`json:"process_ppid"`
	ProcessBit	 		int				`json:"process_bit"`
	ProcessName 		string			`json:"process_name"`
	ProcessExe	 		string			`json:"process_exe"`
	ProcessCmd	 		string			`json:"process_cmd"`
	ProcessCwd	 		string			`json:"process_cwd"`
	ProcessUserName	 	string			`json:"process_user_name"`
	OsType				string			`json:"os_type"`
	HostName 			string			`json:"host_name"`
	Bit 				int				`json:"bit"`
	Platform            string			`json:"platform"`
	PlatformFamily      string			`json:"platform_family"`
	CpuVendorId   		string			`json:"cpu_vendor_id"`
	CpuPhysicalId 		string			`json:"cpu_physical_id"`
	CpuModelName  		string			`json:"cpu_model_name"`
	CpuMhz        		int				`json:"cpu_mhz"`
	CpuCount 			int				`json:"cpu_count"`
	CpuThreads 			int				`json:"cpu_threads"`
	MemSwapTotal		int64			`json:"mem_swap_total"`
	MemVirtualTotal		int64			`json:"men_virtual_total"`
	DiskSizeTotal 		int64			`json:"disk_size_total"`
	IPAddress 			string			`json:"ip_address"`
	PublicAddress 		string			`json:"public_address"`
	MaskAddress 		string			`json:"mask_address"`
	MacAddress 			string			`json:"mac_address"`
	NetSpeed 			int				`json:"net_speed"`
	NetCardName 		string			`json:"net_card_name"`
	CreateTimeStamp 	int64			`json:"create_time_stamp"`
	CreateTime		 	string			`json:"create_time"`
	StartTimeStamp		int64			`json:"start_time_stamp"`
	StartTime			string			`json:"start_time"`
}

type SystemDiff struct {
	CpuTime 				int64			`json:"cpu_time"`
	CpuPercent				float64			`json:"cpu_percent"`
	DiskUsed				int64			`json:"disk_used"`
	RunTime 				int64			`json:"run_time"`
	UnixNanoStamp 			int64			`json:"unix_nano_stamp"`
	LastTime 				string			`json:"last_time"`
	ProcessesTotalInfo
	Memory
	DiskIoDiff
	DiskIoTime
	NetworkIoDiff
}

var (
	systemInfo *host.InfoStat
	systemError error
	osInfo *OsInfo
)


func GetOsInfo() *OsInfo {
	return osInfo
}

func GetSystemStartTime() (string, error) {
	if systemError != nil {
		return "", systemError
	}
	return UnixTimeToDateTime(int64(systemInfo.BootTime)*int64(time.Second)), nil
}

func GetSystemStartTimeStamp() (int64, error) {
	if systemError != nil {
		return 0, systemError
	}
	return int64(systemInfo.BootTime)*int64(time.Second), nil
}

func GetHostId() (string, error) {
	if systemError != nil {
		return "", systemError
	}
	return systemInfo.HostID, nil
}

func GetHostName() (string, error) {
	if systemError != nil {
		return "", systemError
	}
	return systemInfo.Hostname, nil
}

func GetPlatform() (string, error) {
	if systemError != nil {
		return "", systemError
	}
	return systemInfo.Platform, nil
}

func GetPlatformFamily() (string, error) {
	if systemError != nil {
		return "", systemError
	}
	return systemInfo.PlatformFamily, nil
}

func GetOsType() (string, error) {
	return runtime.GOOS, nil
}

func GetOsBit() (int, error) {
	return (32 << (^uint(0) >> 63)), nil
}

func GetRunTime() (int64, error) {
	if startTimeStamp, err := GetSystemStartTimeStamp(); err != nil {
		return 0, err
	} else {
		return int64((GetUnixNanoTime()-startTimeStamp)/int64(time.Second)), nil
	}
}

func GetUnixNanoTime() int64 {
	return time.Now().UnixNano()
}

func UnixTimeToDateTime(t int64) string {
	return time.Unix(0, t).Format("2006-01-02 15:04:05")
}

func getOsInfo() *OsInfo {
	var info OsInfo
	info.ProcessPid = curPid
	info.ProcessPpid = curPpid
	info.ProcessBit = curBit
	info.ProcessName = curName
	info.ProcessExe = curExe
	info.ProcessCmd = curCmd
	info.ProcessCwd = curCwd
	info.ProcessUserName = curUserName
	info.HostID, _ = GetHostId()
	info.OsType, _ = GetOsType()
	info.HostName, _ = GetHostName()
	info.Bit, _ = GetOsBit()
	info.Platform, _ = GetPlatform()
	info.PlatformFamily, _ = GetPlatformFamily()
	info.CpuCount, _ = GetCpuCount()
	info.CpuThreads, _ = GetCpuThreads()
	info.CpuVendorId, _ = GetCpuVendorId()
	info.CpuPhysicalId, _ = GetCpuPhysicalId()
	info.CpuModelName, _ = GetCpuModelName()
	info.CpuMhz, _ = GetCpuMhz()
	info.MemSwapTotal, _ = GetTotalSwapMemory()
	info.MemVirtualTotal, _ = GetTotalVirtualMemory()
	info.StartTimeStamp, _ = GetSystemStartTimeStamp()
	info.StartTime, _ = GetSystemStartTime()
	if diskSize, err := GetTotalDiskSize(); err == nil {
		info.DiskSizeTotal = diskSize.Total
	}
	if netList, err := GetNetCardList(); err == nil {
		info.IPAddress = netList[0].IPAddress
		info.MaskAddress = netList[0].MaskAddress
		info.MacAddress = netList[0].MacAddress
		info.NetSpeed = netList[0].NetSpeed
		info.NetCardName = netList[0].NetCardName
	}
	if proInfo, err := GetProcessTimes(curPid); err == nil {
		info.CreateTimeStamp = proInfo.CreateTimeStamp
		info.CreateTime = UnixTimeToDateTime(info.CreateTimeStamp)
	}
	info.PublicAddress, _ = GetPublicIpAddress()
	return &info
}

func (h OsInfo) String() string {
	s, _ := json.Marshal(h)
	return string(s)
}

func (h SystemDiff) String() string {
	s, _ := json.Marshal(h)
	return string(s)
}