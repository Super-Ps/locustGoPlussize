package monitor

import (
	"encoding/json"
	"errors"
	"fmt"
	pNet "github.com/shirou/gopsutil/net"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type NetworkIo struct {
	BytesRecv 			int64			`json:"bytes_recv"`
	BytesSent 			int64			`json:"bytes_sent"`
	PacketsRecv 		int64			`json:"packets_recv"`
	PacketsSent 		int64			`json:"packets_sent"`
	UnixNanoStamp 		int64			`json:"unix_nano_stamp"`
}

type NetworkIoDiff struct {
	BytesRecvDiff 		int64			`json:"net_bytes_recv_diff"`
	BytesSentDiff 		int64			`json:"net_bytes_sent_diff"`
	PacketsRecvDiff 	int64			`json:"net_packets_recv_diff"`
	PacketsSentDiff 	int64			`json:"net_packets_sent_diff"`
}

type NetCardInfo struct {
	NetCardName 		string			`json:"net_card_name"`
	IPAddress 			string			`json:"ip_address"`
	MaskAddress 		string			`json:"mask_address"`
	MacAddress 			string			`json:"mac_address"`
	NetSpeed 			int				`json:"net_speed"`
}

type PublicResponse struct {
	Code 			int					`json:"code"`
	Data 			PublicData			`json:"data"`
}

type PublicData struct {
	Ip 				string				`json:"ip"`
}

var netMap, netError = getNetCardMap()
var publicIpAddress string


func GetNetWorkIo() (*NetworkIo, error) {
	if ioCounters, err := pNet.IOCounters(false); err != nil {
		return nil, err
	} else {
		if len(ioCounters) == 0 {
			return nil, errors.New("Can not found any netcard")
		}
		var ioCount NetworkIo
		for _, v := range ioCounters {
			ioCount.BytesRecv += int64(v.BytesRecv)
			ioCount.BytesSent += int64(v.BytesSent)
			ioCount.PacketsRecv += int64(v.PacketsRecv)
			ioCount.PacketsSent += int64(v.PacketsSent)
		}
		ioCount.UnixNanoStamp = GetUnixNanoTime()
		return &ioCount, nil
	}
}

func GetNetCardList() ([]*NetCardInfo, error) {
	netInfo, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	netCardMap, _ := getPhysicalNetCardMap()

	var netList []*NetCardInfo
	for _, netCard := range netInfo {
		netFlags := strings.Split(netCard.Flags.String(), "|")
		var isUp bool
		for _, name := range netFlags {
			if name != "up" {
				continue
			}
			isUp = true
			break
		}
		if !isUp {
			continue
		}
		addrInfo, err := netCard.Addrs()
		if err != nil {
			continue
		}
		if netCard.Name == "" {
			continue
		}

		if netCardMap != nil {
			if _, ok := netCardMap[netCard.Name]; !ok {
				continue
			}
		}

		var ipAddress string
		var maskAddress string
		for _, addr := range addrInfo {
			ip, ok := addr.(*net.IPNet)
			if !ok || ip.IP.IsLoopback() {
				continue
			}
			if ip.IP.To4() == nil {
				continue
			}
			if ip.IP.String()[0:3] == "169" || ip.IP.String()[0:3] == "127" {
				continue
			}
			ipAddress = ip.IP.String()
			maskAddress = formatMaskAddress(ip.Mask.String())
			break
		}
		if ipAddress == "" || maskAddress == "" {
			continue
		}

		netSpeed, _ := GetNetCardSpeed(netCard.Name)
		netList = append(netList, &NetCardInfo{
			NetCardName: netCard.Name,
			IPAddress: ipAddress,
			MaskAddress: maskAddress,
			MacAddress: strings.ToUpper(netCard.HardwareAddr.String()),
			NetSpeed: netSpeed})
	}

	if len(netList) == 0 {
		return nil, errors.New("Can not found any netcard")
	}

	return netList, nil
}

func GetNetCardInfo(name string) (*NetCardInfo, error) {
	if netError != nil {
		return nil, netError
	}
	if _, ok := netMap[name]; !ok {
		return nil, errors.New(fmt.Sprintf("Can not found netcard:%s", name))
	}
	return netMap[name], nil
}

func GetNetCardSpeed(name string) (int, error) {
	return getNetCardSpeed(name)
}

func GetPublicIpAddress() (string, error) {
	if publicIpAddress != "" {
		return publicIpAddress, nil
	}
	var jsonBody PublicResponse
	var isSuccess bool
	urlList := []string{"http://ip.taobao.com/service/getIpInfo.php?ip=myip"}
	for _, url := range urlList {
		client := http.Client{Timeout: 1*time.Second}
		res, err := client.Get(url)
		if err != nil {
			continue
		}
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			continue
		}

		if err := json.Unmarshal(body, &jsonBody); err != nil {
			continue
		}

		if jsonBody.Code != 0 {
			continue
		}
		isSuccess = true
		break
	}
	if isSuccess {
		publicIpAddress = jsonBody.Data.Ip
		return publicIpAddress, nil
	}
	return "", errors.New("Can not get public ip")
}

func GetNetworkIoDiff(io1 *NetworkIo, io2 *NetworkIo) *NetworkIoDiff {
	times := float64((io2.UnixNanoStamp-io1.UnixNanoStamp)/int64(time.Second))
	if times <= 0 {
		times = float64(1)
	}

	return &NetworkIoDiff{
		BytesRecvDiff: FormatNegativeInt64(int64(float64(io2.BytesRecv-io1.BytesRecv)/times)),
		BytesSentDiff: FormatNegativeInt64(int64(float64(io2.BytesSent-io1.BytesSent)/times)),
		PacketsRecvDiff: FormatNegativeInt64(int64(float64(io2.PacketsRecv-io1.PacketsRecv)/times)),
		PacketsSentDiff: FormatNegativeInt64(int64(float64(io2.PacketsSent-io1.PacketsSent)/times))}
}

func getNetCardMap() (map[string]*NetCardInfo, error) {
	netList, err := GetNetCardList()
	if err != nil {
		return nil, err
	}

	netCard := make(map[string]*NetCardInfo)
	for _, v := range netList{
		netCard[v.NetCardName] = v
	}
	return netCard, nil
}

func formatMaskAddress(addr string) string {
	var mask string
	for i := 0; i < 8; i+=2 {
		n, _ := strconv.ParseUint(addr[i:i+2],16,32)
		mask += fmt.Sprintf("%s.", strconv.Itoa(int(n)))
	}
	return strings.Trim(mask, ".")
}

func (h NetworkIo) String() string {
	s, _ := json.Marshal(h)
	return string(s)
}

func (h NetworkIoDiff) String() string {
	s, _ := json.Marshal(h)
	return string(s)
}

func (h NetCardInfo) String() string {
	s, _ := json.Marshal(h)
	return string(s)
}

func (h PublicResponse) String() string {
	s, _ := json.Marshal(h)
	return string(s)
}

func (h PublicData) String() string {
	s, _ := json.Marshal(h)
	return string(s)
}