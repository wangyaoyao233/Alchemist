package iface

//路由的抽象接口，数据是IRequest
type IRouter interface {
	//处理业务之前的方法
	PreHandle(IRequest)
	//处理业务的主方法
	Handle(IRequest)
	//处理业务之后的方法
	PostHandle(IRequest)
}
