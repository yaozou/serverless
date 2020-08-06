package main

import (
	"com/aliyun/serverless/scheduler/client"
	"com/aliyun/serverless/scheduler/handler"
	pb "com/aliyun/serverless/scheduler/proto"
	"com/aliyun/serverless/scheduler/server"
	"com/aliyun/serverless/scheduler/utils/groble"
	"fmt"
	"golang.org/x/sys/unix"
	"google.golang.org/grpc"
	"net"
	"os"

	//"com/aliyun/serverless/scheduler/server"
)

func main() {
	pid := unix.Getpid()
	fmt.Printf("pid=%v", pid)
	InitResourceMainEndpoint()
	go handler.AcquireContainerHandler() //启动容器请求处理器
	go handler.ReturnContainerHandler()  //启动容器归还处理器
	go handler.NodeHandler()             //启动Node管理处理器
	go handler.ContainerHandler()        //启动容器处理器
	client.ConnectResourceManagerService(groble.ResourceManagerEndpoint)
	StartSchedulerService()
}

//获取环境变量，资源管理器的地址
func InitResourceMainEndpoint() {
	endpoint := os.Getenv("RESOURCE_MANAGER_ENDPOINT")
	fmt.Println(endpoint)
	if endpoint == "" {
		panic("environment variable RESOURCE_MANAGER_ENDPOINT is not set")
	}
	fmt.Println("get resource manager endpoint is :" + endpoint)
	groble.ResourceManagerEndpoint = endpoint
}

//启动Scheduler服务
func StartSchedulerService() {
	fmt.Println("Hello GoLang")
	listen, err := net.Listen("tcp", groble.SchedulerServerAddress)
	if err != nil {
		fmt.Println(err)
	}
	network := listen.Addr().Network()
	fmt.Println(network)

	addr := listen.Addr()
	fmt.Println(addr.String())
	//实现gRPC服务
	s := grpc.NewServer()
	//注册HelloServer为客户端提供服务
	pb.RegisterSchedulerServer(s, new(server.Server))
	fmt.Println("Listen on " + groble.SchedulerServerAddress)
	//listen.Accept()
	//fmt.Println("connection success ")
	s.Serve(listen)
	fmt.Println("----------------------end--------------------")
}
