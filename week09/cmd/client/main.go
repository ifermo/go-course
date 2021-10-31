package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	address := "127.0.0.1:8888"
	con, err := net.Dial("tcp", address)
	if err != nil {
		log.Printf("%d: dial error: %s", address, err)
		return
	}
	defer func(con net.Conn) {
		_ = con.Close()
	}(con)
	data := "Hello Network!"
	_, err = con.Write([]byte(data))
	if err != nil {
		log.Printf("write data err: %s", err)
	}

	buf := make([]byte, 1024)
	size, err := con.Read(buf)
	if err != nil {
		fmt.Printf("read data err: %s\n", err)
		return
	}
	fmt.Printf("receive data: %s\n", string(buf[0:size]))
}
