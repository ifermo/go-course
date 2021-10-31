package main

import (
	"fmt"
	"net"
	. "week09/internal"
	"week09/pkg/decoder"
	"week09/pkg/encoder"
)

func main() {
	l, err := net.Listen("tcp", ":8888")
	if err != nil {
		fmt.Println("listen error:", err)
		return
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("accept error:", err)
			break
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer func(c net.Conn) {
		_ = c.Close()
	}(conn)
	reader := decoder.NewLengthFieldBaseFrameDecoder(conn, decoder.BigEndian, 1e8, 4)
	writer := encoder.NewLengthFieldPrepender(conn, decoder.BigEndian, 4)
	for {
		buf, err := reader.Read()
		if err != nil {
			fmt.Printf("read data err: %s\n", err)
			return
		}
		msg := Message{}
		err = msg.Read(buf)
		if err != nil {
			fmt.Printf("parsing data err: %s\n", err)
			return
		}
		err = writer.Write(msg.Write())
		if err != nil {
			fmt.Printf("write data err: %s\n", err)
			return
		}
	}
}
