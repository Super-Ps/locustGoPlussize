package service

import (
	"encoding/json"
	"falcon/libs/monitor"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
	"sync"
	"time"
)

type MonitorMapResponse struct {
	Code 					int				`json:"code"`
	Data 					*MonitorMapData	`json:"data"`
}

type MonitorMapData struct {
	ServerName 				string			`json:"server_name"`
	*monitor.OsInfo
	Monitor 				*MonitorMapInfo	`json:"monitor"`
}

type MonitorMapInfo struct {
	System					*monitor.SystemDiff				`json:"system"`
	Process					map[int]*monitor.ProcessDiff	`json:"process"`
}

type MonitorListResponse struct {
	Code 					int					`json:"code"`
	Data 					*MonitorListData	`json:"data"`
}

type MonitorListData struct {
	ServerName 				string				`json:"server_name"`
	*monitor.OsInfo
	Monitor 				*MonitorListInfo	`json:"monitor"`
}

type MonitorListInfo struct {
	System					*monitor.SystemDiff				`json:"system"`
	Process					[]*monitor.ProcessDiff			`json:"process"`
}

var sysInfo = monitor.GetOsInfo()
var updateTime = monitor.GetUnixNanoTime()
var isStart bool
var monitorMapLock sync.RWMutex
var monitorMapInfo = &MonitorMapInfo{}
var monitorListLock sync.RWMutex
var monitorListInfo = &MonitorListInfo{}
var sysLock sync.RWMutex
var sysMonitor = &monitor.SystemDiff{}
var proMapLock sync.RWMutex
var proMapMonitor = make(map[int]*monitor.ProcessDiff)
var proListLock sync.RWMutex
var proListMonitor []*monitor.ProcessDiff


func (s *ServerEntry) HttpMonitorMap(ctx *gin.Context) {
	isMonitor := isStart
	updateTime = monitor.GetUnixNanoTime()
	if !isMonitor {
		time.Sleep(time.Duration(3)*time.Second)
	}
	monitorMapLock.RLock()
	defer monitorMapLock.RUnlock()

	var mapInfo MonitorMapInfo
	proMap := make(map[int]*monitor.ProcessDiff)
	keyWord := strings.TrimSpace(ctx.Query("keyword"))
	if keyWord != "" {
		keyWordList := strings.Split(keyWord, ",")
		if len(keyWordList) == 0 {
			mapInfo = *monitorMapInfo
		} else {
			mapInfo.System = monitorMapInfo.System
			for k, v := range monitorMapInfo.Process {
				for _, word := range keyWordList {
					if word == strconv.Itoa(v.Pid) || word == strconv.Itoa(v.Ppid) || strings.Index(v.Name, word) >= 0 || strings.Index(v.Cmd, word) >= 0 || strings.Index(v.Cwd, word) >= 0 || strings.Index(v.Exe, word) >= 0 {
						proMap[k] = v
						break
					}
				}
			}
			mapInfo.Process = proMap
		}
	} else {
		mapInfo = *monitorMapInfo
	}

	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Access-Control-Allow-Headers", "Content-Type")
	ctx.Header("content-type", "application/json")
	ctx.JSON(200, MonitorMapResponse{0, &MonitorMapData{Param.Name,sysInfo,&mapInfo}})
}

func (s *ServerEntry) HttpMonitorList(ctx *gin.Context) {
	isMonitor := isStart
	updateTime = monitor.GetUnixNanoTime()
	if !isMonitor {
		time.Sleep(time.Duration(3)*time.Second)
	}
	monitorListLock.RLock()
	defer monitorListLock.RUnlock()

	var listInfo MonitorListInfo
	var proList []*monitor.ProcessDiff
	keyWord := strings.TrimSpace(ctx.Query("keyword"))
	if keyWord != "" {
		keyWordList := strings.Split(keyWord, ",")
		if len(keyWordList) == 0 {
			listInfo = *monitorListInfo
		} else {
			listInfo.System = monitorListInfo.System
			for _, v := range monitorListInfo.Process {
				for _, word := range keyWordList {
					if word == strconv.Itoa(v.Pid) || word == strconv.Itoa(v.Ppid) || strings.Index(v.Name, word) >= 0 || strings.Index(v.Cmd, word) >= 0 || strings.Index(v.Cwd, word) >= 0 || strings.Index(v.Exe, word) >= 0 {
						proList = append(proList, v)
						break
					}
				}
			}
			listInfo.Process = proList
		}
	} else {
		listInfo = *monitorListInfo
	}

	//设置ContentType
	ctx.Header("Content-Type", "application/json")

	//允许跨域访问
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	ctx.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

	ctx.JSON(200, MonitorListResponse{0, &MonitorListData{Param.Name,sysInfo,&listInfo}})
}

