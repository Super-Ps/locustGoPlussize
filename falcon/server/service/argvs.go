package service

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type ArgvParam struct {
	Name 			string
	Ip 				string
	Port 			int
	GcTime 			int
	PauseTime 		int
	IntervalTime 	int
	Mode 			string
	PidPath 		string
	Route			RoutePath
	NoOut			bool
}

type RoutePath struct {
	Root 			string
	Monitor 		string
	Address			string
	GetDemo			string
	PostDemo		string
	Restart 		string
	Stop 			string
	Quit			string
	Static			string
}

var Param = &ArgvParam{}

func LoadArgv() {
	defer func() {
		if err := recover(); err != nil {
			os.Exit(1)
		}
	}()
	argv := flag.NewFlagSet(os.Args[0], 2)
	argv.StringVar(&Param.Ip, "ip", "0.0.0.0", "--ip= Bind ip address")
	argv.IntVar(&Param.Port, "port", 5555, "--port= Http listen port")
	argv.IntVar(&Param.GcTime, "gctime", 60, "--gctime= Seconds to GC interval minimum to 60")
	argv.IntVar(&Param.PauseTime, "pausetime", 60, "--pausetime= Seconds to monitor interval, minimum to 60")
	argv.IntVar(&Param.IntervalTime, "intervaltime", 3, "--intervalTime= Seconds to acquisition interval, minimum to 1")
	argv.StringVar(&Param.Name, "name", "server", "--name= show name")
	argv.StringVar(&Param.Route.Root, "route-root", "/", "--route-root= Route for root path")
	argv.StringVar(&Param.Route.Monitor, "route-monitor", "/monitor", "--route-monitor= Route for get service monitor data")
	argv.StringVar(&Param.Route.Address, "route-address", "/address", "--route-address= Route for get client address")
	argv.StringVar(&Param.Route.GetDemo, "route-getdemo", "/getdemo", "--route-getdemo= Route for getdemo")
	argv.StringVar(&Param.Route.PostDemo, "route-postdemo", "/postdemo", "--route-postdemo= Route for postdemo")
	argv.StringVar(&Param.Route.Restart, "route-restart", "/restart", "--route-restart= Route for restart")
	argv.StringVar(&Param.Route.Stop, "route-stop", "/stop", "--route-stop= Route for stop")
	argv.StringVar(&Param.Route.Quit, "route-quit", "/quit", "--route-quit= Route for quit")
	argv.StringVar(&Param.Mode, "mode", gin.ReleaseMode, fmt.Sprintf("--mode= (%s/%s/%s)", gin.ReleaseMode, gin.DebugMode, gin.TestMode))
	argv.StringVar(&Param.PidPath, "pid", "../pid/httpmonitor.pid", "--pid= Pid output path for server")
	argv.BoolVar(&Param.NoOut, "noout", false, "--noout No out put for recv request")
	if !argv.Parsed() {
		_ = argv.Parse(os.Args[1:])
	}
	if Param.Mode != gin.ReleaseMode && Param.Mode != gin.DebugMode && Param.Mode !=  gin.TestMode {
		Printf("Invalid mode '%s', Please choose mode in (%s/%s/%s)",Param.Mode, gin.ReleaseMode, gin.DebugMode, gin.TestMode)
		os.Exit(1)
	}
	if Param.GcTime < 60 {
		Param.GcTime = 60
	}
	if Param.PauseTime < 60 {
		Param.PauseTime = 60
	}
	if Param.IntervalTime < 1 {
		Param.IntervalTime = 1
	}
	Param.PidPath = strings.TrimSpace(Param.PidPath)
	if Param.PidPath != "" {
		Param.PidPath, _ = filepath.Abs(Param.PidPath)
		Param.PidPath = strings.Replace(Param.PidPath, "\\", "/", -1)
		outputPid()
	}
	formatRoute()
}

func formatRoute() {
	Param.Route.Root = fmt.Sprintf("/%s", strings.Trim(Param.Route.Root, "/"))
	Param.Route.Monitor = fmt.Sprintf("/%s", strings.Trim(fmt.Sprintf("%s/%s", strings.Trim(Param.Route.Root, "/"), strings.Trim(Param.Route.Monitor, "/")), "/"))
	Param.Route.Address = fmt.Sprintf("/%s", strings.Trim(fmt.Sprintf("%s/%s", strings.Trim(Param.Route.Root, "/"), strings.Trim(Param.Route.Address, "/")), "/"))
	Param.Route.GetDemo = fmt.Sprintf("/%s", strings.Trim(fmt.Sprintf("%s/%s", strings.Trim(Param.Route.Root, "/"), strings.Trim(Param.Route.GetDemo, "/")), "/"))
	Param.Route.PostDemo = fmt.Sprintf("/%s", strings.Trim(fmt.Sprintf("%s/%s", strings.Trim(Param.Route.Root, "/"), strings.Trim(Param.Route.PostDemo, "/")), "/"))
	Param.Route.Restart = fmt.Sprintf("/%s", strings.Trim(fmt.Sprintf("%s/%s", strings.Trim(Param.Route.Root, "/"), strings.Trim(Param.Route.Restart, "/")), "/"))
	Param.Route.Stop = fmt.Sprintf("/%s", strings.Trim(fmt.Sprintf("%s/%s", strings.Trim(Param.Route.Root, "/"), strings.Trim(Param.Route.Stop, "/")), "/"))
	Param.Route.Quit = fmt.Sprintf("/%s", strings.Trim(fmt.Sprintf("%s/%s", strings.Trim(Param.Route.Root, "/"), strings.Trim(Param.Route.Quit, "/")), "/"))
	Param.Route.Static = fmt.Sprintf("/%s", strings.Trim(fmt.Sprintf("%s/%s", strings.Trim(Param.Route.Root, "/"), "static"), "/"))
}

//输出当前Pid到文件
func outputPid() {
	pidDir := filepath.Dir(Param.PidPath)
	_, err := os.Stat(pidDir)
	if err != nil {
		err = os.MkdirAll(pidDir, os.ModePerm)
		if err != nil {
			return
		}
	}
	f, err := os.OpenFile(Param.PidPath,os.O_RDWR|os.O_CREATE|os.O_APPEND,0644)
	if err != nil {
		return
	}
	defer f.Close()
	_, _ = f.Write([]byte(strconv.Itoa(os.Getpid()) + "\n"))
}


