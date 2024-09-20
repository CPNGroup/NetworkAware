package main

import (
	"context"
	"fmt"
	"log"
	"pro/pb"
	"strconv"
	"time"

	"github.com/go-ping/ping"
	"google.golang.org/grpc"
)

var (
	// 需要改为所有需要测试节点的IP地址
	targetIP = []string{"10.129.32.84", "10.129.13.55", "10.129.173.253", "10.129.184.238"}
	address  = "10.129.184.238:30030" // zookeeper地址
	// 需要改为发出ping测试命令的节点IP，同样是要部署的节点IP
)

// node1, node2, node3, node4

func main() {
	//建立链接
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewZkServiceClient(conn)

	// 每隔30秒ping一次
	for {
		for index, ip := range targetIP {
			// 创建ping对象
			pinger, err := ping.NewPinger(ip)
			if err != nil {
				fmt.Printf("无法创建Pinger: %v\n", err)
				return
			}

			// 设置为ipv4
			pinger.SetPrivileged(true)
			// 设置Ping次数为1次
			pinger.Count = 1

			// 运行ping操作
			err = pinger.Run()
			if err != nil {
				fmt.Printf("Ping失败: %v\n", err)
			} else {
				stats := pinger.Statistics()
				// fmt.Printf("Nuc X Ping Nuc %v: 延迟=%v 丢包率=%.2f%%\n", index+1, stats.AvgRtt, stats.PacketLoss)
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				_, err = c.Set(ctx, &pb.PathAndData{Path: "/latency/node4-node" + strconv.Itoa(index+1), Data: stats.AvgRtt.String()})
				if err != nil {
					log.Printf("/latency/node4-node%v写入错误: %v", index+1, err)
				}
			}
		}
		// 等待30秒再进行下一次ping
		time.Sleep(30 * time.Second)
	}
}
