package monitor

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/shirou/gopsutil/process"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

type ProcessMemory struct {
	Rss 				int64				`json:"rss"`
	Vms 				int64				`json:"vms"`
}

type ProcessTime struct {
	CreateTimeStamp 	int64				`json:"create_time_stamp"`
	UnixNanoStamp 		int64				`json:"unix_nano_stamp"`
	TotalTime			int64				`json:"total_time"`
	CpuTime 			int64				`json:"cpu_time"`
	RunTime 			int64				`json:"run_time"`
	CreateTime			string				`json:"create_time"`
	LastTime 			string				`json:"last_time"`
}

type ProcessIoDiff struct {
	ReadBytesDiff 		int64				`json:"io_read_bytes_diff"`
	WriteBytesDiff 		int64				`json:"io_write_bytes_diff"`
	ReadCountDiff		int64				`json:"io_read_count_diff"`
	WriteCountDiff 		int64				`json:"io_write_count_diff"`
}

type ProcessBaseInfo struct {
	Pid 				int				`json:"pid"`
	Ppid 				int				`json:"ppid"`
	Bit 				int				`json:"bit"`
	Name 				string			`json:"name"`
	UserName 			string			`json:"user_name"`
	Exe 				string			`json:"exe"`
	Cmd 				string			`json:"cmd"`
	Cwd 				string			`json:"cwd"`
}

type ProcessFileInfo struct {
	Size				int64			`json:"size"`
	FileVersion			string			`json:"file_version"`
	ProductVersion		string			`json:"product_version"`
	FileModifyTime		string			`json:"file_modify_time"`
}

type ProcessesTotalInfo struct {
	Processes 			int64			`json:"processes"`
	ProcessesCountInfo
}

type ProcessesCountInfo struct {
	Handles				int64			`json:"handles"`
	Threads				int64			`json:"threads"`
	GdiObject			int64			`json:"gdi_object"`
	UserObject			int64			`json:"user_object"`
}

type ProcessIoData struct {
	IoReadBytes 		int64			`json:"io_read_bytes"`
	IoWriteBytes 		int64			`json:"io_write_bytes"`
	IoReadCount 		int64			`json:"io_read_count"`
	IoWriteCount 		int64			`json:"io_write_count"`
}

type ProcessIoInfo struct {
	ProcessIoData
	IoUnixNanoStamp 	int64			`json:"io_unix_nano_stamp"`
}

type ProcessEntry struct {
	ProcessBaseInfo
	ProcessesCountInfo
	ProcessMemory
	ProcessIoInfo
	ProcessTime
	ProcessFileInfo
}

type ProcessDiff struct {
	ProcessBaseInfo
	ProcessesCountInfo
	ProcessMemory
	ProcessIoData
	ProcessIoDiff
	ProcessTime
	ProcessFileInfo
	CpuPercent 			float64			`json:"cpu_percent"`
}

var curPid, curPpid, curBit int
var curExe, curName, curUserName, curCmd, curCwd string


func GetProcessPidList() ([]int, error) {
	if proList, err := process.Pids(); err != nil {
		return nil, err
	} else {
		var pidList []int
		for _, v := range proList {
			pidList = append(pidList, int(v))
		}
		return pidList, nil
	}
}

func GetProcessList() ([]*ProcessEntry, error) {
	return getProcessList()
}

func GetProcessMap() (map[int]*ProcessEntry, error) {
	return getProcessMap()
}

func GetProcessDiffInfo(pro1 *ProcessEntry, pro2 *ProcessEntry) *ProcessDiff {
	return &ProcessDiff{
		ProcessBaseInfo: pro2.ProcessBaseInfo,
		ProcessesCountInfo: pro2.ProcessesCountInfo,
		ProcessMemory: pro2.ProcessMemory,
		ProcessIoData: pro2.ProcessIoData,
		ProcessIoDiff: *GetProcessIoDiff(&pro1.ProcessIoInfo, &pro2.ProcessIoInfo),
		ProcessTime: pro2.ProcessTime,
		CpuPercent: GetProcessCpuPercent(&pro1.ProcessTime, &pro2.ProcessTime),
	}
}

func OpenProcess(pid int) (*Process, error) {
	return newProcess(pid)
}

