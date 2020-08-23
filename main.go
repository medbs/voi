package main

import (
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

	//create session and join room
	s := core.Session{
		Addr: &net.UDPAddr{
			IP:   net.IPv4(127, 0, 0, 1),
			Port: 7777,
			Zone: "",
		},
		RoomId:   1,
		UserId:   1,
		PingChan: make(chan *core.PingMessage,100000),
	}
	server.JoinRoom(&s)

	server.Wg.Wait()

}
