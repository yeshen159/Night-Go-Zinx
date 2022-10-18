package main

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/aceld/zinx/znet"
)

func main() {
	fmt.Println("client star...")

	time.Sleep(1 * time.Second)

	//1 直接链接远程
	conn, err := net.Dial("tcp", "127.0.0.1:8999")

	if err != nil {
		fmt.Println("client start err, exit")
	}

	for {

		//发送封包的message消息  MsgID:0

		dp := znet.NewDataPack()
		binaryMsg, err := dp.Pack(znet.NewMsgPackage(1, []byte(" client1 Test Message")))
		if err != nil {
			fmt.Println("Pack error:", err)
			return
		}

		if _, err := conn.Write(binaryMsg); err != nil {
			fmt.Println("write error:", err)
			return
		}

		//服务器应该给回复一个message数据, MsgID:1 pingping

		//1.先读取流中的head部分 得到ID 和 dataLen
		binaryHead := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(conn, binaryHead); err != nil {
			fmt.Println("read head error", err)
			break
		}

		//将二进制的head拆包到msg 结构体中
		msgHead, err1 := dp.Unpack(binaryHead)
		//_, err1 := dp.Unpack(binaryHead)
		if err1 != nil {
			fmt.Println("client unpack msgHead error", err1)
			break
		}

		if msgHead.GetDataLen() > 0 {
			//2.再根据dataLen进行二次读取，将data读出来
			msg := msgHead.(*znet.Message)
			msg.Data = make([]byte, msg.GetDataLen())

			if _, err := io.ReadFull(conn, msg.Data); err != nil {
				fmt.Println("read msg data error,", err)
				return
			}

			fmt.Println("--->Recv Server Msg: ID =", msg.ID, ", len =", msg.DataLen, ", data = ", string(msg.Data))

		}

		//cup阻塞
		time.Sleep(1 * time.Second)
	}
}
