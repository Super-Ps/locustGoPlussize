// +build windows

package monitor

import (
	"errors"
	"fmt"
	"github.com/StackExchange/wmi"
	"github.com/gonutz/w32"
	"github.com/shirou/gopsutil/process"
	"golang.org/x/sys/windows"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

type Process struct {
	handle 			windows.Handle
	pid 			int
}

type snapshot struct {
	handle 			windows.Handle
}

type cmdlineMap map[int]string

const (
	SE_DEBUG_PRIVILEGE = 0x14
)

var (
	user32 syscall.Handle
	kernel32 syscall.Handle
	psApi syscall.Handle
	ntDll syscall.Handle
	getGuiResources	uintptr
	getProcessCounters uintptr
	getProcessHandleCount uintptr
	getProcessImageFileNameW uintptr
	getProcessMemoryInfo uintptr
	rtlAdjustPrivilege uintptr
	dosPath map[string]string
)


func newProcess(pid int) (*Process, error) {
	if handle, err := openProcess(pid); err != nil {
		return nil, err
	} else {
		return &Process{handle, pid}, nil
	}
}

func (p *Process) releaseHandle() error {
	return closeHandle(p.handle)
}

func (p *Process) getEntry() (*ProcessEntry, error) {
	if entry, err := getSnapshotEntryByPid(p.pid); err != nil {
		return nil, err
	} else {
		return getEntryInfo(entry, p), nil
	}
}

func (p *Process) getHandles() (int64, error) {
	var handles uint32
	if r1, _, err := syscall.Syscall(uintptr(getProcessHandleCount), 2, uintptr(p.handle), uintptr(unsafe.Pointer(&handles)), 0); r1 == 0 {
		if syscall.GetLastError() != nil {
			return 0, syscall.GetLastError()
		}
		return 0, err
	} else {
		return int64(handles), nil
	}
}

func (p *Process) getGdiObject() (int64, error) {
	if r1, _, err := syscall.Syscall(uintptr(getGuiResources), 2, uintptr(p.handle), uintptr(uint32(0)), 0); r1 == 0 {
		if syscall.GetLastError() != nil {
			return 0, syscall.GetLastError()
		}
		return 0, err
	} else {
		return int64(r1), nil
	}
}

func (p *Process) getUserObject() (int64, error) {
	if r1, _, err := syscall.Syscall(uintptr(getGuiResources), 2, uintptr(p.handle), uintptr(uint32(1)), 0); r1 == 0 {
		if syscall.GetLastError() != nil {
			return 0, syscall.GetLastError()
		}
		return 0, err
	} else {
		return int64(r1), nil
	}
}

func (p *Process) getIoCounters() (*ProcessIoInfo, error) {
	var io_COUNTERS windows.IO_COUNTERS
	if r1, _, err := syscall.Syscall(uintptr(getProcessCounters), 2, uintptr(p.handle), uintptr(unsafe.Pointer(&io_COUNTERS)), 0); r1 == 0 {
		if syscall.GetLastError() != nil {
			return nil, syscall.GetLastError()
		}
		return nil, err
	} else {
		return &ProcessIoInfo{
			ProcessIoData: ProcessIoData{
				IoReadBytes: int64(io_COUNTERS.ReadTransferCount),
				IoWriteBytes: int64(io_COUNTERS.WriteTransferCount),
				IoReadCount: int64(io_COUNTERS.ReadOperationCount),
				IoWriteCount: int64(io_COUNTERS.WriteOperationCount),
			},
			IoUnixNanoStamp: GetUnixNanoTime(),
		}, nil
	}
}

func (p *Process) getTimes() (*ProcessTime, error) {
	var proTime windows.Rusage
	if err := windows.GetProcessTimes(p.handle, &proTime.CreationTime, &proTime.ExitTime, &proTime.KernelTime, &proTime.UserTime); err != nil {
		return nil, err
	} else {
		unixNanoStamp := GetUnixNanoTime()
		createTimeStamp := proTime.CreationTime.Nanoseconds()
		totalTime := int64(proTime.UserTime.HighDateTime << 32 | proTime.UserTime.LowDateTime)+int64(proTime.KernelTime.HighDateTime << 32 | proTime.KernelTime.LowDateTime)
		cpuTime := formatCpuTime(totalTime)
		runTime := (unixNanoStamp - createTimeStamp)/int64(time.Second)
		createTime := UnixTimeToDateTime(createTimeStamp)
		lastTime := UnixTimeToDateTime(unixNanoStamp)

		return &ProcessTime{
			CreateTimeStamp: createTimeStamp,
			TotalTime: totalTime,
			UnixNanoStamp: unixNanoStamp,
			CpuTime: cpuTime,
			RunTime: runTime,
			CreateTime: createTime,
			LastTime: lastTime,
		}, nil
	}
}

func (p *Process) getExePath() (string, error) {
	if p.pid == curPid && curExe != "" {
		return curExe, nil
	}

	buf := make([]uint16, windows.MAX_LONG_PATH)
	size := uint32(windows.MAX_LONG_PATH)

	r1, _, err := syscall.Syscall(uintptr(getProcessImageFileNameW), 3, uintptr(p.handle), uintptr(unsafe.Pointer(&buf[0])), uintptr(size))
	if r1 == 0 {
		return "", err
	}
	path := windows.UTF16ToString(buf[:])
	if path == "" {
		return "", nil
	}

	pathSplit := strings.Split(path, `\`)
	if len(pathSplit) < 3 {
		return path, nil
	}

	rawDrive := strings.Join(pathSplit[:3], `\`)
	if _, ok := dosPath[rawDrive]; ok {
		return filepath.Join(dosPath[rawDrive], path[len(rawDrive):]), nil
	}
	return path, nil
}

func (p *Process) getUserName() (string, error) {
	if p.pid == curPid && curUserName != "" {
		return curUserName, nil
	}

	var token windows.Token
	err := windows.OpenProcessToken(p.handle, windows.TOKEN_QUERY, &token)
	if err != nil {
		return "", err
	}
	defer token.Close()

	tokenUser, err := token.GetTokenUser()
	if err != nil {
		return "", err
	}

	user, _, _, err := tokenUser.User.Sid.LookupAccount("")
	return user, err
}

func (p *Process) getThreads() (int64, error) {
	return getProcessThreads(p.pid)
}

func (p *Process) getParentId() (int, error) {
	return getProcessParentId(p.pid)
}

func (p *Process) getName() (string, error) {
	return getProcessName(p.pid)
}

func (p *Process) getMemory() (*ProcessMemory, error) {
	var memory process.PROCESS_MEMORY_COUNTERS
	if r1, _, err := syscall.Syscall(uintptr(getProcessMemoryInfo), 3, uintptr(p.handle), uintptr(unsafe.Pointer(&memory)), uintptr(unsafe.Sizeof(memory))); r1 == 0 {
		if syscall.GetLastError() != nil {
			return nil, syscall.GetLastError()
		}
		return nil, err
	} else {
		return &ProcessMemory{Rss:int64(memory.WorkingSetSize), Vms: int64(memory.PagefileUsage)}, nil
	}
}

func (p *Process) getCwd() (string, error) {
	if p.pid == curPid && curCwd != "" {
		return curCwd, nil
	}

	return GetProcessWorkDirByHandle(p.handle)
}

func (p *Process) getCmdline() (string, error) {
	if p.pid == curPid && curCmd != "" {
		return curCmd, nil
	}

	return GetProcessCommandLineByHandle(p.handle)
}

func (p *Process) getBit() (int, error) {
	if p.pid == curPid && curBit != 0 {
		return curBit, nil
	}

	var wow64Process bool
	if err := windows.IsWow64Process(p.handle, &wow64Process); err != nil {
		return 0, err
	} else {
		if wow64Process {
			return 32, nil
		} else {
			return 64, nil
		}
	}
}

func getProcessCpuPercent(p1 *ProcessTime, p2 *ProcessTime) float64 {
	if (p2.UnixNanoStamp - p1.UnixNanoStamp) <= 0 {
		return 0
	}
	per, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(float64(p2.TotalTime-p1.TotalTime)/float64(p2.UnixNanoStamp-p1.UnixNanoStamp)/float64(osInfo.CpuThreads))*10000), 64)
	return math.Min(math.Max(per,0), 100)
}

func getProcessThreads(pid int) (int64, error) {
	if proList, err := openProcessList(); err != nil {
		return 0, err
	} else {
		defer proList.closeProcessList()
		if pro, err := proList.getProcessByPid(pid); err != nil {
			return 0, err
		} else {
			return int64(pro.Threads), nil
		}
	}
}

func getProcessParentId(pid int) (int, error) {
	if pid == curPid && curPpid != 0 {
		return curPpid, nil
	}

	if proList, err := openProcessList(); err != nil {
		return 0, err
	} else {
		defer proList.closeProcessList()
		if pro, err := proList.getProcessByPid(pid); err != nil {
			return 0, err
		} else {
			return int(pro.ParentProcessID), nil
		}
	}
}

func getProcessName(pid int) (string, error) {
	if pid == curPid && curName != "" {
		return curName, nil
	}

	if proList, err := openProcessList(); err != nil {
		return "", err
	} else {
		defer proList.closeProcessList()
		if pro, err := proList.getProcessByPid(pid); err != nil {
			return "", err
		} else {
			return formatProcessName(pro.ExeFile[:]), nil
		}
	}
}

func getProcessGdiObject(pid int) (int64, error) {
	if pro, err := newProcess(pid); err != nil {
		return 0, err
	} else {
		defer pro.releaseHandle()
		return pro.getGdiObject()
	}
}

func getProcessUserObject(pid int) (int64, error) {
	if pro, err := newProcess(pid); err != nil {
		return 0, err
	} else {
		defer pro.releaseHandle()
		return pro.getUserObject()
	}
}

func getProcessHandles(pid int) (int64, error) {
	if pro, err := newProcess(pid); err != nil {
		return 0, err
	} else {
		defer pro.releaseHandle()
		return pro.getHandles()
	}
}

func getProcessIoCounters(pid int) (*ProcessIoInfo, error) {
	if pro, err := newProcess(pid); err != nil {
		return nil, err
	} else {
		defer pro.releaseHandle()
		return pro.getIoCounters()
	}
}

func getProcessTimes(pid int) (*ProcessTime, error) {
	if pro, err := newProcess(pid); err != nil {
		return nil, err
	} else {
		defer pro.releaseHandle()
		return pro.getTimes()
	}
}

func getProcessExePath(pid int) (string, error) {
	if pid == curPid && curExe != "" {
		return curExe, nil
	}

	if pro, err := newProcess(pid); err != nil {
		return "", err
	} else {
		defer pro.releaseHandle()
		return pro.getExePath()
	}
}

func getProcessUserName(pid int) (string, error) {
	if pid == curPid && curUserName != "" {
		return curUserName, nil
	}

	if pro, err := newProcess(pid); err != nil {
		return "", err
	} else {
		defer pro.releaseHandle()
		return pro.getUserName()
	}
}

func getProcessCmdline(pid int) (string, error) {
	if pid == curPid && curCmd != "" {
		return curCmd, nil
	}

	return GetProcessCommandLineByPid(pid)
}

func getProcessMemory(pid int) (*ProcessMemory, error) {
	if pro, err := newProcess(pid); err != nil {
		return nil, err
	} else {
		defer pro.releaseHandle()
		return pro.getMemory()
	}
}

func getProcessCwd(pid int) (string, error) {
	if pid == curPid && curCwd != "" {
		return curCwd, nil
	}

	return GetProcessWorkDirByPid(pid)
}

func getProcessBit(pid int) (int, error) {
	if pid == curPid && curBit != 0 {
		return curBit, nil
	}

	if pro, err := newProcess(pid); err != nil {
		return 0, err
	} else {
		defer pro.releaseHandle()
		return pro.getBit()
	}
}

func getProcessEntry(pid int) (*ProcessEntry, error) {
	if pro, err := newProcess(pid); err != nil {
		return nil, err
	} else {
		defer pro.releaseHandle()
		return pro.getEntry()
	}
}

func getProcessesTotalInfo() (*ProcessesTotalInfo, error) {
	var info ProcessesTotalInfo
	if proList, err := getProcessList(); err != nil {
		return nil, err
	} else {
		info.Processes = 0
		for _, v := range proList {
			info.Processes += 1
			info.Handles += v.Handles
			info.Threads += v.Threads
			info.GdiObject += v.GdiObject
			info.UserObject += v.UserObject
		}
		return &info, nil
	}
}

func getTotalHandles() (int64, error) {
	var handles int64
	if proList, err := getProcessList(); err != nil {
		return 0, err
	} else {
		for _, v := range proList {
			handles += v.Handles
		}
		return handles, nil
	}
}

func getTotalThreads() (int64, error) {
	var threads int64
	if proList, err := getProcessList(); err != nil {
		return 0, err
	} else {
		for _, v := range proList {
			threads += v.Threads
		}
		return threads, nil
	}
}

func getTotalGdiObject() (int64, error) {
	var gdiObject int64
	if proList, err := getProcessList(); err != nil {
		return 0, err
	} else {
		for _, v := range proList {
			gdiObject += v.GdiObject
		}
		return gdiObject, nil
	}
}

func getTotalUserObject() (int64, error) {
	var userObject int64
	if proList, err := getProcessList(); err != nil {
		return 0, err
	} else {
		for _, v := range proList {
			userObject += v.UserObject
		}
		return userObject, nil
	}
}

func getTotalProcesses() (int64, error) {
	if proList, err := process.Pids(); err != nil {
		return 0, err
	} else {
		return int64(len(proList)), nil
	}
}

func openProcess(pid int) (windows.Handle, error) {
	return windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION | windows.PROCESS_VM_READ, true, uint32(pid))
}

func getProcessCmdlineFromWmi(pid int) (string, error) {
	type wmiData struct{
		CommandLine 	string
		ProcessId 		int
	}

	var wmiList []wmiData

	if err := wmi.Query(fmt.Sprintf("Select CommandLine,ProcessId from Win32_Process where ProcessId='%d'", pid), &wmiList); err != nil {
		return "", err
	}
	if len(wmiList) == 0 {
		return "", errors.New(fmt.Sprintf("Can not found pid:%d", pid))
	}
	return wmiList[0].CommandLine, nil
}

func openSnapshot() (windows.Handle, error) {
	return windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
}

func closeHandle(handle windows.Handle) error {
	return windows.CloseHandle(handle)
}

func openProcessList() (*snapshot, error) {
	if handle, err := openSnapshot(); err != nil {
		return nil, err
	} else {
		return &snapshot{handle}, nil
	}
}

func (s *snapshot) closeProcessList() error {
	return  closeHandle(s.handle)
}

func (s *snapshot) getProcessByPid(pid int) (*windows.ProcessEntry32, error) {
	proEntry := windows.ProcessEntry32{}
	proEntry.Size = uint32(unsafe.Sizeof(proEntry))

	if err := windows.Process32First(s.handle, &proEntry); err != nil {
		return nil, err
	}

	for {
		if err := windows.Process32Next(s.handle, &proEntry); err != nil {
			break
		} else {
			if int(proEntry.ProcessID) == pid {
				return &proEntry, nil
			}
		}
	}
	return nil, errors.New(fmt.Sprintf("Can not found pid:%d", pid))
}

func getProcessCmdlineMapFromWmi() (cmdlineMap, error) {
	type wmiData struct{
		CommandLine 	string
		ProcessId 		int
	}
	var wmiList []wmiData

	cmdline := make(cmdlineMap)
	err := wmi.Query("Select CommandLine,ProcessId from Win32_Process", &wmiList)
	if err != nil {
		return nil, err
	}

	for _, v := range wmiList {
		cmdline[v.ProcessId] = v.CommandLine
	}
	return cmdline, nil
}

func getProcessList() ([]*ProcessEntry, error) {
	handle, err := openSnapshot()
	if err != nil {
		return nil, err
	}
	defer closeHandle(handle)

	var proEntryList []*ProcessEntry
	entry := windows.ProcessEntry32{}
	entry.Size = uint32(unsafe.Sizeof(entry))
	if err := windows.Process32First(handle, &entry); err != nil {
		return nil, err
	}

	for {
		if err := windows.Process32Next(handle, &entry); err != nil {
			break
		} else {
			pro, err := newProcess(int(entry.ProcessID))
			if err != nil {
				continue
			}
			proEntryList = append(proEntryList, getEntryInfo(&entry, pro))
			_ = pro.releaseHandle()
		}
	}
	return proEntryList, nil
}

func getProcessMap() (map[int]*ProcessEntry, error) {
	handle, err := openSnapshot()
	if err != nil {
		return nil, err
	}
	defer closeHandle(handle)

	proEntryMap := make(map[int]*ProcessEntry)
	entry := windows.ProcessEntry32{}
	entry.Size = uint32(unsafe.Sizeof(entry))
	if err := windows.Process32First(handle, &entry); err != nil {
		return nil, err
	}

	for {
		if err := windows.Process32Next(handle, &entry); err != nil {
			break
		} else {
			pro, err := newProcess(int(entry.ProcessID))
			if err != nil {
				continue
			}
			proEntryMap[int(entry.ProcessID)] = getEntryInfo(&entry, pro)
			_ = pro.releaseHandle()
		}
	}

	return proEntryMap, nil
}

func getEntryInfo(entry *windows.ProcessEntry32, pro *Process) *ProcessEntry {
	var proEntry ProcessEntry

	proEntry.Pid = int(entry.ProcessID)
	proEntry.Ppid = int(entry.ParentProcessID)
	proEntry.Threads = int64(entry.Threads)
	proEntry.Name = formatProcessName(entry.ExeFile[:])

	if proMem, err := pro.getMemory(); err == nil {
		proEntry.ProcessMemory = *proMem
	}

	if proIo, err := pro.getIoCounters(); err == nil {
		proEntry.ProcessIoInfo = *proIo
	}

	if proTimes, err := pro.getTimes(); err == nil {
		proEntry.ProcessTime = *proTimes
	}

	if bit, err := pro.getBit(); err == nil {
		proEntry.Bit = bit
	}

	if handles, err := pro.getHandles(); err == nil {
		proEntry.Handles = handles
	}

	if cwd, err := pro.getCwd(); err == nil {
		proEntry.Cwd = cwd
	}

	if cmd, err := pro.getCmdline(); err == nil {
		proEntry.Cmd = cmd
	}

	if exe, err := pro.getExePath(); err == nil {
		proEntry.Exe = exe
	}

	if user, err := pro.getUserName(); err == nil {
		proEntry.UserName = user
	}

	if gdiObj, err := pro.getGdiObject(); err == nil {
		proEntry.GdiObject = gdiObj
	}

	if userObj, err := pro.getUserObject(); err == nil {
		proEntry.UserObject = userObj
	}

	if fileInfo, err := getProcessFileInfo(proEntry.Exe); err == nil {
		proEntry.Size = fileInfo.Size
		proEntry.FileModifyTime = fileInfo.FileModifyTime
		proEntry.FileVersion = fileInfo.FileVersion
		proEntry.ProductVersion = fileInfo.ProductVersion
	} else {
		if fileInfo != nil {
			proEntry.Size = fileInfo.Size
			proEntry.FileModifyTime = fileInfo.FileModifyTime
		}
	}

	return &proEntry
}

func formatProcessName(name []uint16) string {
	return windows.UTF16ToString(name)
}

func getSnapshotEntryByPid(pid int) (*windows.ProcessEntry32, error) {
	handle, err := openSnapshot()
	if err != nil {
		return nil, err
	}
	defer closeHandle(handle)

	entry := windows.ProcessEntry32{}
	entry.Size = uint32(unsafe.Sizeof(entry))
	if err := windows.Process32First(handle, &entry); err != nil {
		return nil, err
	}

	for {
		if err := windows.Process32Next(handle, &entry); err != nil {
			break
		} else {
			if int(entry.ProcessID) == pid {
				return &entry, nil
			}
		}
	}
	return nil,errors.New(fmt.Sprintf("Can not found pid:%d", pid))
}

func getDOSPathMap() map[string]string {
	pathMap := make(map[string]string)
	for d := 'A'; d <= 'Z'; d++ {
		deviceName := string(d) + ":"
		targetPath := make([]uint16, windows.MAX_LONG_PATH)
		re, _ :=  windows.QueryDosDevice(windows.StringToUTF16Ptr(deviceName), &targetPath[0], windows.MAX_LONG_PATH)
		if re == 0 {
			continue
		}
		pathMap[windows.UTF16ToString(targetPath[:])] = deviceName
	}
	return pathMap
}

func setDebugPrivilege(isDebug bool) bool {
	var turnOn uint32
	if isDebug {
		turnOn = 1
	} else {
		turnOn = 0
	}
	currentThread := uint32(0)
	var enabled1 uint32
	var enabled2 uint32
	_, _, _ = syscall.Syscall6(uintptr(rtlAdjustPrivilege), 4, uintptr(SE_DEBUG_PRIVILEGE), uintptr(turnOn), uintptr(currentThread), uintptr(unsafe.Pointer(&enabled1)), 0, 0)
	_, _, _ = syscall.Syscall6(uintptr(rtlAdjustPrivilege), 4, uintptr(SE_DEBUG_PRIVILEGE), uintptr(turnOn), uintptr(currentThread), uintptr(unsafe.Pointer(&enabled2)), 0, 0)
	if enabled2 != turnOn {
		return false
	}
	return  true
}

func getProcessFileInfo(name string) (*ProcessFileInfo, error) {
	var fileInfo ProcessFileInfo
	statInfo, err := os.Stat(name)
	if err != nil {
		return &fileInfo, err
	}
	fileInfo.Size = statInfo.Size()
	fileInfo.FileModifyTime = statInfo.ModTime().Format("2006-01-02 15:04:05")

	versionSize := w32.GetFileVersionInfoSize(name)
	if versionSize <= 0 {
		return &fileInfo, errors.New("GetFileVersionInfoSize failed")
	}

	versionInfo := make([]byte, versionSize)
	if ok := w32.GetFileVersionInfo(name, versionInfo); !ok {
		return &fileInfo, errors.New("GetFileVersionInfo failed")
	}

	fixed, ok := w32.VerQueryValueRoot(versionInfo)
	if !ok {
		return &fileInfo, errors.New("VerQueryValueRoot failed")
	}
	fileVersionData := uint64(fixed.FileVersionMS)<<32 | uint64(fixed.FileVersionLS)
	ProductVersionData := uint64(fixed.ProductVersionMS)<<32 | uint64(fixed.ProductVersionLS)

	fileInfo.FileVersion = fmt.Sprintf("%d.%d.%d.%d",
		fileVersionData&0xFFFF000000000000>>48,
		fileVersionData&0x0000FFFF00000000>>32,
		fileVersionData&0x00000000FFFF0000>>16,
		fileVersionData&0x000000000000FFFF>>0)

	fileInfo.ProductVersion = fmt.Sprintf("%d.%d.%d.%d",
		ProductVersionData&0xFFFF000000000000>>48,
		ProductVersionData&0x0000FFFF00000000>>32,
		ProductVersionData&0x00000000FFFF0000>>16,
		ProductVersionData&0x000000000000FFFF>>0)

	return &fileInfo, nil
}

func formatCpuTime(times int64) int64 {
	return int64(times/10000000)
}