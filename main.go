package main

import (
	"fmt"
	"log"
	"net"
	"voi/core"
)

func main() {

	server, err := core.NewVoIPServer("127.0.0.1:9091", 2)

	if err != nil {
		log.Fatal(err)
	}

	server.GetOrCreateRoom(1)

	p := core.Packet{
		Data: nil,
		Addr: &net.UDPAddr{
			IP: net.IPv4(127, 0, 0, 1),
		},
	}

	h := core.Header{
		MsgType:   2,
		Timestamp: 1,
		BodySize:  1,
	}

	m, err := core.NewPingMessage(&p, &h)

	if err != nil {
		fmt.Print(err)
	}

	err = m.Process(server)

	if err != nil {
		fmt.Print(err)
	}

	/*s := core.Session{
		Addr: &net.UDPAddr{
			IP:   net.IPv4(127, 0, 0, 1),
			Port: 0,
			Zone: "",
		},
		RoomId:   1,
		UserId:   1,
		PingChan: nil,
	}*/

	//server.JoinRoom(&s)

}
