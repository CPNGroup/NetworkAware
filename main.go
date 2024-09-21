package main

import (
	"context"
	"log"
	"pro/networkaware"
	"pro/pb"
	"time"

	"google.golang.org/grpc"
)

var (
	address = "10.129.32.84:30030" // zookeeper地址
	// 需要改为发出ping测试命令的节点IP，同样是要部署的节点IP
)

func WriteIn() error {
	//建立链接
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
		return err
	}
	defer conn.Close()
	c := pb.NewZkServiceClient(conn)

	latency, err := networkaware.GetLatency()
	if err != nil {
		log.Printf("无法获取网络延迟: %v", err)
		return err
	}
	for key, value := range latency {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, err = c.Set(ctx, &pb.PathAndData{Path: key, Data: value})
		if err != nil {
			log.Printf("latency写入错误: %v", err)
			return err
		}

	}
	return nil
}

func main() {
	// 每隔30秒执行一次
	for {
		err := WriteIn()
		if err != nil {
			log.Printf("写入错误: %v", err)
		}

		// 等待30秒
		time.Sleep(30 * time.Second)
	}
}
