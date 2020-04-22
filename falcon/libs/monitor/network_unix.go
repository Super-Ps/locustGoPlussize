//+build linux freebsd darwin openbsd

package monitor

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)


func getNetCardSpeed(name string) (int, error) {
	return getNetCardSpeedFromCmd(name)
}

func getNetCardSpeedFromCmd(name string) (int, error) {
	out, err := CmdCall(fmt.Sprintf("ethtool %s | grep -o -i -P 'Speed:\\s*\\d*'", name))
	if err != nil {
		return 0, err
	}

	outSplit := strings.Split(out, ":")
	if len(outSplit) < 2 {
		return 0, errors.New(fmt.Sprintf("Can not found netcard:%s", name))
	}

	countString := strings.TrimSpace(outSplit[1])
	count, err := strconv.Atoi(strings.Trim(countString,""))
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

func getPhysicalNetCardMap() (map[string]string, error) {
	return nil, nil
}