package ziface

// IReqProcessingCenter 请求处理模块的抽象层
type IReqProcessingCenter interface {
	ExecReqAppointedHandler(req IRequest)      //使用开发者自定义处理器处理请求
	SendReqToWorkers(req IRequest)             //把请求交给任务协程池处理
	AddHandler(msgId uint32, handler IHandler) //添加处理器
	StartWorkers()                             //开启协程池
	GetWorkPoolSize() uint32                   //获取任务协程池的协程数
}
