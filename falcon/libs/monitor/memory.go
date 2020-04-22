package monitor

import (
	"encoding/json"
	"github.com/shirou/gopsutil/mem"
)


type Memory struct {
	Rss 				int64				`json:"rss"`
	Vms 				int64				`json:"vms"`
}

func GetFreeMemory() (*Memory, error) {
	var memory Memory
	swap, err := GetFreeSwapMemory()
	if err != nil {
		return &memory, err
	} else {
		memory.Rss = swap
	}

	virtual, err := GetFreeVirtualMemory()
	if err != nil {
		return &memory, err
	} else {
		memory.Vms = virtual
	}

	return &memory, nil
}

func GetFreeSwapMemory() (int64, error) {
	if memInfo, err := mem.SwapMemory(); err != nil {
		return 0, err
	}else {
		return int64(memInfo.Free), nil
	}
}

func GetFreeVirtualMemory() (int64, error) {
	if memInfo, err := mem.VirtualMemory(); err != nil {
		return 0, err
	}else {
		return int64(memInfo.Free), nil
	}
}

func GetUsedMemory() (*Memory, error) {
	var memory Memory
	swap, err := GetUsedSwapMemory()
	if err != nil {
		return &memory, err
	} else {
		memory.Rss = swap
	}

	virtual, err := GetUsedVirtualMemory()
	if err != nil {
		return &memory, err
	} else {
		memory.Vms = virtual
	}

	return &memory, nil
}

func GetUsedSwapMemory() (int64, error) {
	if memInfo, err := mem.SwapMemory(); err != nil {
		return 0, err
	}else {
		return int64(memInfo.Used), nil
	}
}

func GetUsedVirtualMemory() (int64, error) {
	if memInfo, err := mem.VirtualMemory(); err != nil {
		return 0, err
	}else {
		return int64(memInfo.Used), nil
	}
}

func GetTotalMemory() (*Memory, error) {
	var memory Memory
	swap, err := GetTotalSwapMemory()
	if err != nil {
		return &memory, err
	} else {
		memory.Rss = swap
	}

	virtual, err := GetTotalVirtualMemory()
	if err != nil {
		return &memory, err
	} else {
		memory.Vms = virtual
	}

	return &memory, nil
}

func GetTotalSwapMemory() (int64, error) {
	if memInfo, err := mem.SwapMemory(); err != nil {
		return 0, err
	}else {
		return int64(memInfo.Total), nil
	}
}

func GetTotalVirtualMemory() (int64, error) {
	if memInfo, err := mem.VirtualMemory(); err != nil {
		return 0, err
	}else {
		return int64(memInfo.Total), nil
	}
}

func (h Memory) String() string {
	s, _ := json.Marshal(h)
	return string(s)
}