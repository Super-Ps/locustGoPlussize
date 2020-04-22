//+build linux freebsd darwin openbsd

package monitor

import (
	"errors"
	"strconv"
	"strings"
)

func getCpuCount() (int, error) {
	return getCpuCountFromCmd()
}

func getCpuCountFromCmd() (int, error) {
	physicalCpuCount, err := getPhysicalCpuCountFromCmd()
	if err != nil {
		return 0, err
	}

	coreCount, err := getCpuCoreCountFromCmd()
	if err != nil {
		return 0, err
	}
	return physicalCpuCount*coreCount, nil
}

func getPhysicalCpuCountFromCmd() (int, error) {
	out, err := CmdCall(`cat /proc/cpuinfo| grep "physical id"| sort| uniq| wc -l`)
	if err != nil {
		return 0, err
	}

	countString := strings.TrimSpace(out)
	count, err := strconv.Atoi(strings.Trim(countString,""))
	if err != nil {
		return 0, err
	}
	return count, nil
}

func getCpuCoreCountFromCmd() (int, error) {
	out, err := CmdCall("cat /proc/cpuinfo| grep 'cpu cores'| uniq")
	if err != nil {
		return 0, err
	}

	outSplit := strings.Split(out, ":")
	if len(outSplit) < 2 {
		return 0, errors.New("Can not found cpuinfo")
	}

	countString := strings.TrimSpace(outSplit[1])
	count, err := strconv.Atoi(strings.Trim(countString,""))
	if err != nil {
		return 0, err
	}
	return count, nil
}