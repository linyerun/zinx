package znet

import "github.com/linyerun/zinx/ziface"

type Request struct {
	conn ziface.IConnection //已经和客户端建立好的链接
	msg  ziface.IMessage    //客户端请求的数据(把之前的[]byte换掉)
}

// NewRequest 比较简单，可以不要构造函数
func NewRequest(connection ziface.IConnection, msg ziface.IMessage) *Request {
	return &Request{
		conn: connection,
		msg:  msg,
	}
}

func (r *Request) GetConnection() ziface.IConnection {
	return r.conn
}

func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

func (r *Request) GetMsgId() uint32 {
	return r.msg.GetMsgId()
}

func (r *Request) GetMsgLen() uint32 {
	return r.msg.GetDataLen()
}
