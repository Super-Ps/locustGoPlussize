//+build linux freebsd darwin openbsd

package monitor

import (
	"fmt"
	"github.com/shirou/gopsutil/process"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Process struct {
	handle			*process.Process
	pid 			int
}

type handleMap map[int]int64


func newProcess(pid int) (*Process, error) {
	if pro, err := process.NewProcess(int32(pid)); err != nil {
		return nil, err
	} else {
		return &Process{pro, pid}, nil
	}
}

func (p *Process) releaseHandle() error {
	return nil
}

func (p *Process) getEntry() (*ProcessEntry, error) {
	return getProcessEntry(p.pid)
}

func (p *Process) getUserName() (string, error) {
	if p.pid == curPid && curUserName != "" {
		return curUserName, nil
	}

	return p.handle.Username()
}

func (p *Process) getParentId() (int, error) {
	return getProcessParentId(p.pid)
}

func (p *Process) getProcessParentId() (int, error) {
	if p.pid == curPid && curPpid != 0 {
		return curPpid, nil
	}

	if ppid, err := p.handle.Ppid(); err != nil {
		return 0, err
	} else {
		return int(ppid), nil
	}
}

func (p *Process) getIoCounters() (*ProcessIoInfo, error) {
	if io, err := p.handle.IOCounters(); err != nil {
		return nil, err
	} else {
		return &ProcessIoInfo{
			ProcessIoData: ProcessIoData{
				IoReadBytes: int64(io.ReadBytes),
				IoWriteBytes: int64(io.WriteBytes),
				IoReadCount: int64(io.ReadCount),
				IoWriteCount: int64(io.WriteCount),
			},
			IoUnixNanoStamp: GetUnixNanoTime(),
		}, nil
	}
}

func (p *Process) getGdiObject() (int64, error) {
	return getProcessGdiObject(p.pid)
}

func (p *Process) getUserObject() (int64, error) {
	return getProcessUserObject(p.pid)
}

func (p *Process) getExePath() (string, error) {
	if p.pid == curPid && curExe != "" {
		return curExe, nil
	}

	return p.handle.Exe()
}

func (p *Process) getName() (string, error) {
	if p.pid == curPid && curName != "" {
		return curName, nil
	}

	return p.handle.Name()
}

func (p *Process) getCwd() (string, error) {
	if p.pid == curPid && curCwd != "" {
		return curCwd, nil
	}

	return p.handle.Cwd()
}

func (p *Process) getCmdline() (string, error) {
	if p.pid == curPid && curCmd != "" {
		return curCmd, nil
	}

	return p.handle.Cmdline()
}

func (p *Process) getBit() (int, error) {
	if curBit != 0 {
		return curBit, nil
	}

	return (32 << (^uint(0) >> 63)), nil
}

func (p *Process) getThreads() (int64, error) {
	if threads, err := p.handle.NumThreads(); err != nil {
		return 0, err
	} else {
		return int64(threads), nil
	}
}

func (p *Process) getHandles() (int64, error) {
	if handles, err := p.handle.NumFDs(); err != nil {
		return 0, err
	} else {
		return int64(handles), nil
	}
}

func (p *Process) getMemory() (*ProcessMemory, error) {
	if memory, err := p.handle.MemoryInfo(); err != nil {
		return nil, err
	} else {
		return &ProcessMemory{Rss:int64(memory.RSS), Vms:int64(memory.VMS)}, nil
	}
}

func (p *Process) getTimes() (*ProcessTime, error) {
	if proTime, err := p.handle.Times(); err != nil {
		return nil, err
	} else {
		unixNanoStamp := GetUnixNanoTime()
		createTimeStamp, _ := p.handle.CreateTime()
		createTimeStamp *= int64(time.Millisecond)
		totalTime := int64(proTime.Total()*float64(time.Second))
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

func getProcessCpuPercent(p1 *ProcessTime, p2 *ProcessTime) float64 {
	if (p2.UnixNanoStamp - p1.UnixNanoStamp) <= 0 {
		return 0
	}
	per, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(p2.TotalTime-p1.TotalTime)/float64(p2.UnixNanoStamp-p1.UnixNanoStamp)/float64(osInfo.CpuThreads)*100), 64)
	return math.Min(math.Max(per,0), 100)
}

func getProcessThreads(pid int) (int64, error) {
	if pro, err := newProcess(pid); err != nil {
		return 0, err
	} else {
		defer pro.releaseHandle()
		return pro.getThreads()
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

func getProcessParentId(pid int) (int, error) {
	if pid == curPid && curPpid != 0 {
		return curPpid, nil
	}

	if pro, err := newProcess(pid); err != nil {
		return 0, err
	} else {
		defer pro.releaseHandle()
		return pro.getProcessParentId()
	}
}

func getProcessName(pid int) (string, error) {
	if pid == curPid && curName != "" {
		return curName, nil
	}

	if pro, err := newProcess(pid); err != nil {
		return "", err
	} else {
		defer pro.releaseHandle()
		return pro.getName()
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

func getProcessMemory(pid int) (*ProcessMemory, error) {
	if pro, err := newProcess(pid); err != nil {
		return nil, err
	} else {
		defer pro.releaseHandle()
		return pro.getMemory()
	}
}

func getProcessCmdline(pid int) (string, error) {
	if pid == curPid && curCmd != "" {
		return curCmd, nil
	}

	if pro, err := newProcess(pid); err != nil {
		return "", err
	} else {
		defer pro.releaseHandle()
		return pro.getCmdline()
	}
}

func getProcessCwd(pid int) (string, error) {
	if pid == curPid && curCwd != "" {
		return curCwd, nil
	}

	if pro, err := newProcess(pid); err != nil {
		return "", err
	} else {
		defer pro.releaseHandle()
		return pro.getCwd()
	}
}

func getProcessBit(pid int) (int, error) {
	if curBit != 0 {
		return curBit, nil
	}

	if pro, err := newProcess(pid); err != nil {
		return 0, err
	} else {
		defer pro.releaseHandle()
		return pro.getBit()
	}
}

func getProcessGdiObject(pid int) (int64, error) {
	return 0, nil
}

func getProcessUserObject(pid int) (int64, error) {
	return 0, nil
}

func getTotalHandles() (int64, error) {
	out, err := CmdCall("lsof | wc -l")
	if err != nil {
		return 0, err
	}

	count, err := strconv.Atoi(strings.Trim(out,""))
	if err != nil {
		return 0, err
	}
	return int64(count), nil
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
	return 0, nil
}

func getTotalUserObject() (int64, error) {
	return 0, nil
}

func getTotalProcesses() (int64, error) {
	if proList, err := process.Pids(); err != nil {
		return 0, err
	} else {
		return int64(len(proList)), nil
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
			info.Threads += v.Threads
			info.Handles += v.Handles
		}
		return &info, nil
	}
}

func getProcessHandleFromCmd(pid int) (int64, error) {
	out, err := CmdCall(fmt.Sprintf("lsof -p %d | wc -l", pid))
	if err != nil {
		return 0, err
	}

	count, err := strconv.Atoi(strings.Trim(out,""))
	if err != nil {
		return 0, err
	}

	return int64(count), nil
}

func getProcessHandleMapFromCmd() (handleMap, error) {
	out, err := CmdCall(`lsof | awk '{print $2}' | sort -n | uniq -c`)
	if err != nil {
		return nil, err
	}

	handle := make(handleMap)
	handleIndex := 1
	pidIndex := 2
	reString := `(\d+)\s+(\d+)`
	split := "\n"

	re := regexp.MustCompile(reString)
	outSplit := strings.Split(out, split)

	for i, data := range outSplit {
		if i == 0 {
			continue
		}
		dataLine := TrimString(data)
		handleSplit := re.FindAllStringSubmatch(dataLine, -1)
		if len(handleSplit) < 1 {
			continue
		}
		if len(handleSplit[0]) < 3 {
			continue
		}
		handles, _ := strconv.Atoi(strings.Trim(handleSplit[0][handleIndex],""))
		pid, _ := strconv.Atoi(strings.Trim(handleSplit[0][pidIndex],""))
		handle[pid] = int64(handles)
	}

	return handle, nil
}

func getProcessList() ([]*ProcessEntry, error) {
	var proEntryList []*ProcessEntry
	if pidList, err := process.Pids(); err != nil {
		return nil, err
	} else {
		for _, v := range pidList {
			proEntryList = append(proEntryList, getEntryInfo(int(v)))
		}
		return proEntryList, nil
	}
}

func getProcessMap() (map[int]*ProcessEntry, error) {
	proEntryMap := make(map[int]*ProcessEntry)
	if pidList, err := process.Pids(); err != nil {
		return nil, err
	} else {
		for _, v := range pidList {
			proEntryMap[int(v)] = getEntryInfo(int(v))
		}
		return proEntryMap, nil
	}
}

func getProcessEntry(pid int) (*ProcessEntry, error) {
	return getEntryInfo(pid), nil
}

func getEntryInfo(pid int) *ProcessEntry {
	var proEntry ProcessEntry
	proEntry.Pid = pid

	pro, err := newProcess(pid)
	if err != nil {
		return &proEntry
	}
	defer pro.releaseHandle()

	if proMem, err := pro.getMemory(); err == nil {
		proEntry.ProcessMemory = *proMem
	}

	if proIo, err := pro.getIoCounters(); err == nil {
		proEntry.ProcessIoInfo = *proIo
	}

	if proTimes, err := pro.getTimes(); err == nil {
		proEntry.ProcessTime = *proTimes
	}

	if ppid, err := pro.getProcessParentId(); err == nil {
		proEntry.Ppid = ppid
	}

	if name, err := pro.getName(); err == nil {
		proEntry.Name = name
	}

	if threads, err := pro.getThreads(); err == nil {
		proEntry.Threads = threads
	}

	if handles, err := pro.getHandles(); err == nil {
		proEntry.Handles = handles
	}

	if bit, err := pro.getBit(); err == nil {
		proEntry.Bit = bit
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

	if fileInfo, err := getProcessFileInfo(proEntry.Exe); err == nil {
		proEntry.Size = fileInfo.Size
		proEntry.FileModifyTime = fileInfo.FileModifyTime
	}

	return &proEntry
}

func getProcessFileInfo(name string) (*ProcessFileInfo, error) {
	var fileInfo ProcessFileInfo
	statInfo, err := os.Stat(name)
	if err != nil {
		return &fileInfo, err
	}
	fileInfo.Size = statInfo.Size()
	fileInfo.FileModifyTime = statInfo.ModTime().Format("2006-01-02 15:04:05")
	return &fileInfo, nil
}

func formatCpuTime(times int64) int64 {
	return int64(times/int64(time.Second))
}