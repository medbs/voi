package main

import (
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
		Data: make([]byte, 30, 30),
		Addr: &net.UDPAddr{
			IP: net.IPv4(127, 0, 0, 1),
		},
	}

	//create ping message
	m, err := core.NewPingMessage(&p, &h)

	fmt.Print(m)

	s, err := net.ResolveUDPAddr("udp4", "127.0.0.1:9091")
	c, err := net.DialUDP("udp4", nil, s)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("The UDP server is %s\n", c.RemoteAddr().String())
	defer c.Close()

	byteKey := []byte(fmt.Sprintf("%v", m))
	_, err = c.Write(byteKey)

	if err != nil {
		fmt.Println(err)
		return
	}

}