func UpdateSysMonitor(interval time.Duration) {
	for {
		if !isStart {
			break
		}
		if info, err := monitor.GetSystemDiff(interval); err == nil {
			sysLock.Lock()
			sysMonitor = info
			sysLock.Unlock()
		}
	}
}

func UpdateProMapMonitor(interval time.Duration) {
	for {
		if !isStart {
			break
		}
		if info, err := monitor.GetProcessMapDiff(interval); err == nil {
			proMapLock.Lock()
			proMapMonitor = info
			proMapLock.Unlock()
		}
	}
}

func UpdateProListMonitor(interval time.Duration) {
	for {
		if !isStart {
			break
		}
		time.Sleep(interval)
		proMapLock.RLock()
		var proList []*monitor.ProcessDiff
		for _, v := range proMapMonitor {
			proList = append(proList, v)
		}
		proMapLock.RUnlock()

		proListLock.Lock()
		proListMonitor = proList
		proListLock.Unlock()
	}
}

func UpdateMonitorMapInfo(interval time.Duration) {
	for {
		if !isStart {
			break
		}
		time.Sleep(interval)
		monitorMapLock.Lock()
		sysLock.RLock()
		monitorMapInfo.System = sysMonitor
		sysLock.RUnlock()
		proMapLock.RLock()
		monitorMapInfo.Process = proMapMonitor
		proMapLock.RUnlock()
		if len(monitorMapInfo.Process) > 0 {
			monitorMapInfo.System.Processes = int64(len(monitorMapInfo.Process))
		}
		monitorMapLock.Unlock()
	}
}

func UpdateMonitorListInfo(interval time.Duration) {
	for {
		if !isStart {
			break
		}
		time.Sleep(interval)
		monitorListLock.Lock()
		sysLock.RLock()
		monitorListInfo.System = sysMonitor
		sysLock.RUnlock()
		proListLock.RLock()
		monitorListInfo.Process = proListMonitor
		proListLock.RUnlock()
		if len(monitorListInfo.Process) > 0 {
			monitorListInfo.System.Processes = int64(len(monitorListInfo.Process))
		}
		monitorListLock.Unlock()
	}
}

func UpdateStart()  {
	for {
		time.Sleep(time.Duration(1)*time.Second)
		if ((monitor.GetUnixNanoTime()-updateTime)/int64(time.Second) > int64(Param.PauseTime)) {
			if isStart {
				isStart = false
			}
			continue
		}
		if !isStart {
			isStart = true
			go UpdateSysMonitor(time.Duration(Param.IntervalTime)*time.Second)
			go UpdateProMapMonitor(time.Duration(Param.IntervalTime)*time.Second)
			go UpdateProListMonitor(time.Duration(1)*time.Second)
			go UpdateMonitorMapInfo(time.Duration(1)*time.Second)
			go UpdateMonitorListInfo(time.Duration(1)*time.Second)
		}
	}
}

func (h MonitorMapResponse) String() string {
	s, _ := json.Marshal(h)
	return string(s)
}

func (h MonitorMapData) String() string {
	s, _ := json.Marshal(h)
	return string(s)
}

func (h MonitorMapInfo) String() string {
	s, _ := json.Marshal(h)
	return string(s)
}

func (h MonitorListResponse) String() string {
	s, _ := json.Marshal(h)
	return string(s)
}

func (h MonitorListData) String() string {
	s, _ := json.Marshal(h)
	return string(s)
}

func (h MonitorListInfo) String() string {
	s, _ := json.Marshal(h)
	return string(s)
}