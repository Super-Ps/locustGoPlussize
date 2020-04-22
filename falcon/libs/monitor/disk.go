package monitor

import (
	"encoding/json"
)

type DiskSize struct {
	Total           		int64  				`json:"total"`
	Free            		int64  				`json:"free"`
	Used            		int64  				`json:"used"`
}

type DiskIoBase struct {
	ReadBytes 				int64			`json:"read_bytes"`
	WriteBytes				int64			`json:"write_bytes"`
	ReadCount				int64			`json:"read_count"`
	WriteCount				int64			`json:"write_count"`
	MergedReadCount			int64			`json:"merged_read_count"`
	MergedWriteCount		int64			`json:"merged_write_count"`
	UnixNanoStamp 			int64			`json:"unix_nano_stamp"`
}

type DiskIoTime struct {
	ReadTime				int64			`json:"io_read_time"`
	WriteTime				int64			`json:"io_write_time"`
	IoPsInProgress			int64			`json:"io_ps_in_progress"`
	IoTime					int64			`json:"io_time"`
	IoWeighted				int64			`json:"io_weighted"`
}

type DiskIoDiff struct {
	ReadBytesDiff 			int64			`json:"io_read_bytes_diff"`
	WriteBytesDiff 			int64			`json:"io_write_bytes_diff"`
	ReadCountDiff			int64			`json:"io_read_count_diff"`
	WriteCountDiff			int64			`json:"io_write_count_diff"`
	MergedReadCountDiff		int64			`json:"io_merged_read_count_diff"`
	MergedWriteCountDiff	int64			`json:"io_merged_write_count_diff"`
}

type DiskIoInfo struct {
	DiskIoBase
	DiskIoTime
}

type DiskDiff struct {
	DiskIoDiff
	DiskIoTime
}

type DiskSizeMap map[string]*DiskSize


func GetDiskIoInfo() (*DiskIoInfo, error) {
	return getDiskIoInfo()
}

func GetDiskSizeMap() (DiskSizeMap, error) {
	return getDiskSizeMap()
}

func GetTotalDiskSize() (*DiskSize, error) {
	diskMap, err := GetDiskSizeMap()
	if err != nil {
		return nil, err
	}

	var totalDisk DiskSize
	for _, v := range diskMap {
		totalDisk.Total += v.Total
		totalDisk.Used += v.Used
		totalDisk.Free += v.Free
	}
	return &totalDisk, nil
}

func GetDiskIoDiff(io1 *DiskIoBase, io2 *DiskIoBase) *DiskIoDiff {
	return getDiskIoDiff(io1, io2)
}

func (h DiskSize) String() string {
	s, _ := json.Marshal(h)
	return string(s)
}

func (h DiskIoBase) String() string {
	s, _ := json.Marshal(h)
	return string(s)
}

func (h DiskIoTime) String() string {
	s, _ := json.Marshal(h)
	return string(s)
}

func (h DiskIoDiff) String() string {
	s, _ := json.Marshal(h)
	return string(s)
}

func (h DiskIoInfo) String() string {
	s, _ := json.Marshal(h)
	return string(s)
}

func (h DiskDiff) String() string {
	s, _ := json.Marshal(h)
	return string(s)
}