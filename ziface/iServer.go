package ziface

// IServer 服务的抽象层。接口的方法是提供给开发者可以使用的
type IServer interface {
	Serve()                                        //运行服务器
	AddHandler(msgId uint32, handler IHandler)     //添加处理器
	AddOnConnStart(f func(connection IConnection)) //添加Conn的start hook
	AddOnConnStop(f func(connection IConnection))  //添加Conn的stop hook
}
