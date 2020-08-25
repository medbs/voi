package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
	"voi/core"
)

func main() {

	ip, _, err := net.ParseCIDR("127.0.0.1/24")

	h := core.Header{
		MsgType: 2,
	}

	p := core.Packet{
		Data: make([]byte, 30, 30),
		Addr: &net.UDPAddr{
			IP:   ip,
			Port: 7777,
			Zone: "",
		},
	}

	//create ping message
	m, err := core.NewPingMessage(&p, &h)

	if err != nil {
		fmt.Println(err)
	}

	//m.SentTime = time.Now()

	s, err := net.ResolveUDPAddr("udp4", "127.0.0.1:9091")
	c, err := net.DialUDP("udp4", nil, s)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("The UDP server is %s\n", c.RemoteAddr().String())
	defer c.Close()

	var b bytes.Buffer

	err = gob.NewEncoder(&b).Encode(m)

	if err != nil {
		fmt.Println(err)
		return
	}

	//byteKey := []byte(m)
	for i := 0; i < 10; i++ {
		_, err = c.Write(b.Bytes())
	}
	if err != nil {
		fmt.Println(err)
		return
	}

}
