package core

import "sync"

//存放container信息
type Container struct {
	FunName   string //函数名字
	Id        string //容器id
	UsedCount int    //使用数量
	UsedMem   int64  //使用内存
	lock      sync.RWMutex
}

//得到容器使用内存大小
func (container *Container) GetUsedMem() int64 {
	container.lock.RLock()
	defer container.lock.RUnlock()
	return container.UsedMem
}

//设置内存使用大小
func (container *Container) SetUsedMem(usedMem int64) {
	container.lock.Lock()
	defer container.lock.Unlock()
	container.UsedMem = usedMem
}

//不用修改内存返回true
func (container *Container) UpUsedCount() bool {
	container.lock.Lock()
	defer container.lock.Unlock()
	b := container.UsedCount == 0
	container.UsedCount++
	return b
}

//不用修改内存返回true
func (container *Container) DownUsedCount() bool {
	container.lock.Lock()
	defer container.lock.Unlock()
	b := container.UsedCount == 1
	container.UsedCount--
	return b
}
