package ziface

// IMessage 消息的抽象层
type IMessage interface {
	GetMsgId() uint32      //获取消息ID
	SetMsgId(id uint32)    //设置消息ID
	GetDataLen() uint32    //获取消息长度
	SetDataLen(len uint32) //设置消息长度
	GetData() []byte       //获取消息内容
	SetData(data []byte)   //设置消息内容
}
