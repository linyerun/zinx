package znet

import (
	"errors"
	"fmt"
	"github.com/linyerun/zinx/global_properties"
	"net"
	"sync"
)

type Connection struct {
	server     *Server                //方便我们使用Conn的两个钩子函数
	Conn       *net.TCPConn           //当前链接的socket TCP套接字
	ConnID     uint32                 //链接ID
	isClose    bool                   //判断连接是否关闭
	exitChan   chan struct{}          //告知当前链接已经退出/停止的 channel (无缓冲的其实也行，好像还是有缓存的好)
	msgChan    chan []byte            //消息队列，sendMsg先发到此处等待Writer处理
	properties map[string]interface{} //连接属性集合
	propLock   sync.RWMutex           //保护连接属性的读写锁
}

func newConnection(server *Server, conn *net.TCPConn, connID uint32) *Connection {
	c := &Connection{
		server:     server,
		Conn:       conn,
		ConnID:     connID,
		msgChan:    make(chan []byte, global_properties.GlobalObject.MaxMsgChanSize),
		exitChan:   make(chan struct{}),
		properties: map[string]interface{}{},
	}
	//每当Connection被创建的时候，都应该注册到connsMgr当中
	c.server.ConnsMgr.Add(c)
	return c
}

func (c *Connection) startReader() {
	fmt.Println("Reader Goroutine is running...")

	defer c.stop() //读协程停止的时候，告知连接停止
	defer fmt.Println("ConnID =", c.ConnID, "[Reader] is exit, remote addr is", c.RemoteAddr().String())

	dataPack := NewDataPack()

	for {
		//有个疑惑，没有阻塞，读取完用户又不操作不就直接没了。只要你不关闭，读到末尾了就阻塞
		msg, err := dataPack.Unpack(c.Conn)
		if _, ok := err.(*net.OpError); ok {
			fmt.Println("Remote host take the initiative to quit!")
			c.stop() //关闭连接，回收资源
			return
		}
		if err != nil {
			fmt.Println("receive buf err:", err)
			//出现错误应该可以救一下，没必要直接关闭连接，可以重新尝试一下？
			return
		}
		//把req交给任务线程池处理
		req := NewRequest(c, msg)
		c.server.ReqProcessingCenter.SendReqToWorkers(req)
	}
}

func (c *Connection) startWriter() {
	fmt.Println("Writer Goroutine is running...")
	defer fmt.Println("ConnID =", c.ConnID, "[Writer] is exit, remote addr is", c.RemoteAddr().String())
	for {
		select {
		case data := <-c.msgChan:
			if _, err := c.Conn.Write(data); err != nil {
				//正常连接是不会出问题的，除非远程客户端真的挂了
				fmt.Println("write data err", err)
				return
			}
		case <-c.exitChan:
			//关闭写协程信号
			return
		}
	}
}

func (c *Connection) start() {
	fmt.Println("Conn start()... ConnID =", c.ConnID)
	go c.startReader()          //读协程
	go c.startWriter()          //写协程
	c.server.callOnConnStart(c) //调用连接创建的hook函数
}

func (c *Connection) stop() { //这个方法只需要给Reader协程调用就好了，Writer协程就不需要管它了，极端情况下怎么搞以后再说
	fmt.Println("Conn stop()... ConnID =", c.ConnID)

	//读协程已经关闭
	//关闭写协程
	c.exitChan <- struct{}{} //这里做成无缓冲的吧,这样可以确保写协程是关闭的

	//在连接关闭之前调用对于hook函数
	//服务端就不能对客户端进行读写操作了
	c.server.callOnConnStop(c)

	//关闭socket连接
	if err := c.Conn.Close(); err != nil {
		panic(err)
	}

	//删除在connMgr里面的记录
	c.server.ConnsMgr.Remove(c.ConnID)

	//回收资源
	close(c.exitChan)
	close(c.msgChan)

	//改变isClose值
	c.isClose = true
}

func (c *Connection) SendMsg(msgId uint32, msgData []byte) error {
	if c.isClose {
		return errors.New("connection has closed")
	}
	dataPack := NewDataPack()
	bytes, err := dataPack.Pack(NewMessage(msgId, msgData))
	if err != nil {
		return err
	}
	c.msgChan <- bytes //使用这个msg管道发送到写协程进行处理
	return nil
}

func (c *Connection) SendMsgToAll(msgId uint32, msgData []byte) {
	connsManager := c.server.ConnsMgr.(*ConnsManager)
	connsManager.connLock.RLock()
	defer connsManager.connLock.RUnlock()
	for _, connection := range connsManager.conns {
		if err := connection.SendMsg(msgId, msgData); err != nil {
			fmt.Printf("message group sending: ConnId = %d has err = %v\n", connection.GetConnID(), err)
			continue
		}
	}
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection) SetProperty(key string, value interface{}) error {
	c.propLock.RLock()
	_, ok := c.properties[key]
	c.propLock.RUnlock()
	if ok {
		return errors.New(fmt.Sprintf("%s is already in the map", key))
	}
	c.propLock.Lock()
	c.properties[key] = value
	c.propLock.Unlock()
	return nil
}

func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propLock.RLock()
	val, ok := c.properties[key]
	c.propLock.RUnlock()
	if !ok {
		return nil, errors.New("can not found in map")
	}
	return val, nil
}

func (c *Connection) RemoveProperty(key string) error {
	c.propLock.RLock()
	_, ok := c.properties[key]
	c.propLock.RUnlock()
	if !ok {
		return errors.New("can not found in map")
	}
	c.propLock.Lock()
	delete(c.properties, key)
	c.propLock.Unlock()
	return nil
}
