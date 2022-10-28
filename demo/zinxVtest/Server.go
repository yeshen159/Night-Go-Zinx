package main

import (
	"fmt"

	"github.com/aceld/zinx/znet"
)

func main() {
	fmt.Println("vim-go")
	s := znet.NewServer("dsa")

	s.Server()
}
