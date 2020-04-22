// +build windows

package monitor

import (
	"errors"
	"github.com/StackExchange/wmi"
	"github.com/shirou/gopsutil/disk"
	"time"
	"unsafe"
	"golang.org/x/sys/windows"
)

func getDiskIoDiff(io1 *DiskIoBase, io2 *DiskIoBase) *DiskIoDiff {
	return getDiskIoDiffFromProcesses(io1, io2)
}

func getDiskIoDiffFromWmi(io1 *DiskIoBase, io2 *DiskIoBase) *DiskIoDiff {
	times := float64((io2.UnixNanoStamp-io1.UnixNanoStamp)/int64(time.Second))
	if times <= 0 {
		times = float64(1)
	}
	return &DiskIoDiff{
		ReadBytesDiff: FormatNegativeInt64(int64(io2.ReadBytes)),
		WriteBytesDiff: FormatNegativeInt64(int64(io2.WriteBytes)),
		ReadCountDiff: FormatNegativeInt64(int64(io2.ReadCount)),
		WriteCountDiff: FormatNegativeInt64(int64(io2.WriteCount)),
		MergedReadCountDiff: FormatNegativeInt64(int64(io2.MergedReadCount)),
		MergedWriteCountDiff: FormatNegativeInt64(int64(io2.MergedWriteCount)),
	}
}

func getDiskIoDiffFromProcesses(io1 *DiskIoBase, io2 *DiskIoBase) *DiskIoDiff {
	times := float64((io2.UnixNanoStamp-io1.UnixNanoStamp)/int64(time.Second))
	if times <= 0 {
		times = float64(1)
	}
	return &DiskIoDiff{
		ReadBytesDiff: FormatNegativeInt64(int64(float64(io2.ReadBytes-io1.ReadBytes)/times)),
		WriteBytesDiff: FormatNegativeInt64(int64(float64(io2.WriteBytes-io1.WriteBytes)/times)),
		ReadCountDiff: FormatNegativeInt64(int64(float64(io2.ReadCount-io1.ReadCount)/times)),
		WriteCountDiff: FormatNegativeInt64(int64(float64(io2.WriteCount-io1.WriteCount)/times)),
		MergedReadCountDiff: FormatNegativeInt64(int64(float64(io2.MergedReadCount-io1.MergedReadCount)/times)),
		MergedWriteCountDiff: FormatNegativeInt64(int64(float64(io2.MergedWriteCount-io1.MergedWriteCount)/times)),
	}
}

func getDiskSizeMap() (DiskSizeMap, error) {
	diskSizeMap := make(DiskSizeMap)
	if diskPart, err := disk.Partitions(true); err != nil {
		return nil, err

	} else {
		for _, diskInfo := range diskPart {
			if diskSize, err := disk.Usage(diskInfo.Device); err == nil {
				diskSizeMap[diskInfo.Device] = &DiskSize{int64(diskSize.Total), int64(diskSize.Free), int64(diskSize.Used)}
			}
		}
		if len(diskSizeMap) == 0 {
			return nil, errors.New("Can not found any disk")
		}
		return diskSizeMap, nil
	}
}

func getDiskIoInfo() (*DiskIoInfo, error) {
	return getDiskIoFromProcesses()
}

func getDiskIoInfoFromWmi() (*DiskIoInfo, error) {
	type wmiData struct{
		AvgDiskReadQueueLength 		int64
		AvgDiskWriteQueueLength 	int64
		AvgDiskBytesPerRead 		int64
		AvgDiskBytesPerWrite 		int64
		AvgDisksecPerRead 			int64
		AvgDisksecPerWrite 			int64
	}
	var wmiList []wmiData

	err := wmi.Query("Select AvgDiskReadQueueLength,AvgDiskWriteQueueLength,AvgDiskBytesPerRead,AvgDiskBytesPerWrite,AvgDisksecPerRead,AvgDisksecPerWrite from Win32_PerfFormattedData_PerfDisk_LogicalDisk where Name='_Total'", &wmiList)
	if err != nil {
		return nil, err
	}

	if len(wmiList) == 0 {
		return nil, errors.New("Can not found any disk")
	}

	diskIoInfo := &DiskIoInfo{}
	diskIoInfo.ReadBytes = int64(wmiList[0].AvgDiskBytesPerRead)
	diskIoInfo.WriteBytes =  int64(wmiList[0].AvgDiskBytesPerWrite)
	diskIoInfo.ReadCount =  int64(wmiList[0].AvgDiskReadQueueLength)
	diskIoInfo.WriteCount =  int64(wmiList[0].AvgDiskWriteQueueLength)
	diskIoInfo.ReadTime =  int64(wmiList[0].AvgDisksecPerRead)
	diskIoInfo.WriteTime =  int64(wmiList[0].AvgDisksecPerWrite)
	diskIoInfo.MergedReadCount =  0
	diskIoInfo.MergedWriteCount =  0
	diskIoInfo.IoPsInProgress =  0
	diskIoInfo.IoWeighted =  0
	diskIoInfo.IoTime =  0
	diskIoInfo.UnixNanoStamp = GetUnixNanoTime()

	return diskIoInfo, nil
}

func getDiskIoFromProcesses() (*DiskIoInfo, error) {
	var diskIO DiskIoInfo
	if handle, err := openSnapshot(); err != nil {
		return nil, err
	} else {
		defer closeHandle(handle)
		proEntry := windows.ProcessEntry32{}
		proEntry.Size = uint32(unsafe.Sizeof(proEntry))
		if err := windows.Process32First(handle, &proEntry); err != nil {
			return nil, err
		}
		for {
			if err := windows.Process32Next(handle, &proEntry); err != nil {
				break
			} else {
				hPro, err := newProcess(int(proEntry.ProcessID))
				if err != nil {
					continue
				}
				io, err := hPro.GetIoCounters()
				if err != nil {
					hPro.CloseHandle()
					continue
				}
				diskIO.ReadBytes += io.IoReadBytes
				diskIO.WriteBytes += io.IoWriteBytes
				diskIO.ReadCount += io.IoReadCount
				diskIO.WriteCount += io.IoWriteCount
				hPro.CloseHandle()
			}
		}
		diskIO.UnixNanoStamp = GetUnixNanoTime()
		return &diskIO, nil
	}
}