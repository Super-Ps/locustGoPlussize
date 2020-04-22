package monitor

import (
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"os"
)

func init(){
	_init()
	curPid = os.Getpid()
	curPpid = os.Getppid()
	curExe = os.Args[0]
	curCmd = GetCurrentCmdline()
	curName, _ = GetProcessName(curPid)
	curUserName, _ = GetProcessUserName(curPid)
	curCwd, _ = os.Getwd()
	curBit, _ = GetProcessBit(curPid)
	systemInfo, systemError = host.Info()
	cpuInfo, cpuError = cpu.Info()
	osInfo = getOsInfo()
}