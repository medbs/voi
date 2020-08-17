package core

import (
	"errors"
)


type JoinMessage struct {
	*Packet
	*Header
	key []byte
}

func NewJoinMessage(p *Packet, h *Header) (*JoinMessage, error) {
	m := &JoinMessage{p, h, make([]byte, 0, 0)}
	err := m.Parse()
	return m, err
}

func (m *JoinMessage) Parse() error {
	if len(m.Data) < 11+int(m.BodySize) {
		return errors.New("body size is too small in join message")
	}
	m.key = m.Data[11 : 11+int(m.BodySize)]
	return nil
}

func (m *JoinMessage) Process(vs *VoIPServer) error {
	//uid, rid, err := vs.CheckJoinKey(m.key)
	/*if err != nil {
		return err
	}*/
	vs.JoinRoom(NewSession(1, 1, m.Addr, vs.GetConn()))
	return nil
}

