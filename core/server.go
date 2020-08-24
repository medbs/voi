package core

import (
	"log"
	"net"
	"sync"
	"time"
)

type VoIPServer struct {
	conn         *net.UDPConn
	PacketChan   chan *Packet
	Wg           *sync.WaitGroup
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
		PacketChan:   make(chan *Packet, 100000),
		Wg:           new(sync.WaitGroup),
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
	log.Print("waiting all goroutines are stopped")
	vs.Wg.Wait()
	// finlize...
}

//read the packet from UDP
func (vs *VoIPServer) readLoop() {
	var buf [2048]byte
	vs.Wg.Add(1)
	defer func() {
		vs.Wg.Done()
		//vs.Shutdown()
	}()
	for {
		n, addr, err := vs.conn.ReadFromUDP(buf[0:])
		if err != nil {
			return
		}
		vs.PacketChan <- &Packet{buf[:n], addr, time.Now(),time.Time{}}
	}
}

//transform packet to message
func (vs *VoIPServer) analyzeLoop() {
	vs.Wg.Add(1)
	defer func() {
		vs.Wg.Done()
		//vs.Shutdown()
	}()
	for {
		select {
		case p := <-vs.PacketChan:

			if p == nil {
				continue
			}
			msg, err := p.ToMessage()
			if err != nil {
				log.Print(err.Error())
				break
			}

			err = msg.Process(vs)
			if err != nil {
				log.Print(err.Error())
			}
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
	//todo fix the addrStr random port changing
	//return nil
	return vs.sessions["127.0.0.1:7777"]
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
