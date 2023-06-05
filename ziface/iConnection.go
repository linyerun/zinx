package ziface

import "net"

// IConnection 连接模块的抽象层。接口的方法是提供给开发者使用的
type IConnection interface {
	GetTCPConnection() *net.TCPConn                  //获取当前连接的绑定socket.conn
	GetConnID() uint32                               //获取当前连接模块的链接ID
	RemoteAddr() net.Addr                            //获取远程连接客户的TCP状态、IP和Port
	SendMsg(msgId uint32, msgData []byte) error      //发送数据，将数据发送给远程客户端
	SendMsgToAll(msgId uint32, msgData []byte)       //群发数据给当前在线的连接
	SetProperty(key string, value interface{}) error //设置连接属性
	GetProperty(key string) (interface{}, error)     //获取连接属性
	RemoveProperty(key string) error                 //删除一个连接属性
}
