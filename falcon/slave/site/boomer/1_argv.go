package boomer

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
)

type ArgvParam struct {
	MasterHost 			string
	MasterPort 			int
	ConfigPath 			string
	PidPath 			string
	AfterTime  			int64
	LocalConfig 		bool
	LogLevel 			string
	LogPath 			string
	LogUpload			bool
	LogAppend			bool
	GcTime				int64
	MinWait             int64
	MaxWait             int64
	Duration         	int64
	MonitorInterval 	int64
	HeartbeatInterval   int64
	RendezvousInterval  int64
	RunTimes            int64
	RandomMode          bool
	SlaveMode          	bool
	KeepAlive          	bool
	logLevelNum 		int
}

var Param = &ArgvParam{}


//加载命令行参数
func LoadArgv() {
	defer func() {
		if err := recover(); err != nil {
			os.Exit(1)
		}
	}()
	argv := flag.NewFlagSet(os.Args[0], 2)
	argv.StringVar(&Param.MasterHost, "master-host", "127.0.0.1", "--master-host= Host or IP address of locust master for distributed load testing")
	argv.IntVar(&Param.MasterPort, "master-port", 5557, "--master-port= The port to connect to that is used by the locust master for distributed load testing")
	argv.StringVar(&Param.ConfigPath, "config", "../conf/main.yml", "--config= Config file path")
	argv.StringVar(&Param.PidPath, "pid", "../pid/slave.pid", "--pid= Pid output path for this process")
	argv.StringVar(&Param.LogPath, "logfile", "", "--logfile= Path to log file. If not set, log will go to stdout/stderr")
	argv.StringVar(&Param.LogLevel, "loglevel", "INFO", "--loglevel= Choose between CRITICAL/ERROR/WARNING/INFO/DEBUG")
	argv.Int64Var(&Param.AfterTime, "after", 0, "--after= Set number of milliseconds before running wait")
	argv.BoolVar(&Param.LocalConfig, "local", false, "--local Set used local config file only")
	argv.BoolVar(&Param.LogAppend, "logappend", false, "--logappend Set append log mode")
	argv.BoolVar(&Param.LogUpload, "logupload", true, "--logupload Set auto upload log mode")
	argv.Int64Var(&Param.MinWait, "minwait", 0, "--minwait= Before next run to min wait milliseconds")
	argv.Int64Var(&Param.MaxWait, "maxwait", 0, "--maxwait= Before next run to max wait milliseconds")
	argv.Int64Var(&Param.Duration, "duration", 0, "--duration= Set number of milliseconds to stop,zero is never auto stop")
	argv.Int64Var(&Param.GcTime, "gctime-interval", 300000, "--gctime-interval= Set number of milliseconds to gc time interval. minimum to 60000")
	argv.Int64Var(&Param.MonitorInterval, "monitor-interval", 3000, "--monitor-interval= Set number of milliseconds monitor to master interval,zero is not send")
	argv.Int64Var(&Param.RendezvousInterval, "rendezvous-interval", 1, "--rendezvous-interval= Set number of milliseconds rendezvous to master interval,zero is not send")
	argv.Int64Var(&Param.RunTimes, "runtimes", 0, "--runtimes= Set number of runtimes to stop,zero is never auto stop")
	argv.BoolVar(&Param.RandomMode, "random", false, "--random Set task mode is random")
	argv.BoolVar(&Param.SlaveMode, "slave", true, "--SlaveMode Set running mode with this process as slave")
	argv.BoolVar(&Param.KeepAlive, "keepalive", false, "--KeepAlive When heartbeat time out will be auto exit")
	if !argv.Parsed() {
		_ = argv.Parse(os.Args[1:])
	}

	if Param.GcTime < 60 {
		Param.GcTime = 60
	}

	Param.ConfigPath, _ = filepath.Abs(Param.ConfigPath)
	Param.PidPath, _ = filepath.Abs(Param.PidPath)
	Param.ConfigPath = strings.Replace(Param.ConfigPath, "\\", "/", -1)
	Param.PidPath = strings.Replace(Param.PidPath, "\\", "/", -1)
	if Param.LogPath != "" {
		Param.LogPath, _ = filepath.Abs(Param.LogPath)
	}

	outputPid()
}
