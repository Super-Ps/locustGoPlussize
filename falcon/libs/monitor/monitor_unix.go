//+build linux freebsd darwin openbsd

package monitor

import "time"

func getDiskDiff(interval... time.Duration) (*DiskDiff, error) {
	var intervalTime time.Duration
	if len(interval) == 0 {
		intervalTime = time.Duration(1)*time.Second
	} else {
		intervalTime = interval[0]
	}

	info1, err := GetDiskIoInfo()
	if err != nil {
		return nil, err
	}
	time.Sleep(intervalTime)
	info2, err := GetDiskIoInfo()
	if err != nil {
		return nil, err
	}
	ioDiff := GetDiskIoDiff(&info1.DiskIoBase, &info2.DiskIoBase)
	var diffInfo DiskDiff
	diffInfo.DiskIoDiff = *ioDiff
	diffInfo.DiskIoTime = info2.DiskIoTime
	return &diffInfo, nil
}