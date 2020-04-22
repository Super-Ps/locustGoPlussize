//+build linux freebsd darwin openbsd

package monitor

import (
	"errors"
	"github.com/shirou/gopsutil/disk"
	"regexp"
	"strconv"
	"strings"
	"time"
)


func getDiskIoDiff(io1 *DiskIoBase, io2 *DiskIoBase) *DiskIoDiff {
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
	return getDiskSizeMapFromCmd()
}

func getDiskSizeMapFromCmd() (DiskSizeMap, error) {
	out, err := CmdCall(`df | awk '{print $1 " " $4 " " $2}'`)
	if err != nil {
		return nil, err
	}

	diskMap := make(DiskSizeMap)
	pathIndex := 1
	freeIndex := 2
	totalIndex := 3
	reString := `([^\s]+)\s+(\d+)\s+(\d+)`
	split := "\n"

	re := regexp.MustCompile(reString)
	outSplit := strings.Split(out, split)
	for i, data := range outSplit {
		if i == 0 {
			continue
		}
		dataLine := TrimString(data)
		diskSplit := re.FindAllStringSubmatch(dataLine, -1)
		if len(diskSplit) < 1 {
			continue
		}
		if len(diskSplit[0]) < 4 {
			continue
		}
		path := diskSplit[0][pathIndex]
		totalSize, _ := strconv.Atoi(strings.Trim(diskSplit[0][totalIndex],""))
		freeSize, _ := strconv.Atoi(strings.Trim(diskSplit[0][freeIndex],""))
		diskMap[path] = &DiskSize{
			Total: int64(totalSize)*1024,
			Free: int64(freeSize)*1024,
			Used: int64(totalSize-freeSize)*1024,
		}
	}
	return diskMap, nil
}

func getDiskIoInfo() (*DiskIoInfo, error) {
	var diskIoInfo DiskIoInfo
	var diskList []string
	if diskPart, err := disk.Partitions(true); err == nil {
		for _, diskInfo := range diskPart {
			diskList = append(diskList, diskInfo.Device)
		}
	}
	if len(diskList) == 0 {
		return nil, errors.New("Can not found any disk")
	}

	if diskIoMap, err := disk.IOCounters(diskList...) ; err == nil {
		for _, v := range diskIoMap {
			diskIoInfo.ReadBytes += int64(v.ReadBytes)
			diskIoInfo.WriteBytes +=  int64(v.WriteBytes)
			diskIoInfo.ReadCount +=  int64(v.ReadCount)
			diskIoInfo.WriteCount +=  int64(v.WriteCount)
			diskIoInfo.ReadTime +=  int64(v.ReadTime)
			diskIoInfo.WriteTime +=  int64(v.WriteTime)
			diskIoInfo.MergedReadCount +=  int64(v.MergedReadCount)
			diskIoInfo.MergedWriteCount +=  int64(v.MergedWriteCount)
			diskIoInfo.IoPsInProgress +=  int64(v.IopsInProgress)
			diskIoInfo.IoWeighted +=  int64(v.WeightedIO)
			diskIoInfo.IoTime +=  int64(v.IoTime)
		}
	}
	diskIoInfo.UnixNanoStamp = GetUnixNanoTime()
	return &diskIoInfo, nil
}