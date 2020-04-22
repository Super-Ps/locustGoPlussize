package boomer

import (
	"encoding/json"
	"falcon/libs/monitor"
	"math"
	"os"
	"runtime"
	"time"
)

type monitorInfo struct {
	HostID 				string			`json:"host_id"`
	ClientID 			string			`json:"client_id"`
	UserCount			int				`json:"user_count"`
	DiskUsed 			int64			`json:"disk_used"`
	monitor.ProcessDiff
	monitor.NetworkIoDiff
}

var startTime = int64(0)
var lastHeartbeat = MilliNow()

//监控心跳
func monitorHeartbeat() {
	for {
		time.Sleep(time.Duration(1)*time.Second)
		if (MilliNow()-lastHeartbeat) > (int64(heartbeatInterval/time.Second)*1000+3000) {
			control.log.info("Heartbeat timeout from master")
			onExit()
		}
		if (startTime > 0) && (Param.Duration > 0) && ((MilliNow()-startTime) >= Param.Duration) && (control.isRunning == true) && (control.isSendStop == false) && (Slave.Boom.slaveRunner.state == "running" || Slave.Boom.slaveRunner.state == "hatching") {
			control.isSendStop = true
			Slave.Boom.slaveRunner.client.sendChannel() <- newMessage("stop", nil, Slave.Boom.slaveRunner.nodeID)
			control.log.info("Duration over %dms to stop", Param.Duration)
			break
		}
	}

}

//监控用户启动事件
func monitorStartChannel() {
	control.log.info("monitorStartChannel start")
	defer func() {
		if err := recover(); err != nil {
			control.log.error("monitorStartChannel fail:%s", err)
		}
	}()

	userCount := int64(0)
	for {
		if !control.isRunning {
			break
		}
		select {
		case <- control.stopChannel:
			control.log.info("monitorStartChannel stop")
			runtime.Goexit()
		case <- control.startListChannel:
			userCount += 1
			if userCount >= info.workerCount {
				Slave.Boom.slaveRunner.client.sendChannel() <- newMessage("start", nil, Slave.Boom.slaveRunner.nodeID)
				control.log.info("monitorStartChannel stop")
				runtime.Goexit()
			}
		}
	}
}

//监控用户停止事件
func monitorStopChannel() {
	control.log.info("monitorStopChannel start")
	defer func() {
		if err := recover(); err != nil {
			control.log.error("monitorStopChannel fail:%s", err)
		}
	}()
	userCount := int64(0)
	for {
		select {
		case <- control.stopOverChannel:
			control.log.info("monitorStopChannel stop")
			runtime.Goexit()
		case <- control.stopListChannel:
			userCount += 1
			if userCount >= info.workerCount {
				if control.isRunning == true && control.isSendStop == false {
					time.Sleep(time.Duration(1)*time.Second)
					control.isSendStop = true
					Slave.Boom.slaveRunner.client.sendChannel() <- newMessage("stop", nil, Slave.Boom.slaveRunner.nodeID)
					time.Sleep(time.Duration(1)*time.Second)
				}
				closeControlChannel(control.stopOverChannel)
			}
		}
	}
}

//监控用户任务事件
func monitorTaskChannel() {
	control.log.info("monitorTaskChannel start")
	defer func() {
		if err := recover(); err != nil {
			control.log.error("monitorTaskChannel fail:%s", err)
		}
	}()
	userCount := int64(0)
	var task_times = int64(0)
	for {
		select {
		case <- control.stopChannel:
			if (userCount >= info.workerCount) && (len(control.taskListChannel) > 0) {
				if Param.RendezvousInterval > 0 {
					Slave.Boom.slaveRunner.client.sendChannel() <- newMessage("task", nil, Slave.Boom.slaveRunner.nodeID)
				} else {
					onTaskOver()

				}
			}
			control.log.info("monitorTaskChannel stop")
			runtime.Goexit()
		case <- control.taskListChannel:
			userCount += 1
			if userCount >= info.workerCount {
				userCount = int64(0)
				if task_times < math.MaxInt64 {
					task_times += 1
				}
				if Param.RendezvousInterval > 0 {
					time.Sleep(time.Duration(Param.RendezvousInterval)*time.Millisecond)
					Slave.Boom.slaveRunner.client.sendChannel() <- newMessage("task", nil, Slave.Boom.slaveRunner.nodeID)
				} else {
					time.Sleep(time.Duration(1)*time.Millisecond)
					onTaskOver()
				}
				if Param.RunTimes > 0 && task_times >= Param.RunTimes*caseCount && control.isSendStop == false {
					control.isSendStop = true
					Slave.Boom.slaveRunner.client.sendChannel() <- newMessage("stop", nil, Slave.Boom.slaveRunner.nodeID)
					control.log.info("Run times over %d to stop.", Param.RunTimes)
					break
				}
			}
		}
	}
}

//获取本地性能损耗
func getMonitorInfo(interval int64) {
	control.log.info("getMonitorInfo start")
	defer func() {
		if err := recover(); err != nil {
			control.log.info("getMonitorInfo fail:%s", err)
		}
	}()
	monitorChannel := make(chan int)
	close(monitorChannel)

	for {
		select {
		case <- control.stopChannel:
			control.log.info("getMonitorInfo stop")
			runtime.Goexit()
		case <- monitorChannel:
			if control.isRunning == false {
				control.log.info("getMonitorInfo stop")
				runtime.Goexit()
			}
			pro1, _ := monitor.GetProcessEntry(os.Getpid())
			net1, _ := monitor.GetNetWorkIo()
			time.Sleep(time.Duration(interval)*time.Millisecond)
			pro2, _ := monitor.GetProcessEntry(os.Getpid())
			net2, _ := monitor.GetNetWorkIo()
			diskTotal, _ := monitor.GetTotalDiskSize()

			sendMonitorInfo(&monitorInfo{
				HostID: info.systemInfo.HostID,
				ClientID: Slave.Boom.slaveRunner.nodeID,
				UserCount: int(Slave.Boom.slaveRunner.numClients),
				DiskUsed: diskTotal.Used,
				ProcessDiff: *monitor.GetProcessDiffInfo(pro1, pro2),
				NetworkIoDiff: *monitor.GetNetworkIoDiff(net1, net2),
			})
		}
	}
}

func (h monitorInfo) String() string {
	s, _ := json.Marshal(h)
	return string(s)
}