func (p *Process) CloseHandle() error {
	return p.releaseHandle()
}

func (p *Process) GetEntry() (*ProcessEntry, error) {
	return p.getEntry()
}

func (p *Process) GetHandles() (int64, error) {
	return p.getHandles()
}

func (p *Process) GetGdiObject() (int64, error) {
	return p.getGdiObject()
}

func (p *Process) GetUserObject() (int64, error) {
	return p.getUserObject()
}

func (p *Process) GetIoCounters() (*ProcessIoInfo, error) {
	return p.getIoCounters()
}

func (p *Process) GetTimes() (*ProcessTime, error) {
	return p.getTimes()
}

func (p *Process) GetExePath() (string, error) {
	return p.getExePath()
}

func (p *Process) GetUserName() (string, error) {
	return p.getUserName()
}

func (p *Process) GetThreads() (int64, error) {
	return p.getThreads()
}

func (p *Process) GetParentId() (int, error) {
	return p.getParentId()
}

func (p *Process) GetName() (string, error) {
	return p.getName()
}

func (p *Process) GetMemory() (*ProcessMemory, error) {
	return p.getMemory()
}

func (p *Process) GetCwd() (string, error) {
	return p.getCwd()
}

func (p *Process) GetCmdline() (string, error) {
	return p.getCmdline()
}

func (p *Process) GetBit() (int, error) {
	return p.getBit()
}

func GetProcessEntry(pid int) (*ProcessEntry, error) {
	return getProcessEntry(pid)
}

func GetProcessThreads(pid int) (int64, error) {
	return getProcessThreads(pid)
}

func GetProcessParentId(pid int) (int, error) {
	return getProcessParentId(pid)
}

func GetProcessName(pid int) (string, error) {
	return getProcessName(pid)
}

func GetProcessGdiObject(pid int) (int64, error) {
	return getProcessGdiObject(pid)
}

func GetProcessUserObject(pid int) (int64, error) {
	return getProcessUserObject(pid)
}

func GetProcessHandles(pid int) (int64, error) {
	return getProcessHandles(pid)
}

func GetProcessIoCounters(pid int) (*ProcessIoInfo, error) {
	return getProcessIoCounters(pid)
}

func GetProcessTimes(pid int) (*ProcessTime, error) {
	return getProcessTimes(pid)
}

func GetProcessExePath(pid int) (string, error) {
	return getProcessExePath(pid)
}

func GetProcessUserName(pid int) (string, error) {
	return getProcessUserName(pid)
}

func GetProcessCmdline(pid int) (string, error) {
	return getProcessCmdline(pid)
}

func GetProcessMemory(pid int) (*ProcessMemory, error) {
	return getProcessMemory(pid)
}

func GetProcessCwd(pid int) (string, error) {
	return getProcessCwd(pid)
}

func GetProcessBit(pid int) (int, error) {
	return getProcessBit(pid)
}

func ProcessExistFromExePath(exe string) (bool, error) {
	if proList, err := process.Pids(); err != nil {
		return false, err
	} else {
		for _, v := range proList {
			if proPath, err := GetProcessExePath(int(v)); err != nil {
				continue
			} else {
				if exe == proPath {
					return true, nil
				}
			}
		}
		return false, errors.New(fmt.Sprintf("Can not found process path: %d", exe))
	}
}

func ProcessExistFromName(name string) (bool, error) {
	if proList, err := process.Pids(); err != nil {
		return false, err
	} else {
		for _, v := range proList {
			if proName, err := GetProcessName(int(v)); err != nil {
				continue
			} else {
				if name == proName {
					return true, nil
				}
			}
		}
		return false, errors.New(fmt.Sprintf("Can not found process name: %d", name))
	}
}

func ProcessExistFromPid(pid int) (bool, error) {
	if proList, err := process.Pids(); err != nil {
		return false, err
	} else {
		for _, v := range proList {
			if pid == int(v) {
				return true, nil
			}
		}
		return false, errors.New(fmt.Sprintf("Can not found process pid: %d", pid))
	}
}

func KillProcess(pid int) error {
	if pro, err := process.NewProcess(int32(pid)); err != nil {
		return nil
	} else {
		return pro.Kill()
	}
}

