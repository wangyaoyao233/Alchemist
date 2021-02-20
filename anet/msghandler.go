package anet

import (
	"Alchemist/iface"
	"Alchemist/utils"
	"fmt"
	"strconv"
)

type MsgHandle struct {
	Apis map[uint32]iface.IRouter
	//Worker取任务的消息队列
	TaskQueue []chan iface.IRequest
	//Worker工作池的数量
	WorkerPoolSize uint32
}

//创建MsgHandler方法
func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:           make(map[uint32]iface.IRouter),
		TaskQueue:      make([]chan iface.IRequest, int(utils.GlobalObject.WorkerPoolSize)),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
	}
}

//执行对应的router消息处理方法
func (mh *MsgHandle) DoMsgHandler(request iface.IRequest) {
	//1.从request中找到msgID
	handle, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgID: ", request.GetMsgID(), "is not found")
	}
	//2.根据msgID调度对应的router
	handle.PreHandle(request)
	handle.Handle(request)
	handle.PostHandle(request)
}

//添加router路由
func (mh *MsgHandle) AddRouter(msgID uint32, router iface.IRouter) {
	//1.判断当前绑定的api处理方法是否存在
	if _, ok := mh.Apis[msgID]; ok {
		panic("repeate api, msgID " + strconv.Itoa(int(msgID)))
	}
	//2.不存在则添加
	mh.Apis[msgID] = router
	fmt.Println("Add api MsgID: ", msgID, "succ")
}

//启动一个Worker工作池(开启工作池的动作只能发生一次， 一个框架只有一个Worker工作池)
func (mh *MsgHandle) StartWorkerPool() {
	//根据WorkerPoolSize分别开启Worker，每个Worker用一个go
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		//1.给当前的worker对应的消息队列开辟空间
		mh.TaskQueue[i] = make(chan iface.IRequest, int(utils.GlobalObject.MaxWorkerTaskLen))
		//2.启动当前的Worker,阻塞等待消息从channel传递进来
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}

//启动一个Worker工作流程
func (mh *MsgHandle) StartOneWorker(workerID int, taskQueue chan iface.IRequest) {
	fmt.Println("Worker ID=", workerID, "is started..")

	//阻塞等待对应消息队列的消息
	for {
		select {
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
		}
	}
}

func (mh *MsgHandle) SendMsgToTaskQueue(request iface.IRequest) {
	//1.平均分配给不同Worker
	id := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	fmt.Println("Add connID=", request.GetConnection().GetConnID(), "request MsgID:", request.GetMsgID(), "to WorkerID:", id)

	//2.将消息发送给对应的worker的TaskQueue
	mh.TaskQueue[id] <- request
}
