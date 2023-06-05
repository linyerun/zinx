package znet

import (
	"fmt"
	"net"
	"testing"
	"time"
)

//测试DataPack这三个方法
func TestDataPack(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		t.Error(err)
		return
	}
	//创建一个goroutine负责处理业务(服务端)
	go func() {
		conn, err := listener.Accept() //看源码可以知道, conn每次获取连接之后会new一个新的,是指针,我们直接引用传递也没关系
		if err != nil {
			t.Error(err)
			return
		}
		//为每个客户端分配一个goroutine, conn不需要传，直接根据闭包特性引用传递也不会出错
		go func() {
			dp := NewDataPack()
			for {
				//可以连续读conn数据
				msg, err := dp.Unpack(conn)
				if err != nil {
					fmt.Println("unpack msg error: ", err)
					return //反正是测试，直接让它停了
				}
				fmt.Println("receive msg==>", "msgLen:", msg.GetDataLen(), "msgId:", msg.GetMsgId(), "msgData:", string(msg.GetData()))
			}
		}()
	}()

	//模拟客户端
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		t.Error(err)
		return
	}
	dataPack := NewDataPack()
	for i := 0; i < 10; i++ {
		bs := []byte(fmt.Sprintf("Hello,World!I'm %d", i))
		b1, err := dataPack.Pack(NewMessage(uint32(i), bs))
		b2, err := dataPack.Pack(NewMessage(uint32(i), bs))
		msgBytes := append(b1, b2...)
		if err != nil {
			fmt.Println(i, ":", err)
			continue
		}
		_, err = conn.Write(msgBytes)
		if err != nil {
			fmt.Println(i, ",", err)
			continue
		}
		time.Sleep(time.Second * 3)
	}
	println("客户端退出！！！")
}
