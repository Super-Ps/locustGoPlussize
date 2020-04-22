// +build windows

package monitor

import (
	"errors"
	"github.com/StackExchange/wmi"
)

func getCpuCount() (int, error) {
	return getCpuCountFromWmi()
}

func getCpuCountFromWmi() (int, error) {
	type wimData struct{
		NumberOfCores 		int
	}
	var wmiList []wimData

	err := wmi.Query("Select NumberOfCores from Win32_Processor", &wmiList)
	if err != nil {
		return 0, err
	}

	if len(wmiList) == 0 {
		return 0, errors.New("Can not found any cpus")
	}

	var cpuCount int
	for _, v := range wmiList {
		cpuCount += v.NumberOfCores
	}

	return cpuCount, nil
}