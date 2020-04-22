package boomer

import (
	"time"
)

type WorkerEntry struct {
	Index       int64				//Worker在本Slave的序号
	Num         int64				//Worker在所有Slave中的序号
	Times       int64				//Worker在本Slave的当前运行次数
	IsRunning   *bool				//Slave当前运行状态
	Log 		*userLog			//日志输出方法
	Config 		*BoomerConfig		//配置文件内容
	Param 		*ArgvParam			//命令行参数
}

//向Master发送成功事件
func (w *WorkerEntry) RecordSuccess(requestType, name string, responseTime int64, responseLength int64) {
	Slave.Boom.RecordSuccess(requestType, name, responseTime, responseLength)
}

//向Master发送失败事件
func (w *WorkerEntry) RecordFailure(requestType, name string, responseTime int64, exception string) {
	Slave.Boom.RecordFailure(requestType, name, responseTime, exception)
}

//获取Slave分配的用户数
func (w *WorkerEntry) GetSlaveTotalUsers() int64 {
	return info.workerCount
}

//获取Slave当前启动的用户数
func (w *WorkerEntry) GetSlaveRunningUsers() int64 {
	return int64(Slave.Boom.slaveRunner.numClients)
}

//获取Slave的启动效率
func (w *WorkerEntry) GetSlaveHatchRate() int64 {
	return info.hatchRate
}

//获取Master的总用户数
func (w *WorkerEntry) GetMasterTotalUsers() int64 {
	return info.totalCount
}

//获取Master的当前已启动的用户数
func (w *WorkerEntry) GetMasterCompleteUsers() int64 {
	return info.hatchComplete
}

//获取Slave的ID
func (w *WorkerEntry) GetSlaveID() string {
	return Slave.Boom.slaveRunner.nodeID
}

//获取机器ID
func (w *WorkerEntry) GetMachineID() string {
	return info.systemInfo.HostID
}

//随机等待(毫秒)
func (w *WorkerEntry) RandomWait(minWait time.Duration, maxWait time.Duration) {
	RandomWait(minWait, maxWait, w.Times)
}

//生成随机数
func (w *WorkerEntry) RandInt64(min int64, max int64, index int64) int64 {
	return RandInt64(min, max, index)
}

//获取当前时间戳(毫秒)
func (w *WorkerEntry) MilliNow() int64 {
	return MilliNow()
}
