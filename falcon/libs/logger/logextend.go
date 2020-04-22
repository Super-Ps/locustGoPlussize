package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

const (
	SystemLog     = iota
	UserLog
)

var logTypeList = [UserLog + 1]string{
	"SYSTEM",
	"USER",
}

type logJson struct {
	HostID  	string		`json:"host_id"`
	ClientID	string		`json:"client_id"`
	Time    	string		`json:"time"`
	Level   	string		`json:"level"`
	Type 		string		`json:"type"`
	Fn      	string		`json:"fn"`
	Path    	string		`json:"path"`
	Content 	string		`json:"content"`
}

func PaincSys(f interface{}, v ...interface{}) {
	defaultLogger.PanicSys(formatLog(f, v...))
}

func FatalSys(f interface{}, v ...interface{}) {
	defaultLogger.FatalSys(formatLog(f, v...))
}

func EmerSys(f interface{}, v ...interface{}) {
	defaultLogger.EmerSys(formatLog(f, v...))
}

func AlertSys(f interface{}, v ...interface{}) {
	defaultLogger.AlertSys(formatLog(f, v...))
}

func CritSys(f interface{}, v ...interface{}) {
	defaultLogger.CritSys(formatLog(f, v...))
}

func ErrorSys(f interface{}, v ...interface{}) {
	defaultLogger.ErrorSys(formatLog(f, v...))
}

func WarnSys(f interface{}, v ...interface{}) {
	defaultLogger.WarnSys(formatLog(f, v...))
}

func InfoSys(f interface{}, v ...interface{}) {
	defaultLogger.InfoSys(formatLog(f, v...))
}

func DebugSys(f interface{}, v ...interface{}) {
	defaultLogger.DebugSys(formatLog(f, v...))
}

func TraceSys(f interface{}, v ...interface{}) {
	defaultLogger.TraceSys(formatLog(f, v...))
}

func (this *LocalLogger) FatalSys(format string, args ...interface{}) {
	this.EmerSys("###Exec Panic:"+format, args...)
	os.Exit(1)
}

func (this *LocalLogger) PanicSys(format string, args ...interface{}) {
	this.EmerSys("###Exec Panic:"+format, args...)
	panic(fmt.Sprintf(format, args...))
}

// Emer Log EMERGENCY level message.
func (this *LocalLogger) EmerSys(format string, v ...interface{}) {
	this.writeMsgSys(LevelEmergency, format, v...)
}

// Alert Log ALERT level message.
func (this *LocalLogger) AlertSys(format string, v ...interface{}) {
	this.writeMsgSys(LevelAlert, format, v...)
}

// Crit Log CRITICAL level message.
func (this *LocalLogger) CritSys(format string, v ...interface{}) {
	this.writeMsgSys(LevelCritical, format, v...)
}

// Error Log ERROR level message.
func (this *LocalLogger) ErrorSys(format string, v ...interface{}) {
	this.writeMsgSys(LevelError, format, v...)
}

// Warn Log WARNING level message.
func (this *LocalLogger) WarnSys(format string, v ...interface{}) {
	this.writeMsgSys(LevelWarning, format, v...)
}

// Info Log INFO level message.
func (this *LocalLogger) InfoSys(format string, v ...interface{}) {
	this.writeMsgSys(LevelInformational, format, v...)
}

// Debug Log DEBUG level message.
func (this *LocalLogger) DebugSys(format string, v ...interface{}) {
	this.writeMsgSys(LevelDebug, format, v...)
}

// Trace Log TRAC level message.
func (this *LocalLogger) TraceSys(format string, v ...interface{}) {
	this.writeMsgSys(LevelTrace, format, v...)
}

func (this *LocalLogger) writeMsgSys(logLevel int, msg string, v ...interface{}) error {
	if !this.init {
		this.SetLogger(AdapterConsole)
	}
	msgSt := new(loginfo)
	src := ""
	if len(v) > 0 {
		msg = fmt.Sprintf(msg, v...)
	}
	when := time.Now()
	pc, file, lineno, ok := runtime.Caller(this.callDepth)
	fn := runtime.FuncForPC(pc).Name()
	fnList := strings.Split(fn, "/")
	fnName := fnList[len(fnList)-1]

	var strim string = "src/"
	if this.usePath != "" {
		strim = this.usePath
	}
	if ok {
		src = strings.Replace(
			fmt.Sprintf("%s:%d", stringTrim(file, strim), lineno), "%2e", ".", -1)
	}

	msgSt.Level = levelPrefix[logLevel]
	msgSt.Path = src
	msgSt.Content = msg
	msgSt.Fn = fnName
	msgSt.Name = this.appName
	msgSt.Time = when.Format(this.timeFormat)
	msgSt.Type = "SYSTEM"
	this.writeToLoggers(when, msgSt, logLevel)
	return nil
}

func LogToJson(host string, client string, logLevel int, logType int, f interface{}, v ...interface{}) []byte {
	if !defaultLogger.init {
		defaultLogger.SetLogger(AdapterConsole)
	}
	msgSt := new(logJson)
	src := ""
	var msg string
	if len(v) > 0 {
		msg = formatLog(f, v...)
	} else {
		msg = f.(string)
	}
	when := time.Now()
	pc, file, lineno, ok := runtime.Caller(3)
	fn := runtime.FuncForPC(pc).Name()
	fnList := strings.Split(fn, "/")
	fnName := fnList[len(fnList)-1]

	var strim string = "src/"
	if defaultLogger.usePath != "" {
		strim = defaultLogger.usePath
	}
	if ok {
		src = strings.Replace(
			fmt.Sprintf("%s:%d", stringTrim(file, strim), lineno), "%2e", ".", -1)
	}

	msgSt.HostID = host
	msgSt.ClientID = client
	msgSt.Level = levelPrefix[logLevel]
	msgSt.Path = src
	msgSt.Content = msg
	msgSt.Fn = fnName
	msgSt.Time = when.Format(defaultLogger.timeFormat)
	msgSt.Type = logTypeList[logType]
	msgJson, err := json.Marshal(msgSt)
	if err != nil {
		fmt.Println("Umarshal failed:", err)
		return nil
	}
	return msgJson
}