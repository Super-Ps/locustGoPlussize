// +build windows

package monitor

import "time"

func getDiskDiff(interval... time.Duration) (*DiskDiff, error) {
	info, err := GetDiskIoInfo()
	if err != nil {
		return nil, err
	}
	var diffInfo DiskDiff
	diffInfo.ReadBytesDiff = info.ReadBytes
	diffInfo.WriteBytesDiff = info.WriteBytes
	diffInfo.ReadCountDiff = info.ReadCount
	diffInfo.WriteCountDiff = info.WriteCount
	diffInfo.MergedReadCountDiff = info.MergedReadCount
	diffInfo.MergedWriteCountDiff = info.MergedWriteCount
	diffInfo.ReadTime = info.ReadTime
	diffInfo.WriteTime = info.WriteTime
	diffInfo.IoPsInProgress = info.IoPsInProgress
	diffInfo.IoTime = info.IoTime
	diffInfo.IoWeighted = info.IoWeighted
	return &diffInfo, nil
}