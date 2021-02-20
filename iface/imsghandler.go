package iface

type IMsgHandle interface {
	//执行对应的router消息处理方法
	DoMsgHandler(request IRequest)
	//添加router路由
	AddRouter(msgID uint32, router IRouter)
	//启动Worker工作池
	StartWorkerPool()
	//发送消息给消息任务队列
	SendMsgToTaskQueue(request IRequest)
}
