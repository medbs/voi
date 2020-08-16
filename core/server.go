package core

import (
	"log"
	"net"
	"sync"
)

type VoIPServer struct {
	conn         *net.UDPConn
	packetChan   chan *Packet
	wg           *sync.WaitGroup
	shutdownChan chan struct{}
	sessions     map[string]*Session
	sessionM     sync.RWMutex
	rooms        map[int]*Room
	roomM        sync.RWMutex
}

func NewVoIPServer(addr string, numLoop int) (*VoIPServer, error) {
	udpAddr, err := net.ResolveUDPAddr("udp4", addr)
	if err != nil {
		return nil, err
	}
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return nil, err
	}
	vs := VoIPServer{
		conn:         conn,
		packetChan:   make(chan *Packet, 100000),
		wg:           new(sync.WaitGroup),
		shutdownChan: make(chan struct{}),
		sessions:     make(map[string]*Session),
		sessionM:     sync.RWMutex{},
		rooms:        make(map[int]*Room),
		roomM:        sync.RWMutex{},
	}
	go vs.readLoop()
	for i := 0; i < numLoop; i++ {
		go vs.analyzeLoop()
	}
	return &vs, nil
}

// voip
func (vs *VoIPServer) Shutdown() {
	close(vs.shutdownChan) //goroutine
	vs.conn.Close()        // socket
	log.Print("waiting all gorutine are stoped")
	vs.wg.Wait()
	// finlize...
}

//UDP
func (vs *VoIPServer) readLoop() {
	var buf [2048]byte
	vs.wg.Add(1)
	defer func() {
		vs.wg.Done()
		vs.Shutdown()
	}()
	for {
		n, addr, err := vs.conn.ReadFromUDP(buf[0:])
		if err != nil {
			// socketが死んでるので落とす -> timeoutかEOS
			return
		}
		vs.packetChan <- &Packet{buf[:n], addr}
	}
}

func (vs *VoIPServer) analyzeLoop() {
	vs.wg.Add(1)
	defer func() {
		vs.wg.Done()
		vs.Shutdown()
	}()
	for {
		select {
		case p := <-vs.packetChan:

			if p == nil {
				continue
			}
			msg, err := p.ToMessage()
			if err == nil {
				continue
			}
			msg.Process(vs)
		case _, ok := <-vs.shutdownChan:
			if !ok {
				log.Print("shutdown analyze loop")
				return
			}
		}
	}
}


func (vs *VoIPServer) GetConn() *net.UDPConn {
	return vs.conn
}

func (vs *VoIPServer) GetSession(addrStr string) *Session {
	vs.sessionM.RLock()
	defer vs.sessionM.RUnlock()
	if session, ok := vs.sessions[addrStr]; ok {
		return session
	}
	return nil
}

func (vs *VoIPServer) GetRoom(roomId int) *Room {
	vs.roomM.RLock()
	defer vs.roomM.RUnlock()
	if room, ok := vs.rooms[roomId]; ok {
		return room
	}
	return nil
}

func (vs *VoIPServer) GetOrCreateRoom(roomId int) (*Room, bool) {
	vs.roomM.Lock()
	defer vs.roomM.Unlock()
	if room, ok := vs.rooms[roomId]; ok {
		return room, false
	} else {
		room = NewRoom(vs, roomId)
		vs.rooms[roomId] = room
		log.Print("create %v", room)
		return room, true
	}
}

func (vs *VoIPServer) JoinRoom(s *Session) {
	addrStr := s.Addr.String()

	if ses := vs.GetSession(addrStr); ses != nil {
		log.Print("session already exits")
		return
	}
	room, _ := vs.GetOrCreateRoom(s.RoomId)
	// join to room
	err := room.JoinRoom(s)
	if err == nil {
		// add to session
		vs.sessionM.Lock()
		vs.sessions[addrStr] = s
		vs.sessionM.Unlock()
	} else {
		log.Print(err)
	}
}