func GetTotalHandles() (int64, error) {
	return getTotalHandles()
}

func GetTotalThreads() (int64, error) {
	return getTotalThreads()
}

func GetTotalGdiObject() (int64, error) {
	return getTotalGdiObject()
}

func GetTotalUserObject() (int64, error) {
	return getTotalUserObject()
}

func GetTotalProcesses() (int64, error) {
	return getTotalProcesses()
}

func GetProcessesTotalInfo() (*ProcessesTotalInfo, error) {
	return getProcessesTotalInfo()
}

func GetProcessCpuPercent(p1 *ProcessTime, p2 *ProcessTime) float64 {
	return getProcessCpuPercent(p1, p2)
}

func GetProcessIoDiff(d1 *ProcessIoInfo, d2 *ProcessIoInfo) *ProcessIoDiff {
	var ioDiff ProcessIoDiff
	times := float64((d2.IoUnixNanoStamp-d1.IoUnixNanoStamp)/int64(time.Second))
	if times <= 0 {
		times = float64(1)
	}
	ioDiff.ReadBytesDiff = FormatNegativeInt64(int64(float64(d2.IoReadBytes-d1.IoReadBytes)/times))
	ioDiff.WriteBytesDiff = FormatNegativeInt64(int64(float64(d2.IoWriteBytes-d1.IoWriteBytes)/times))
	ioDiff.ReadCountDiff = FormatNegativeInt64(int64(float64(d2.IoReadCount-d1.IoReadCount)/times))
	ioDiff.WriteCountDiff = FormatNegativeInt64(int64(float64(d2.IoWriteCount-d1.IoWriteCount)/times))
	return &ioDiff
}

func FormatProcessDiff(proEntry *ProcessEntry, ioDiff *ProcessIoDiff, cpuPer float64) *ProcessDiff {
	var proDiff ProcessDiff
	proDiff.ProcessBaseInfo = proEntry.ProcessBaseInfo
	proDiff.ProcessesCountInfo = proEntry.ProcessesCountInfo
	proDiff.ProcessMemory = proEntry.ProcessMemory
	proDiff.ProcessTime = proEntry.ProcessTime
	proDiff.ProcessIoData = proEntry.ProcessIoData
	proDiff.ProcessIoDiff = *ioDiff
	proDiff.ProcessFileInfo = proEntry.ProcessFileInfo
	proDiff.CpuPercent = cpuPer
	return &proDiff
}

func GetCurrentCmdline() string {
	var cmd string
	for _, v := range os.Args {
		cmd += fmt.Sprintf("%s ", v)
	}
	return TrimString(cmd)
}

func CmdCall(cmdLine string) (string, error) {
	var bin, arg string
	var stdout, stderr bytes.Buffer
	if runtime.GOOS == "windows" {
		bin = "cmd"
		arg = "/C"
	} else {
		bin = "bash"
		arg = "-c"
	}
	cmd := exec.Command(bin, arg, cmdLine)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", err
	} else {
		return TrimString(stdout.String()), nil
	}
}

func TrimString(s string) string {
	s = strings.TrimSpace(s)
	s = strings.Trim(s, "")
	s = strings.Trim(s, " ")
	s = strings.Trim(s, "\n")
	s = strings.Trim(s, "\r\n")
	return s
}

func (h ProcessMemory) String() string {
	s, _ := json.Marshal(h)
	return string(s)
}

func (h ProcessTime) String() string {
	s, _ := json.Marshal(h)
	return string(s)
}

func (h ProcessIoDiff) String() string {
	s, _ := json.Marshal(h)
	return string(s)
}

func (h ProcessBaseInfo) String() string {
	s, _ := json.Marshal(h)
	return string(s)
}

func (h ProcessIoInfo) String() string {
	s, _ := json.Marshal(h)
	return string(s)
}

func (h ProcessesTotalInfo) String() string {
	s, _ := json.Marshal(h)
	return string(s)
}

func (h ProcessEntry) String() string {
	s, _ := json.Marshal(h)
	return string(s)
}

func (h ProcessDiff) String() string {
	s, _ := json.Marshal(h)
	return string(s)
}