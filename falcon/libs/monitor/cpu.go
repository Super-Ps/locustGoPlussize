package monitor

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"math"
	"reflect"
	"strconv"
)

type CpuTime struct {
	TotalTime 			float64			`json:"cpu_total_time"`
	IdleTime	 		float64			`json:"idle_time"`
	BusyTime	 		float64			`json:"busy_time"`
	UnixNanoStamp 		int64			`json:"unix_nano_stamp"`
}

var (
	cpuInfo []cpu.InfoStat
	cpuError error
)


func GetCpuCount() (int, error) {
	if count, err := getCpuCount(); err != nil {
		return cpu.Counts(false)
	} else {
		return count, nil
	}
}

func GetCpuThreads() (int, error) {
	return cpu.Counts(true)
}

func GetCpuVendorId() (string, error) {
	if cpuError != nil {
		return "", cpuError
	}
	if len(cpuInfo) == 0 {
		return "", errors.New("Can not found cpu")
	}
	return cpuInfo[0].VendorID, nil
}

func GetCpuPhysicalId() (string, error) {
	if cpuError != nil {
		return "", cpuError
	}
	if len(cpuInfo) == 0 {
		return "", errors.New("Can not found cpu")
	}
	physicalId := cpuInfo[0].PhysicalID
	if physicalId == "0" {
		return "", nil
	} else {
		return physicalId, nil
	}
}

func GetCpuModelName() (string, error) {
	if cpuError != nil {
		return "", cpuError
	}
	if len(cpuInfo) == 0 {
		return "", errors.New("Can not found cpu")
	}
	return cpuInfo[0].ModelName, nil
}

func GetCpuMhz() (int, error) {
	if cpuError != nil {
		return 0, cpuError
	}
	if len(cpuInfo) == 0 {
		return 0, errors.New("Can not found cpu")
	}
	return int(cpuInfo[0].Mhz/100)*100, nil
}

func GetCpuTimes() (*CpuTime, error) {
	var cpuTime CpuTime
	if info , err := cpu.Times(true); err == nil {
		for _, data := range info {
			key := reflect.TypeOf(data)
			value := reflect.ValueOf(data)
			for i := 0; i < key.NumField(); i++ {
				if value.Field(i).Type().String() != "float64" {
					continue
				}
				cpuTime.TotalTime += value.Field(i).Float()
			}
			cpuTime.IdleTime += data.Idle
		}
		cpuTime.BusyTime = cpuTime.TotalTime-cpuTime.IdleTime
	} else {
		return nil, err
	}
	cpuTime.UnixNanoStamp = GetUnixNanoTime()
	return &cpuTime, nil
}

func GetCpuPercent(p1 *CpuTime, p2 *CpuTime) float64 {
	if (p2.TotalTime - p1.TotalTime) <= 0 {
		return 0
	}
	per, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", (p2.BusyTime-p1.BusyTime)/(p2.TotalTime-p1.TotalTime)*100), 64)
	return math.Min(math.Max(per,0), 100)
}

func (h CpuTime) String() string {
	s, _ := json.Marshal(h)
	return string(s)
}
