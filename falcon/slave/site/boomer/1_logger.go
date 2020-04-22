package boomer

import (
	"falcon/libs/logger"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	LevelEmergency     = iota // 系统级紧急，比如磁盘出错，内存异常，网络不可用等
	LevelAlert                // 系统级警告，比如数据库访问异常，配置文件出错等
	LevelCritical             // 系统级危险，比如权限出错，访问异常等
	LevelError                // 用户级错误
	LevelWarning              // 用户级警告
	LevelInformational        // 用户级信息
	LevelDebug                // 用户级调试
	LevelTrace                // 用户级基本输出
)

const (
	SystemLog     = iota
	UserLog
)

type sysLog struct {}
type userLog struct {}

// 日志等级和描述映射关系
var LevelMap = map[string]int{
	"EMERGENCY": LevelEmergency,
	"ALERT":     LevelAlert,
	"CRITICAL":  LevelCritical,
	"ERROR":     LevelError,
	"WARNING":   LevelWarning,
	"INFO":      LevelInformational,
	"DEBUG":     LevelDebug,
	"TRACE":     LevelTrace,
}

//func (l *userLog) Painc(f interface{}, v ...interface{}) {
//	logger.Painc(f, v...)
//}
//
//func (l *userLog) Fatal(f interface{}, v ...interface{}) {
//	logger.Fatal(f, v...)
//}

//func (l *userLog) Emer(f interface{}, v ...interface{}) {
//	logger.Emer(f, v...)
//	uploadUserLog(LevelEmergency, f, v...)
//}

//func (l *userLog) Alert(f interface{}, v ...interface{}) {
//	logger.Alert(f, v...)
//	uploadUserLog(LevelAlert, f, v...)
//}

func (l *userLog) Crit(f interface{}, v ...interface{}) {
	logger.Crit(f, v...)
	uploadUserLog(LevelCritical, f, v...)
}

func (l *userLog) Error(f interface{}, v ...interface{}) {
	logger.Error(f, v...)
	uploadUserLog(LevelError, f, v...)
}

func (l *userLog) Warn(f interface{}, v ...interface{}) {
	logger.Warn(f, v...)
	uploadUserLog(LevelWarning, f, v...)
}

func (l *userLog) Info(f interface{}, v ...interface{}) {
	logger.Info(f, v...)
	uploadUserLog(LevelInformational, f, v...)
}

func (l *userLog) Debug(f interface{}, v ...interface{}) {
	logger.Debug(f, v...)
	uploadUserLog(LevelDebug, f, v...)
}

//func (l *userLog) Trace(f interface{}, v ...interface{}) {
//	logger.Trace(f, v...)
//	uploadUserLog(LevelTrace, f, v...)
//}

//func (l *sysLog) painc(f interface{}, v ...interface{}) {
//	logger.PaincSys(f, v...)
//}
//
//func (l *sysLog) fatal(f interface{}, v ...interface{}) {
//	logger.FatalSys(f, v...)
//}

//func (l *sysLog) emer(f interface{}, v ...interface{}) {
//	logger.EmerSys(f, v...)
//	uploadSystemLog(LevelEmergency, f, v...)
//}

//func (l *sysLog) alert(f interface{}, v ...interface{}) {
//	logger.AlertSys(f, v...)
//	uploadSystemLog(LevelAlert, f, v...)
//}

func (l *sysLog) crit(f interface{}, v ...interface{}) {
	logger.CritSys(f, v...)
	uploadSystemLog(LevelCritical, f, v...)
}

func (l *sysLog) error(f interface{}, v ...interface{}) {
	logger.ErrorSys(f, v...)
	uploadSystemLog(LevelError, f, v...)
}

func (l *sysLog) warn(f interface{}, v ...interface{}) {
	logger.WarnSys(f, v...)
	uploadSystemLog(LevelWarning, f, v...)
}

func (l *sysLog) info(f interface{}, v ...interface{}) {
	logger.InfoSys(f, v...)
	uploadSystemLog(LevelInformational, f, v...)
}

func (l *sysLog) debug(f interface{}, v ...interface{}) {
	logger.DebugSys(f, v...)
	uploadSystemLog(LevelDebug, f, v...)
}

//func (l *sysLog) trace(f interface{}, v ...interface{}) {
//	logger.TraceSys(f, v...)
//	uploadSystemLog(LevelTrace, f, v...)
//}

//格式化日志
func formatLogger() {
	var logConfig string
	Param.logLevelNum = LevelMap[Param.LogLevel]

	if Param.LogPath == "" {
		logConfig = fmt.Sprintf(`{
					"Console": {
						"level": "%s",
						"color": true
					}}`, Param.LogLevel)
	} else {
		logPath := getFileAbsPath(Param.LogPath)
		logDirPath := filepath.Dir(logPath)
		logConfig = fmt.Sprintf(`{
					"Console": {
						"level": "%s",
						"color": true
					},
					"File": {
						"filename": "%s",
						"level": "%s",
						"daily": true,
						"maxlines": 1000000,
						"maxsize": 4096,
						"maxdays": -1,
						"append": true,
						"permit": "0660"
				}}`, Param.LogLevel, logPath, Param.LogLevel)
		if !isExist(logDirPath) {
			_ = os.MkdirAll(logDirPath, os.ModePerm)
		}

		if (isExist(logPath) == true) && (Param.LogAppend == false) {
			_ = os.Remove(logPath)
		}
	}

	logger.Reset()
	_ = logger.SetLogger(logConfig)
	return
}

func uploadSystemLog(logLevel int, f interface{}, v ...interface{}) {
	if logLevel > Param.logLevelNum {
		return
	}

	if Param.LogUpload {
		data := make(map[string]interface{})
		data["info"] = logger.LogToJson(info.systemInfo.HostID, Slave.Boom.slaveRunner.nodeID, logLevel, SystemLog, f, v...)
		Slave.Boom.slaveRunner.client.sendChannel() <- newMessage("system_log", data, Slave.Boom.slaveRunner.nodeID)
	}

}

func uploadUserLog(logLevel int, f interface{}, v ...interface{}) {
	if logLevel > Param.logLevelNum {
		return
	}

	if Param.LogUpload {
		data := make(map[string]interface{})
		data["info"] = logger.LogToJson(info.systemInfo.HostID, Slave.Boom.slaveRunner.nodeID, logLevel, UserLog, f, v...)
		Slave.Boom.slaveRunner.client.sendChannel() <- newMessage("user_log", data, Slave.Boom.slaveRunner.nodeID)
	}
}

func isExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	}
	return true
}

func getFileAbsPath(file string) string {
	path, _ := filepath.Abs(file)
	return  strings.Replace(path, "\\", "/", -1)
}
