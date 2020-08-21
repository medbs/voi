package main

import (
	"bufio"
	"fmt"
	"net"
	"voi/core"
)

func main() {

	h := core.Header{
		MsgType:   2,
		Timestamp: 1,
		BodySize:  1,
	}

	p := core.Packet{
		Data: make([]byte, 30,30),
		Addr: &net.UDPAddr{
			IP: net.IPv4(127, 0, 0, 1),
		},
	}

	//create ping message
	m, err := core.NewPingMessage(&p, &h)


	//pa :=  make([]byte, 1024)

	conn, err := net.Dial("udp", "127.0.0.1:9091")
	if err != nil {
		fmt.Printf("Some error %v", err)
		return
	}
	fmt.Fprintf(conn, "Hi UDP Server, How are you doing?")
	_, err = bufio.NewReader(conn).Read(m.Data)
	if err == nil {
		fmt.Printf("%s\n", p)
	} else {
		fmt.Printf("Some error %v\n", err)
	}
	conn.Close()

}
