package core

import (
	"bytes"
	"errors"
	"fmt"
	"time"
)

type PingMessage struct {
	*Packet
	*Header
	payload []byte
}


func NewPingMessage(p *Packet, h *Header) (*PingMessage, error) {
	m := &PingMessage{p, h, make([]byte, 0, 0)}
	err := m.Parse()
	return m, err
}

func (m *PingMessage) Parse() error {
	if len(m.Data) < 11+int(m.BodySize) {
		return errors.New("body size is too small in ping message")
	}
	m.payload = m.Data[11 : 11+int(m.BodySize)]
	return nil
}

func (m *PingMessage) ToPong() []byte {
	var b bytes.Buffer
	return b.Bytes()
}

func (m *PingMessage) Process(vs *VoIPServer) error {
	s := vs.GetSession(m.Addr.String())
	if s == nil {
		return errors.New("not authed")
	}
	m.ReceivedTime = time.Now()
	s.PingChan <- m
	fmt.Print("processed")
	// send pong
	// room := s.GetRoom()
	return nil
}
