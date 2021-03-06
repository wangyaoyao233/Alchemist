package utils

import (
	"Alchemist/iface"
	"encoding/json"
	"io/ioutil"
)

//存储框架的全局参数，提供其他模块使用
//一些参数是通过conf.json由用户进行配置
type GlobalObj struct {
	//全局的Server对象
	TcpServer iface.IServer
	//当前服务器监听的IP
	Host string
	//当前服务器监听的端口号
	TcpPort int
	//当前服务器名称
	Name string

	//框架
	Version          string
	MaxConn          int
	MaxPackageSize   uint32
	WorkerPoolSize   uint32 //Worker池的Goroutine数量
	MaxWorkerTaskLen uint32 //框架允许用户最多开辟多少个Worker
	MaxMsgChanLen    uint32 //SendBuffMsg缓冲的最大长度
}

//定义一个全局的对外GlobalObj
var GlobalObject *GlobalObj

//提供一个init方法，初始化全局GlobalObject
func init() {
	//如果配置文件没有加载，默认的值
	GlobalObject = &GlobalObj{
		Name:             "ServerApp",
		Version:          "v0.0",
		TcpPort:          8999,
		Host:             "0.0.0.0",
		MaxConn:          1000,
		MaxPackageSize:   4096,
		WorkerPoolSize:   10,
		MaxWorkerTaskLen: 1024,
		MaxMsgChanLen:    1024,
	}

	//从conf.json加载用户自定义值
	GlobalObject.Reload()
}

func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("conf/setting.json")
	if err != nil {
		panic(err)
	}
	//将json文件解析到struct中
	err = json.Unmarshal(data, g)
	if err != nil {
		panic(err)
	}
}
