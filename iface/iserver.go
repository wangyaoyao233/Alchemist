package iface

//定义一个服务器接口
type IServer interface {
	//启动服务器
	Start()
	//停止服务器
	Stop()
	//运行服务器
	Serve()

	//添加路由: 给当前的server注册一个路由方法，供客户端的连接使用
	AddRouter(msgID uint32, router IRouter)

	//获取ConnMgr连接管理器方法
	GetConnMgr() IConnManager

	//注册连接之间的hook函数
	SetOnConnStart(func(conn IConnection))
	//注册连接销毁之前的hook函数
	SetOnConnStop(func(conn IConnection))
	//调用hook函数
	CallOnConnStart(conn IConnection)
	CallOnConnStop(conn IConnection)
}
