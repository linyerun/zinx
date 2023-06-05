package znet

import (
	"fmt"
	"github.com/linyerun/zinx/global_properties"
	"github.com/linyerun/zinx/ziface"
	"net"
)

// Server 实现IServer接口
type Server struct {
	Name                string                              //服务器名称
	IPVersion           string                              //服务器绑定的ip版本
	IP                  string                              //服务器监听的IP
	Port                int                                 //服务器监听的端口
	ReqProcessingCenter ziface.IReqProcessingCenter         //加入、运行并管理多个handler
	ConnsMgr            ziface.IConnsManager                //集成连接管理模块
	OnConnStart         func(connection ziface.IConnection) //创建连接之后立马调用的hook
	OnConnStop          func(connection ziface.IConnection) //停止连接之前那瞬间调用hook
}

func NewServer() (server ziface.IServer) {
	// 初始化全局变量
	global_properties.Init()

	// 使用全局变量
	server = &Server{
		Name:                global_properties.GlobalObject.Name,
		IPVersion:           "tcp4",
		IP:                  global_properties.GlobalObject.Host,
		Port:                global_properties.GlobalObject.TcpPort,
		ReqProcessingCenter: newReqProcessingCenter(),
		ConnsMgr:            newConnsManager(),
	}
	return
}

func (s *Server) Serve() {
	s.start()
}

func (s *Server) AddHandler(msgId uint32, handler ziface.IHandler) {
	s.ReqProcessingCenter.AddHandler(msgId, handler)
}

func (s *Server) AddOnConnStart(f func(connection ziface.IConnection)) {
	s.OnConnStart = f
}

func (s *Server) AddOnConnStop(f func(connection ziface.IConnection)) {
	s.OnConnStop = f
}

func (s *Server) callOnConnStart(connection ziface.IConnection) {
	if connection == nil {
		fmt.Println("connection is nil, can not call the func of OnConnStart!")
		return
	}
	if s.OnConnStart == nil {
		fmt.Println("未注册OnConnStart,无需调用！")
		return
	}
	s.OnConnStart(connection)
}

func (s *Server) callOnConnStop(connection ziface.IConnection) {
	if connection == nil {
		fmt.Println("connection is nil, can not call the func of OnConnStop!")
		return
	}
	if s.OnConnStop == nil {
		fmt.Println("未注册OnConnStop,无需调用！")
		return
	}
	s.OnConnStop(connection)
}

func (s *Server) start() {
	fmt.Printf("[Zinx] Server Name: %s, Listening at IP: %s, Port: %d is starting!\n", global_properties.GlobalObject.Name, global_properties.GlobalObject.Host, global_properties.GlobalObject.TcpPort)
	fmt.Printf("[Zinx] Version %s, MaxConn: %d, MaxPacketSize: %d\n", global_properties.GlobalObject.Version, global_properties.GlobalObject.MaxConn, global_properties.GlobalObject.MaxPacketSize)
	fmt.Printf("[start] Server Listener at IP: %s, Port: %d, is starting\n", s.IP, s.Port)

	//开启工作池
	s.ReqProcessingCenter.StartWorkers() //开启工作池

	//1. 获取一个TCP的Addr
	addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
	if err != nil {
		fmt.Println("resolve tcp addr error: ", err)
		return
	}

	//2. 监听服务器地址
	listener, err := net.ListenTCP(s.IPVersion, addr)
	if err != nil {
		fmt.Println("listen", s.IPVersion, "err", err)
		return
	}
	fmt.Println("start Zinx server succeed,", s.Name, "succeed, Listening...")

	//3. 阻塞的等待客户端连接，处理客户端连接业务(读写)
	var connId uint32 = 0
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Println("Accept err", err)
			continue
		}

		//超出最大连接数之后就拒绝连接
		if s.ConnsMgr.Len() == global_properties.GlobalObject.MaxConn {
			if err := conn.Close(); err != nil { // 关闭连接
				fmt.Println("关闭连接失败, err:", err)
				continue
			}
			fmt.Println("to many connections,over MaxConnections which is equal", global_properties.GlobalObject.MaxConn)
			continue
		}

		connection := newConnection(s, conn, connId)
		connId++ //这一步不要忘记了
		connection.start()
	}
}

func (s *Server) stop() {
	fmt.Println("[STOP] Zinx server name =", s.Name)
	s.ConnsMgr.ClearConns()
}
