package handler

import (
	"com/aliyun/serverless/scheduler/client"
	"com/aliyun/serverless/scheduler/core"
	pb "com/aliyun/serverless/scheduler/proto"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"sync"
	"time"
)

//异步创建容器
func AsyncCreateContainer(funcName string, handler string, timeoutInMs int64, memoryInBytes int64) {
	go CreateContainer(funcName, handler, timeoutInMs, memoryInBytes)
}
func CreateContainer(funcName string, handler string, timeoutInMs int64, memoryInBytes int64) {
	node := core.GetSuitableNode(funcName, memoryInBytes)
	if node == nil {
		return
	}

	st := time.Now().UnixNano()
	for {
		//创建一个container
		reply, err := client.CreateContainer(node.Client, uuid.NewV4().String(), funcName+uuid.NewV4().String(),
			funcName, handler, timeoutInMs, 4*1024*1024*1024)

		if err != nil {
			fmt.Printf("FuncName:%v, Mem:%v, error: %v", funcName, memoryInBytes/1048576, err)
			return
		}

		//将container添加到node中
		container := &core.Container{
			ContainerId:      reply.ContainerId,
			TotalMem:         4 * 1024 * 1024 * 1024,
			UsageMem:         128 * 1024 * 1024,
			CpuUsagePct:      0,
			FuncName:         funcName,
			UseCount:         0,
			ConcurrencyCount: 2,
			Node:             node,
		}

		node.AddContainer(container)
		core.AddContainer(container)

		et := time.Now().UnixNano()
		fmt.Printf("create container,FuncName:%v, Mem:%v, time=%v, nodeId=%v\n", funcName, memoryInBytes/1024/1024, (et-st)/1000000, node.NodeID)
		return
	}

}

//
//import (
//	"com/aliyun/serverless/scheduler/client"
//	"com/aliyun/serverless/scheduler/core"
//	pb "com/aliyun/serverless/scheduler/proto"
//	"fmt"
//	uuid "github.com/satori/go.uuid"
//	"math/rand"
//	"sync"
//	"time"
//)

/*
	处理container的加载
*/
//var funcQueue = make(chan *pb.AcquireContainerRequest, 40000)
var NodeMaxContainerCount = 15 //node加载container最大数量

//func AddAcquireContainerToContainerHandler(req *pb.AcquireContainerRequest) {
//	funcQueue <- req
//}

//var CreateContainerWG sync.WaitGroup

var FuncNameMap = make(map[string]*pb.AcquireContainerRequest) // New empty set
var FuncNameMapLock sync.Mutex

//var locker = new(sync.Mutex)
//var cond = sync.NewCond(locker)
//var IsChange = false

////添加req到map中
//func AddReq(req *pb.AcquireContainerRequest) {
//	FuncNameMapLock.Lock()
//	if FuncNameMap[req.FunctionName] == nil {
//		FuncNameMap[req.FunctionName] = req
//		go CreateContainerHandler(req)
//	}
//	FuncNameMapLock.Unlock()
//}

////复制一份出来
//func GetReq() map[string]*pb.AcquireContainerRequest {
//	tmp := make(map[string]*pb.AcquireContainerRequest) // New empty set
//	FuncNameMapLock.Lock()
//	for k, v := range FuncNameMap {
//		tmp[k] = v
//	}
//	FuncNameMapLock.Unlock()
//	return tmp
//}

//
//func CreateContainerHandler(req *pb.AcquireContainerRequest) {
//	nodes := core.GetNodes()
//	for i := 0; i < core.CollectionMaxCapacity; i++ {
//		var wg sync.WaitGroup
//		for _, node := range nodes { //为每个node添加函数
//			wg.Add(1)
//			go HandleFuncName(node, req, &wg)
//		}
//		wg.Wait()
//		randomTime := rand.Intn(60)
//		time.Sleep(time.Second * time.Duration(120+randomTime)) //睡眠一段时间再去创建第二个
//		//******************log*************************
//		//core.PrintNodes(" create container ")
//		//******************log*************************
//	}
//}
//
////为新node加载函数实例
//func LoadFuncForNewNode(node *core.Node, allWg *sync.WaitGroup) {
//	defer allWg.Done()
//	reqMap := GetReq()
//	//未每个node加载
//	var i int64 = 0
//	for ; i < node.CollectionMaxCapacity; i++ {
//		var wg sync.WaitGroup
//		for _, req := range reqMap {
//			wg.Add(1)
//			go HandleFuncName(node, req, &wg)
//		}
//		wg.Wait()
//		//******************log*************************
//		//core.PrintNodes(" create container ")
//		//******************log*************************
//	}
//
//}
//
////处理一个函数的加载
//func HandleFuncName(node *core.Node, req *pb.AcquireContainerRequest, wg *sync.WaitGroup) {
//	defer wg.Done()
//	//判断这个node是否缺乏这个函数实例
//	if node.Lack(req.FunctionName) {
//		container := CreateContainer(node, req)
//		node.AddContainer(container)
//	}
//}
//
//////保证创建一个container
//func CreateContainer(node *core.Node, req *pb.AcquireContainerRequest) *core.Container {
//	st := time.Now().UnixNano()
//	for {
//		//创建一个container
//		reply, err := client.CreateContainer(
//			node.Client,
//			req.RequestId,                          //demo是这样
//			req.FunctionName+uuid.NewV4().String(), //demo是这样
//			req.FunctionName,
//			req.FunctionConfig.Handler,
//			req.FunctionConfig.TimeoutInMs,
//			req.FunctionConfig.MemoryInBytes,
//		)
//
//		if err != nil {
//			fmt.Printf("FuncName:%v, Mem:%v, error: %v", req.FunctionName, req.FunctionConfig.MemoryInBytes/1048576, err)
//			return nil
//		}
//
//		//将container添加到node中
//		container := &core.Container{FuncName: req.FunctionName, Id: reply.ContainerId,
//			MaxUsedMem: req.FunctionConfig.MemoryInBytes, MaxUsedCount: core.DefaultMaxUsedCount}
//		et := time.Now().UnixNano()
//		fmt.Printf("create container,FuncName:%v, Mem:%v, time=%v, nodeId=%v\n", req.FunctionName, req.FunctionConfig.MemoryInBytes/1048576, (et-st)/1000000, node.NodeID)
//		return container
//	}
//}
