package boomer

import (
	"falcon/libs/monitor"
	"os"
	"os/exec"
	"runtime"
	"time"
)

var caseCount = int64(0)

//启动用户事件
func onHatch(msg *message) {
	startTime = MilliNow()
	rate, _ := msg.Data["hatch_rate"]
	clients, _ := msg.Data["num_clients"]
	hatchRate := int(rate.(float64))
	workers := 0
	if _, ok := clients.(uint64); ok {
		workers = int(clients.(uint64))
	} else {
		workers = int(clients.(int64))
	}
	if workers == 0 || hatchRate == 0 {
		control.log.info("Invalid hatch message from master, num_clients is %d, hatch_rate is %d\n",
			workers, hatchRate)
	} else {
		control.log.info("Hatching and swarming", workers, "clients at the rate", hatchRate, "clients/s...")
	}
}

//初始化用户事件
func onCreate(msg *message) {
	defer func() {
		if err := recover(); err != nil {
			control.log.error("onCreate fail:%s", err)
		}
	}()
	if control.isRunning {
		return
	}
	Conf.BaseConfig = LoadBaseConfig()
	control.isRunning = true
	control.isSendStop = false
	start, _ := msg.Data["start_index"]
	client, _ := msg.Data["num_clients"]
	total, _ := msg.Data["total_clients"]
	hatch, _ := msg.Data["hatch_rate"]
	index := int64(0)
	num := int64(0)
	all := int64(0)
	rate := int64(0)
	if _, ok := start.(uint64); ok {
		index = int64(start.(uint64))
	} else {
		index = int64(start.(int64))
	}

	if _, ok := client.(uint64); ok {
		num = int64(client.(uint64))
	} else {
		num = int64(client.(int64))
	}

	if _, ok := total.(uint64); ok {
		all = int64(total.(uint64))
	} else {
		all = int64(total.(int64))
	}

	if _, ok := hatch.(uint64); ok {
		rate = int64(hatch.(uint64))
	} else {
		rate = int64(hatch.(int64))
	}

	info.totalCount = all
	info.hatchRate = rate
	info.workerCount = num
	info.weightCount = loadCaseWeight()
	info.caseConfig = loadCaseConfig()
	caseCount = int64(len(info.caseConfig))
	for i := int64(0); i < num; i++ {
		entry := CaseEntry()
		workerEntry := loadWorkerEntry(entry)
		caseFunc := loadCaseFunc(entry)
		info.caseFunc[i+1] = caseFunc
		control.userChannel <- &userInfo{i + 1, index + i, entry, workerEntry}
	}
	close(control.createChannel)
}

//心跳事件
func onHeartbeat() {
	lastHeartbeat = MilliNow()
}

//用户启动完成事件
func onHatchOver() {
	defer func() {
		if err := recover(); err != nil {
			control.log.error("onHatchOver fail:%s", err)
		}
	}()
	close(control.startChannel)
}

//用户执行完成事件
func onTaskOver() {
	defer func() {
		if err := recover(); err != nil {
			control.log.error("onTaskOver fail:%s", err)
		}
	}()
	if control.isRunning == false {
		return
	}
	closeControlChannel(control.taskOverChannel)
	control.taskOverChannel = make(chan int)
}

//修改配置事件
func onConfig(msg *message) {
	if Param.LocalConfig {
		return
	}
	conf := LoadDataConfig(msg.Data["config"].([]byte))
	if conf != nil {
		Conf.masterConfig = conf
		Conf.BaseConfig = conf
	}
}

//更新当前master已启动的用户数
func onComplete(msg *message) {
	complete, _ := msg.Data["complete"]
	count := int64(0)
	if _, ok := complete.(uint64); ok {
		count = int64(complete.(uint64))
	} else {
		count = int64(complete.(int64))
	}
	info.hatchComplete = count
}

//监控性能事件
func onMonitor() {
	monitorInterval := Param.MonitorInterval
	if monitorInterval < 1000 {
		return
	}
	go getMonitorInfo(monitorInterval)
}

//停止任务事件
func onStop() {
	if !control.isRunning {
		return
	}
	control.log.info("stopTask start")
	control.isRunning = false
	defer func() {
		if err := recover(); err != nil {
			control.log.error("stopTask fail:%s", err)
		}
		control.log.info("stop task over.")
		control.log.info("Recv stop message from master, all the goroutines are stopped")
	}()
	closeControlChannel(control.stopChannel)
	_ = <- control.stopOverChannel
	time.Sleep(time.Duration(1)*time.Second)
	closeControlChannel(control.startListChannel)
	closeControlChannel(control.taskListChannel)
	closeControlChannel(control.stopListChannel)
	closeControlChannel(control.startChannel)
	closeControlChannel(control.createChannel)
	closeControlChannel(control.taskOverChannel)
	close(control.userChannel)

	control.startChannel = make(chan int)
	control.startListChannel = make(chan int, 100000)
	control.createChannel = make(chan int)
	control.taskOverChannel = make(chan int)
	control.taskListChannel = make(chan int, 100000)
	control.stopChannel = make(chan int)
	control.stopListChannel = make(chan int, 100000)
	control.stopOverChannel = make(chan int)
	control.userChannel = make(chan *userInfo, 100000)

	runtime.GC()
	return
}

//退出任务事件
func onQuit() {
	onStop()
}

//退出事件
func onExit() {
	disconnectedMaster()
	//deletePidFile()
	time.Sleep(time.Duration(1)*time.Second)
	os.Exit(0)
}

//重启事件
func onRestart() {
	disconnectedMaster()
	exe, _ := monitor.GetProcessExePath(os.Getpid())
	dir, _ := os.Getwd()
	cmd := &exec.Cmd{
		Path: exe,
		Args: os.Args,
		Dir: dir,
	}
	_ = cmd.Start()
	time.Sleep(time.Duration(1)*time.Second)
	os.Exit(0)
}

//断开Master连接
func disconnectedMaster() {
	Slave.Boom.slaveRunner.client.sendChannel() <- newMessage("client_stopped", nil, Slave.Boom.slaveRunner.nodeID)
	Slave.Boom.slaveRunner.client.sendChannel() <- newMessage("quit", nil, Slave.Boom.slaveRunner.nodeID)
	Slave.Boom.slaveRunner.client.disconnectedChannel()
}