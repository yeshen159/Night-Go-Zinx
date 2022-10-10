package zserver

import (
	"Night/zinx/utils"
	"Night/zinx/ziserver"
	"bytes"
	"encoding/binary"
	"errors"
)

/*
	封包，拆包 具体模块
	直接面向TCP连接中的数据流，用于处理TCP粘包问题
*/

type DataPack struct{}

//拆包封包实例的一个初始化方法

func NewDataPack() *DataPack {

	return &DataPack{}
}

//获取包的头的长度方法
func (dp *DataPack) GetHeadLen() uint32 {
	//DataLen uint32(4字节) + ID uint32(4字节)
	return 8
}

//封包方法
// |datalen|msgID|data|
func (dp *DataPack) Pack(msg ziserver.IMessage) ([]byte, error) {
	//创建一个存放bytes字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})

	//将dataLen 写进databuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgLen()); err != nil {
		return nil, err
	}

	//将MsgId 写进databuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}

	//将data数据 写进databuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

//拆包方法(将包的Head信息读出来) 之后再根据head信息里的data的长度，再进行一次读
func (dp *DataPack) Unpack(binaryData []byte) (ziserver.IMessage, error) {
	//创建一个从输入二进制数据的ioReader
	dataBuff := bytes.NewReader(binaryData)

	//只解压head信息，得到datalen和MsgID
	msg := &Message{}

	//读datalen
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}
	//读MsgID
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	//判断datalen是否已经超出了允许的最大包长
	if utils.GlobalObject.MaxPackageSize > 0 && msg.DataLen > utils.GlobalObject.MaxPackageSize {
		return nil, errors.New("too large msg data recv!")
	}

	return msg, nil
}
