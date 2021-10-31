package main

import (
	"fmt"
	"net"
)

func main() {
	l, err := net.Listen("tcp", ":8888")
	if err != nil {
		fmt.Println("listen error:", err)
		return
	}

	for {
		con, err := l.Accept()
		if err != nil {
			fmt.Println("accept error:", err)
			break
		}
		go handleConn(con)
	}
}

func handleConn(con net.Conn) {
	defer func(con net.Conn) {
		_ = con.Close()
	}(con)
	for {
		buf:=make([]byte,1024)
		size, err := con.Read(buf)
		if err != nil {
			fmt.Printf("read data err: %s\n",err)
			return
		}
		_, err = con.Write(buf[0:size])
		if err != nil {
			fmt.Printf("write data err: %s\n",err)
			return
		}
	}
}
