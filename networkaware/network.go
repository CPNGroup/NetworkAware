package networkaware

import (
	"fmt"
	"strconv"

	"github.com/go-ping/ping"
)

func GetLatency() (map[string]string, error) {
	LatencyMap := make(map[string]string)
	// 需要改为所有需要测试节点的IP地址
	targetIP := []string{"10.129.32.84", "10.129.13.55", "10.129.173.253", "10.129.184.238"}
	for index, ip := range targetIP {
		// 创建ping对象
		pinger, err := ping.NewPinger(ip)
		if err != nil {
			fmt.Printf("无法创建Pinger: %v\n", err)
			return LatencyMap, err
		}

		// 设置为ipv4
		pinger.SetPrivileged(true)
		// 设置Ping次数为1次
		pinger.Count = 1

		// 运行ping操作
		err = pinger.Run()
		if err != nil {
			fmt.Printf("Ping失败: %v\n", err)
			return LatencyMap, err
		} else {
			stats := pinger.Statistics()
			LatencyMap["/latency/node1-node"+strconv.Itoa(index+1)] = stats.AvgRtt.String()
		}
	}
	return LatencyMap, nil
}
