// +build windows

package monitor

import (
	"errors"
	"fmt"
	"github.com/StackExchange/wmi"
)


func getNetCardSpeed(name string) (int, error) {
	return getNetCardSpeedFromWmi(name)
}

func getNetCardSpeedFromWmi(name string) (int, error) {
	type wimData struct{
		Speed 		int64
	}
	var wmiList []wimData

	err := wmi.Query(fmt.Sprintf("Select Speed from Win32_NetworkAdapter where NetConnectionID='%s'", name), &wmiList)
	if err != nil {
		return 0, err
	}

	if len(wmiList) == 0 {
		return 0, errors.New(fmt.Sprintf("Can not found NetConnectionID:%s", name))
	}
	return int(wmiList[0].Speed/1000000), nil
}

func getPhysicalNetCardMap() (map[string]string, error) {
	return getPhysicalFromWmi()
}

func getPhysicalFromWmi() (map[string]string, error) {
	type wimData struct{
		NetConnectionID string
		PNPDeviceID 	string
	}
	var wmiList []wimData

	err := wmi.Query("Select NetConnectionID,PNPDeviceID from Win32_NetworkAdapter", &wmiList)
	if err != nil {
		return nil, err
	}

	if len(wmiList) == 0 {
		return nil, errors.New("Can not found NetCard")
	}

	netCardMap := make(map[string]string)
	for _, v := range wmiList {
		if (len(v.PNPDeviceID) < 3) || (len(v.NetConnectionID) == 0) {
			continue
		}
		if v.PNPDeviceID[0:3] != "PCI" && v.PNPDeviceID[0:3] != "USB" {
			continue
		}
		netCardMap[v.NetConnectionID] = v.PNPDeviceID
	}
	if len(netCardMap) == 0 {
		return nil, errors.New("Can not found Physical NetCard")
	}

	return netCardMap, nil
}