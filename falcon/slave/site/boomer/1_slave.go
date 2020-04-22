package boomer

import (
	"bytes"
	"encoding/json"
	"falcon/libs/exception"
	"falcon/libs/monitor"
	"math/rand"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type SlaveEntry struct {
	Boom 					*Boomer
	Exception				*exception.Trier
}

type slaveInfo struct {
	systemInfo 				systemMonitor
	workerCount 			int64
	totalCount 				int64
	hatchRate 				int64
	hatchComplete 			int64
	weightCount 			int64
	caseConfig 				[]caseInfo
	caseFunc 				map[int64]map[string]reflect.Value
}

type systemMonitor struct {
	ClientID 				string		`json:"client_id"`
	ClientType 				string		`json:"client_type"`
	ClientVersion 			string		`json:"client_version"`
	HeartbeatInterval 		int64		`json:"heartbeat_interval"`
	monitor.OsInfo
}

type slaveControl struct {
	isRunning 				bool
	isSendStop 				bool
	log 					*sysLog
	startChannel			chan int
	startListChannel 		chan int
	createChannel 			chan int
	createListChannel 		chan int
	taskOverChannel 		chan int
	taskListChannel 		chan int
	stopChannel 			chan int
	stopListChannel 		chan int
	stopOverChannel 		chan int
	userChannel             chan *userInfo
}

var Slave *SlaveEntry
var info *slaveInfo
var control *slaveControl


//运行任务
func RunTasks() {
	go gc_monitor()
	time.Sleep(time.Duration(Param.AfterTime)*time.Millisecond)
	Slave = &SlaveEntry{
		Boom: NewBoomer(Param.MasterHost, Param.MasterPort),
		Exception: exception.New(),
	}

	info = &slaveInfo{
		systemInfo: systemMonitor{
			OsInfo: *monitor.GetOsInfo(),
		},
		caseFunc: make(map[int64]map[string]reflect.Value),
	}

	Conf = &SlaveConfig{
		BaseConfig:	LoadBaseConfig(),
		mutexInterval: time.Duration(1)*time.Millisecond,
		checkInterval: time.Duration(1)*time.Millisecond,
		breakMillisecond: int64(1000),
		restartWait: 5000,
	}

	control = &slaveControl{
		isRunning: false,
		log: &sysLog{},
		startChannel: make(chan int),
		startListChannel: make(chan int, 100000),
		createChannel: make(chan int),
		taskOverChannel: make(chan int),
		taskListChannel: make(chan int, 100000),
		stopChannel: make(chan int),
		stopListChannel: make(chan int, 100000),
		stopOverChannel: make(chan int),
		userChannel: make(chan *userInfo, 100000),
	}

	formatLogger()
	Slave.Boom.SetMode(DistributedMode)

	//Slave.Boom.Run(&Task{1, task, "task"})
	var tasks []*Task
	tasks = append(tasks, &Task{1, task, "task"})
	Slave.Boom.slaveRunner = newSlaveRunner(Slave.Boom.masterHost, Slave.Boom.masterPort, tasks, Slave.Boom.rateLimiter)
	for _, o := range Slave.Boom.outputs {
		Slave.Boom.slaveRunner.addOutput(o)
	}
	Slave.Boom.slaveRunner.getMessage()
	Slave.Boom.slaveRunner._run()
	time.Sleep(time.Duration(1)*time.Second)
	if Param.KeepAlive {
		go monitorHeartbeat()
	}
	info.systemInfo.ClientID = Slave.Boom.slaveRunner.nodeID
	info.systemInfo.ClientType = "go"
	info.systemInfo.ClientVersion = runtime.Version()[3:]
	info.systemInfo.HeartbeatInterval = int64(heartbeatInterval/time.Second)*1000
	sendSystemInfo()

	caseConfig := loadCaseConfig()
	entry := CaseEntry()
	caseFunc := loadCaseFunc(entry)
	workerEntry := loadWorkerEntry(entry)
	weightCount := loadCaseWeight()

	if weightCount == 0 {
		control.log.error("no case to run.")
		onExit()
	}
	for i := range caseConfig {
		name := caseConfig[i].name
		if checkEntryFunc(caseFunc, name) == false {
			control.log.error("case %s not exist.", name)
			onExit()
		}
	}
	if workerEntry == nil {
		control.log.error("can not find WorkerEntry.")
		onExit()
	}

	_ = Events.Subscribe("boomer:create", func(workers int, hatchRate int) {
		info.workerCount = int64(workers)
	})

	quitByMe := false
	_ = Events.Subscribe("boomer:quit", func() {
		if !quitByMe {
			control.log.info("shut down.")
			quit(1)
		}
	})
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
	quitByMe = true
	control.log.info("shut down by yourself.")
	//deletePidFile()
	quit(1)
}

//随机等待(毫秒)
func RandomWait(minWait time.Duration, maxWait time.Duration, seed... int64) {
	if (int64(minWait) == 0) && (int64(maxWait) == 0) {
		return
	}

	var randomSeed int64
	if len(seed) > 0 {
		randomSeed = seed[0]
	} else {
		randomSeed = 1
	}

	randTime := RandInt64(int64(minWait), int64(maxWait), randomSeed)
	if randTime < 0 {
		randTime = 0
	}
	time.Sleep(time.Duration(randTime))
}

//生成随机数
func RandInt64(min int64, max int64, seed int64) int64 {
	if (max < min) || (min == max) {
		return max
	}
	rad := rand.New(rand.NewSource(time.Now().UnixNano() + seed))
	return min + rad.Int63n(max-min+1)
}

//获取当前时间戳(毫秒)
func MilliNow() int64 {
	return time.Now().UnixNano()/int64(time.Millisecond)
}

//获取协程ID
func GetGoroutineId() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

func checkArgv(argv string) bool {
	for i := 1; i < len(os.Args); i++ {
		if len(os.Args[i]) < len(argv) {
			continue
		}
		if os.Args[i][0:len(argv)] == argv {
			return true
		}
	}
	return false
}

//获取默认配置文件路径
func GetDefaultConfigPath() string {
	absPath, _ := filepath.Abs(path.Join("..", "conf", "main.yml"))
	return strings.Replace(absPath, "\\", "/", -1)
}

//获取默认Pid文件路径
func GetDefaultPidPath() string {
	absPath, _ := filepath.Abs(path.Join("..", "pid", "slave.pid"))
	return strings.Replace(absPath, "\\", "/", -1)
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

//关闭ControlChannel
func closeControlChannel(c chan int) {
	defer func() {
		if err := recover(); err != nil {
		}
	}()
	close(c)
}

//退出slave
func quit(sleep int) {
	time.Sleep(time.Duration(sleep) * time.Second)
	go exit()
	Slave.Boom.Quit()
}

//退出程序
func exit() {
	time.Sleep(time.Duration(1) * time.Second)
	os.Exit(0)
}

//发送系统信息至Master事件
func sendSystemInfo() {
	data := make(map[string]interface{})
	data["info"], _ = json.Marshal(info.systemInfo)
	Slave.Boom.slaveRunner.client.sendChannel() <- newMessage("client_system", data, Slave.Boom.slaveRunner.nodeID)
}

//发送监控信息至Master事件
func sendMonitorInfo(info *monitorInfo) {
	data := make(map[string]interface{})
	data["info"], _ = json.Marshal(info)
	Slave.Boom.slaveRunner.client.sendChannel() <- newMessage("client_monitor", data, Slave.Boom.slaveRunner.nodeID)
}

//GC清理
func gc_monitor(){
	for {
		time.Sleep(time.Duration(Param.GcTime)*time.Millisecond)
		runtime.GC()
	}
}

//删除pid文件
func deletePidFile() {
	_, err := os.Stat(Param.PidPath)
	if err == nil {
		_ = os.Remove(Param.PidPath)
	}
}