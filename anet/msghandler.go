package anet

import (
	"Alchemist/iface"
	"fmt"
	"strconv"
)

type MsgHandle struct {
	Apis map[uint32]iface.IRouter
}

//创建MsgHandler方法
func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis: make(map[uint32]iface.IRouter),
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
