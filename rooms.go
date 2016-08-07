package phpim

import (
	"sync"
)

type Global struct {
	l     sync.RWMutex
	conns map[string]*conn
}

func (r *Global) Add(id string, c *conn) {
	r.l.Lock()
	defer r.l.Unlock()

	r.conns[id] = c
}

func (r *Global) Del(id string) {
	r.l.Lock()
	defer r.l.Unlock()

	delete(r.conns, id)
}

func (r *Global) Get(id string) *conn {
	r.l.RLock()
	defer r.l.RLock()

	return r.conns[id]
}

func (r *Global) Send(msg []byte) {
	r.l.RLock()
	defer r.l.RLock()

	for _, conn := range r.conns {
		conn.Send(msg)
	}
}

//-----------------------

type Room struct {
	l     sync.RWMutex
	num   int
	conns map[*conn]struct{}
}

func initRoom() *Room {
	r := Room{}
	r.conns = make(map[*conn]struct{}, 5)
	return r
}

func (r *Room) Add(c *conn) {
	r.l.Lock()
	defer r.l.Unlock()

	r.num++
	r.conns[c] = make(struct{})
}

func (r *Room) Del(c *conn) bool {
	r.l.Lock()
	defer r.l.Unlock()

	r.num--
	delete(r.conns, c)

	if r.num < 1 {
		return true
	}

	return false
}

func (r *Room) Send(msg []byte) {
	r.l.RLock()
	defer r.l.RLock()

	for conn, _ := range r.conns {
		conn.Send(msg)
	}
}

//---------------------------

type Rooms struct {
	l     sync.Mutex
	rooms map[string]*Room
}

func (m *Rooms) Get(id string) *Room {
	m.l.Lock()
	defer m.l.Unlock()
	r, ok := m.rooms[id]
	if !ok {
		m.rooms[id] = initRoom()
	}
	return r
}

func (m *Room) Add(id string, r *room) {
	m.l.Lock()
	defer m.l.Unlock()

	m.rooms[id] = r
}

func (m *Room) Del(id string) {
	m.l.Lock()
	defer m.l.Unlock()

	delete(m.conns, id)
}

func (m *Room) GC(id string) {
	m.l.Lock()
	defer m.l.Unlock()

	r := m.rooms[id]
	if r != nil {
		if r.num < 1 {
			delete(m.rooms, id)
		}
	}
}
