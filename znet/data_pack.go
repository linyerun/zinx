package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	. "github.com/linyerun/zinx/global_properties"
	"github.com/linyerun/zinx/ziface"
	"io"
	"net"
	"sync"
)

type DataPack struct {
}

//单例模式创建
var dp ziface.IDataPack
var once sync.Once

func NewDataPack() ziface.IDataPack {
	once.Do(func() {
		dp = new(DataPack)
	})
	return dp
}

func (d *DataPack) GetHeadLen() uint32 {
	return 8
}

func (d *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {

	//判断一下，以防万一
	if msg.GetDataLen() != uint32(len(msg.GetData())) {
		return nil, errors.New(fmt.Sprintf("msgId=%v, msgLen != uint32(len(msg.GetData))", msg.GetMsgId()))
	}

	//创建一个存放byte字节的缓冲
	dataBuff := bytes.NewBuffer(make([]byte, 0))

	//将dataLen写进dataBuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetDataLen()); err != nil {
		return nil, err
	}

	//将MsgId写进dataBuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}

	//将MsgData写进dataBuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

func (d *DataPack) Unpack(conn net.Conn) (ziface.IMessage, error) {

	msg := new(Message)

	//获取msgLen
	if err := binary.Read(conn, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	//数据长度过长(我绝对是总体不可以超过)
	if msg.DataLen > GlobalObject.MaxPacketSize {
		return nil, errors.New("the packet length exceeded the user-defined packet length")
	}

	//获取msgId
	if err := binary.Read(conn, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	//获取msgData
	msg.Data = make([]byte, msg.DataLen)
	n, err := io.ReadFull(conn, msg.Data)
	if err != nil {
		return nil, err
	}
	if uint32(n) != msg.DataLen {
		return nil, errors.New("body len not equals msgLen")
	}

	return msg, nil
}
