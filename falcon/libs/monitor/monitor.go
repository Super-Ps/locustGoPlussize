package monitor

import (
	"time"
)


func GetSystemDiff(interval... time.Duration) (*SystemDiff, error) {
	var intervalTime time.Duration
	if len(interval) == 0 {
		intervalTime = time.Duration(1)*time.Second
	} else {
		intervalTime = interval[0]
	}

	var sysDiff SystemDiff
	io1, err := GetNetWorkIo()
	if err != nil {
		return &sysDiff, nil
	}

	disk1, err := GetDiskIoInfo()
	if err != nil {
		return &sysDiff, nil
	}

	cpu1, err := GetCpuTimes()
	if err != nil {
		return &sysDiff, nil
	}

	time.Sleep(intervalTime)

	cpu2, err := GetCpuTimes()
	if err != nil {
		return &sysDiff, nil
	} else {
		sysDiff.CpuPercent = GetCpuPercent(cpu1, cpu2)
		sysDiff.CpuTime = int64(int64(cpu2.TotalTime)/int64(osInfo.CpuThreads))
	}

	io2, err := GetNetWorkIo()
	if err != nil {
		return &sysDiff, nil
	} else {
		sysDiff.NetworkIoDiff = *GetNetworkIoDiff(io1, io2)
	}

	disk2, err := GetDiskIoInfo()
	if err != nil {
		return &sysDiff, nil
	}
	sysDiff.DiskIoDiff = *GetDiskIoDiff(&disk1.DiskIoBase, &disk2.DiskIoBase)

	proTotal, err := GetProcessesTotalInfo()
	if err != nil {
		return &sysDiff, nil
	} else {
		sysDiff.ProcessesTotalInfo = *proTotal
	}

	memTotal, err := GetUsedMemory()
	if err != nil {
		return &sysDiff, nil
	} else {
		sysDiff.Memory = *memTotal
	}

	diskTotal, err := GetTotalDiskSize()
	if err != nil {
		return &sysDiff, nil
	} else {
		sysDiff.DiskUsed = diskTotal.Used
	}

	runTime, err :=  GetRunTime()
	if err != nil {
		return &sysDiff, nil
	} else {
		sysDiff.RunTime = runTime
	}

	sysDiff.UnixNanoStamp = GetUnixNanoTime()
	sysDiff.LastTime = UnixTimeToDateTime(sysDiff.UnixNanoStamp)
	return &sysDiff, nil
}

func GetProcessDiff(pid int, interval... time.Duration) (*ProcessDiff, error) {
	var intervalTime time.Duration
	if len(interval) == 0 {
		intervalTime = time.Duration(1)*time.Second
	} else {
		intervalTime = interval[0]
	}

	p, err := OpenProcess(pid)
	if err != nil {
		return nil, err
	}
	defer p.CloseHandle()

	io1, err := p.GetIoCounters()
	if err != nil {
		return nil, err
	}

	time1, err := p.GetTimes()
	if err != nil {
		return nil, err
	}

	time.Sleep(intervalTime)

	io2, err := p.GetIoCounters()
	if err != nil {
		return nil, err
	}

	time2, err := p.GetTimes()
	if err != nil {
		return nil, err
	}

	proEntry, err := p.GetEntry()
	if err != nil {
		return nil, err
	}

	ioDiff := GetProcessIoDiff(io1, io2)
	cpuPer := GetProcessCpuPercent(time1, time2)
	return FormatProcessDiff(proEntry, ioDiff, cpuPer), nil
}

func GetProcessMapDiff(interval... time.Duration) (map[int]*ProcessDiff, error) {
	var intervalTime time.Duration
	if len(interval) == 0 {
		intervalTime = time.Duration(1)*time.Second
	} else {
		intervalTime = interval[0]
	}
	proMap1, err := GetProcessMap()
	if err != nil {
		return nil, err
	}
	time.Sleep(intervalTime)
	proMap2, err := GetProcessMap()
	if err != nil {
		return nil, err
	}

	proMapDiff := make(map[int]*ProcessDiff)
	for k, v := range proMap2 {
		if _, ok := proMap1[k]; !ok {
			continue
		}
		ioDiff := GetProcessIoDiff(&proMap1[k].ProcessIoInfo, &v.ProcessIoInfo)
		cpuPer := GetProcessCpuPercent(&proMap1[k].ProcessTime, &v.ProcessTime)
		proMapDiff[k] = FormatProcessDiff(v, ioDiff, cpuPer)
	}
	return proMapDiff, nil
}

func GetDiskDiff(interval... time.Duration) (*DiskDiff, error) {
	return getDiskDiff(interval...)
}

func GetNetworkDiff(interval... time.Duration) (*NetworkIoDiff, error) {
	var intervalTime time.Duration
	if len(interval) == 0 {
		intervalTime = time.Duration(1)*time.Second
	} else {
		intervalTime = interval[0]
	}

	io1, err := GetNetWorkIo()
	if err != nil {
		return nil, err
	}
	time.Sleep(intervalTime)
	io2, err := GetNetWorkIo()
	if err != nil {
		return nil, err
	}
	return GetNetworkIoDiff(io1, io2), nil
}

func GetCpuPercentDiff(interval... time.Duration) (float64, error) {
	var intervalTime time.Duration
	if len(interval) == 0 {
		intervalTime = time.Duration(1)*time.Second
	} else {
		intervalTime = interval[0]
	}

	cpu1, err := GetCpuTimes()
	if err != nil {
		return 0, err
	}
	time.Sleep(intervalTime)
	cpu2, err := GetCpuTimes()
	if err != nil {
		return 0, err
	}
	return GetCpuPercent(cpu1, cpu2), nil
}