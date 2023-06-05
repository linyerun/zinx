package ziface

import "net"

// IDataPack 封包拆包类的抽象层
type IDataPack interface {
	GetHeadLen() uint32                     //获取信息头长度
	Pack(msg IMessage) ([]byte, error)      //封包方法
	Unpack(conn net.Conn) (IMessage, error) //拆包方法
}
