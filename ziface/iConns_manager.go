package ziface

// IConnsManager 连接管理模块的抽象层
type IConnsManager interface {
	Add(conn IConnection)                   //添加连接
	Remove(connId uint32)                   //删除连接
	Get(connId uint32) (IConnection, error) //获取连接
	Len() int                               //获取连接总数
	ClearConns()                            //关闭所有连接
}
