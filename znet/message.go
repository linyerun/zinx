package znet

import "github.com/linyerun/zinx/ziface"

type Message struct {
	Id      uint32 //消息的ID
	DataLen uint32 //消息的长度
	Data    []byte //消息的内容
}

func NewMessage(id uint32, data []byte) ziface.IMessage {
	return &Message{
		Id:      id,
		DataLen: uint32(len(data)),
		Data:    data,
	}
}

func (m *Message) GetMsgId() uint32 {
	return m.Id
}

func (m *Message) SetMsgId(id uint32) {
	m.Id = id
}

func (m *Message) GetDataLen() uint32 {
	return m.DataLen
}

func (m *Message) SetDataLen(len uint32) {
	m.DataLen = len
}

func (m *Message) GetData() []byte {
	return m.Data
}

func (m *Message) SetData(data []byte) {
	m.Data = data
}
