package core

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"net"
	"time"
)

const (
	MSG_TYPE_JOIN = uint8(1)
	MSG_TYPE_PING = uint8(2)
)

type Packet struct {
	Data []byte
	Addr *net.UDPAddr
	SentTime time.Time
	ReceivedTime time.Time
}

type Header struct {
	MsgType   uint8
	Timestamp uint64
	BodySize  uint16
}

func readHeader(packet []byte) (*Header, error) {
	if len(packet) < 11 {
		// lower than defined header size
		return nil, errors.New("invalid packet size")
	}

	pm := PingMessage{}
	err := gob.NewDecoder(bytes.NewReader(packet)).Decode(&pm)

	if err != nil {
		fmt.Print(err)
	}

	return &Header{pm.MsgType, pm.Timestamp, pm.BodySize}, nil
}

func (p *Packet) ToMessage() (Message, error) {
	header, err := readHeader(p.Data)
	if err != nil {
		return nil, err
	}
	switch header.MsgType {
	case MSG_TYPE_JOIN:
		return NewJoinMessage(p, header)
	case MSG_TYPE_PING:
		return NewPingMessage(p, header)
	default:
		return nil, errors.New("packet is not a voip message")
	}
}
