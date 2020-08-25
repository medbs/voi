package main

import (
	"fmt"
	"log"
	"net"
	"voi/core"
)

func main() {

	//xs := make([]float64, 10)

	pms := make([]core.PingMessage, 0, 10)
	//var min, max, avg float64

	msgChannel := make(chan core.PingMessage)
	server, err := core.NewVoIPServer("127.0.0.1:9091", 10, msgChannel)

	if err != nil {
		log.Fatal(err)
	}

	server.GetOrCreateRoom(1)

	ip, _, err := net.ParseCIDR("127.0.0.1/24")
	//create session and join room
	s := core.Session{
		Addr: &net.UDPAddr{
			IP:   ip,
			Port: 7777,
			Zone: "",
		},
		RoomId:   1,
		UserId:   1,
		PingChan: make(chan *core.PingMessage, 100000),
	}
	server.JoinRoom(&s)

	/*ses := server.GetSession("127.O.0.1:7777")
	pm := <-ses.PingChan*/

	ses := server.GetSession("127.O.0.1:7777")

	/*for {
		select {
		case pm := <-ses.PingChan:
			pms = append(pms, *pm)

	}*/

	for {
		pm := <-ses.PingChan
		if pm == nil {
			break
		}
		pms = append(pms, *pm)

		if len(pms) > 9 {
			fmt.Print("keksimus", len(pms))
		}

	}

	//fmt.Print("keksimus", len(pms))

	server.Wg.Wait()

}
