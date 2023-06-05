package ziface

// IRequest 请求的抽象层
type IRequest interface {
	GetConnection() IConnection //获取请求对应的连接
	GetData() []byte            //获取请求信息的内容
	GetMsgId() uint32           //获取请求信息的ID
	GetMsgLen() uint32          //获取请求信息的长度
}
