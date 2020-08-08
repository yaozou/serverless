package handler

import (
	"com/aliyun/serverless/scheduler/client"
	"com/aliyun/serverless/scheduler/core"
	"encoding/json"
	"fmt"
	"time"
)

/*
	node-manager负责探测node资源的使用率，
	当使用率高的时候就去申请资源，
	当使用率低的时候就释放资源
*/
//const ReservePress = 100                             //申请压力
//const ReleasePress = 0.3                             //释放压力
const AccountId = "1317891723692367"      //TODO 线上可能会变化
const MinNodeCount = 4                    //最少节点数量
const MaxNodeCount = 20                   //最大节点数量
const SleepTime = time.Millisecond * 2000 //睡眠时间
const ReservePress = 65                   //预定node的cpu压力
const ReleasePress = 30                   //释放node的cpu使用率
//const NodeSniffIntervalTime = time.Millisecond * 2000 //Node嗅探间隔时间

//MinNodeCount=a,MaxNodeCount=b
//(0,a)申请资源
//[a,a]只能申请资源
//(a,b)申请或者释放资源
//[b,)只能释放资源

func NodeHandler() {
	for {
		size := core.GetNodeCount()
		//(0,a)不满足最低要求，无条件直接申请资源
		if size < MinNodeCount {
			node := ReserveOneNode(5)
			core.AddNode(node)
			core.PrintNodes("reserve node ")
			continue
		}
		//press := core.CalcNodesPress() //计算节点压力
		press := SniffAllNodeAvgPress() //嗅探节点平均压力

		//[a,a]只能申请资源
		if size == MinNodeCount {
			if press > ReservePress {
				node := ReserveOneNode(2)
				core.AddNode(node)
				LoadFuncForNewNode(node) //为新节点加载历史函数
				fmt.Println(node)
			} else {
				time.Sleep(SleepTime)
			}
			continue
		}

		//(a,b)申请或者释放资源
		if size > MinNodeCount && size < MaxNodeCount {
			if press > ReservePress { //当压力达到0.7就申请一个node
				node := ReserveOneNode(2)
				LoadFuncForNewNode(node)
				core.AddNode(node)
				LoadFuncForNewNode(node) //为新节点加载历史函数
				fmt.Println(node)
			} else if press < ReleasePress { //当压力小于0.4就释放一个
				ReleaseOneNode()
			} else {
				time.Sleep(SleepTime)
			}
			continue
		}

		if size >= MaxNodeCount {
			if press < ReleasePress {
				ReleaseOneNode()
			} else {
				time.Sleep(SleepTime)
			}
			continue
		}
	}
}

//嗅探所有节点平均压力
func SniffAllNodeAvgPress() float64 {
	nodes := core.GetNodes()
	core.PrintNodes("local node status")
	var count = 0
	var totalPress float64 = 0
	fmt.Printf("****************************%v*******************************\n", "remote node stats")
	for _, n := range nodes {
		res := client.GetStats(n.Client, "")
		if res == nil {
			continue
		}
		count++
		totalPress += res.NodeStats.CpuUsagePct

		jsonStr, _ := json.Marshal(res)
		if jsonStr != nil {
			fmt.Println(string(jsonStr))
		}
	}
	fmt.Printf("**************************************************************\n\n")
	if count == 0 {
		return 0 //如果获取状态都失败，那么就直接返回0，表示没有压力
	}
	var avgPress = totalPress / float64(count) //计算平均压力
	return avgPress
}

func PrintNodeStats() {
	for {
		time.Sleep(time.Millisecond * 10000) //没10秒打印一次node状态
		nodes := core.GetNodes()
		core.PrintNodes("local node status")
		fmt.Printf("****************************%v*******************************\n", "remote node stats")
		for _, n := range nodes {
			reply := client.GetStats(n.Client, "")
			jsonStr, err := json.Marshal(reply)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println(string(jsonStr))
		}
		fmt.Printf("**************************************************************\n\n")
	}
}

//这个方法需要保证一定要申请一个Node,TODO 需要为节点实例话所已知的函数
func ReserveOneNode(collectionMaxCapacity int64) *core.Node {
	st := time.Now().UnixNano()
	for {
		//预约一个node
		reply, err := client.ReserveNode("", AccountId)
		if err != nil || reply == nil || reply.Node == nil {
			fmt.Println("error ", err)
			time.Sleep(time.Second * 1) //一秒过后再重试
			continue
		}

		//ReservedTimeTimestampMs ReleasedTimeTimestampMs
		nodeClient, err := client.ConnectNodeService(reply.Node.Id, reply.Node.Address, reply.Node.NodeServicePort)
		if err != nil {
			fmt.Println("error ", err)
			continue
		}
		//requestId := uuid.NewV4().String()
		//statsReply := client.GetStats(nodeClient, requestId)
		//totalMem := statsReply.GetNodeStats().TotalMemoryInBytes
		//usedMem := statsReply.GetNodeStats().MemoryUsageInBytes
		//创建成功node并且连接成功，进行节点添加
		node := core.NewNode(reply.Node.Id, reply.Node.Address, reply.Node.NodeServicePort, reply.Node.MemoryInBytes,
			0, nodeClient, collectionMaxCapacity)
		et := time.Now().UnixNano()
		fmt.Printf("---- reserve node, time=%v, node:%v \n", (et-st)/1000000, node)
		return node
	}
}

func ReleaseOneNode() {
	node := core.RemoveLastNode()
	time.Sleep(time.Millisecond * 5000) //5秒钟过后在释放
	client.ReleaseNode("", node.NodeID)
}
