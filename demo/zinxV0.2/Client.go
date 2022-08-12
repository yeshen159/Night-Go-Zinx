package main

import (
	"fmt"
	"net"
	"time"
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
		//2 链接调用Write 写数据
		_, err := conn.Write([]byte("Hello Zinx V0.2.."))

		if err != nil {
			fmt.Println("write conn  err, exit")
			return
		}

		buf := make([]byte, 512)

		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read buf  err, exit")
			return
		}

		fmt.Printf("server call back : %s, cnt = %d\n", buf, cnt)

		//cup阻塞
		time.Sleep(1 * time.Second)
	}
}
