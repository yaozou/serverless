package main

import (
	pb "com/aliyun/serverless/nodeservice/proto"
	"context"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
	"net"
	"time"
)

type NodeService struct {
}

//预定
func (s NodeService) Reserve(ctx context.Context, in *pb.ReserveRequest) (*pb.ReserveReply, error) {
	fmt.Printf("call function: NodeService.Reserve, %v\n", in)
	time.Sleep(time.Millisecond * 500)
	res := new(pb.ReserveReply)
	return res, nil
}

//创建容器
func (s NodeService) CreateContainer(ctx context.Context, in *pb.CreateContainerRequest) (*pb.CreateContainerReply, error) {
	fmt.Printf("call function: NodeService.CreateContainer, %v\n", in)
	res := new(pb.CreateContainerReply)
	res.ContainerId = uuid.NewV4().String()
	return res, nil
}

//销毁容器
func (s NodeService) RemoveContainer(ctx context.Context, in *pb.RemoveContainerRequest) (*pb.RemoveContainerReply, error) {
	fmt.Printf("call function: NodeService.RemoveContainer, %v\n", in)
	res := new(pb.RemoveContainerReply)
	return res, nil
}

//调用函数
func (s NodeService) InvokeFunction(in *pb.InvokeFunctionRequest, out pb.NodeService_InvokeFunctionServer) error {
	fmt.Printf("call function: NodeService.InvokeFunction, %v\n", in)
	return nil
}

//得到容器状态
func (s NodeService) GetStats(ctx context.Context, in *pb.GetStatsRequest) (*pb.GetStatsReply, error) {
	fmt.Printf("call function: NodeService.GetStats, %v\n", in)
	res := new(pb.GetStatsReply)
	res.NodeStats = &pb.NodeStats{TotalMemoryInBytes: 3 * 1024 * 1024 * 1024, MemoryUsageInBytes: 128 * 1024 * 1024}
	return res, nil
}

func main() {
	listen, err := net.Listen("tcp", "127.0.0.1:30000")
	if err != nil {
		fmt.Println(err)
	}

	//实现gRPC服务
	s := grpc.NewServer()
	//注册HelloServer为客户端提供服务
	pb.RegisterNodeServiceServer(s, NodeService{})

	fmt.Println("Node Service Listen on 127.0.0.1:30000")
	s.Serve(listen)
}
