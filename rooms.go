package phpim

import (
	"sync"
)

type Global struct {
	l     sync.RWMutex
	num   int32
	conns map[string]*connection
}

func (g *Global) init() {
	g.conns = make(map[string]*connection, 10000)
}

func (r *Global) Add(id string, c *connection) {
	r.l.Lock()
	defer r.l.Unlock()

	r.num++
	r.conns[id] = c
}

func (r *Global) Del(id string) {
	r.l.Lock()
	defer r.l.Unlock()

	r.num--
	delete(r.conns, id)
}

func (r *Global) Get(id string) *connection {
	r.l.RLock()
	defer r.l.RUnlock()

	return r.conns[id]
}

func (r *Global) Send(msg []byte) {
	r.l.RLock()
	defer r.l.RUnlock()

	for _, conn := range r.conns {
		conn.Send(msg)
	}
}

//-----------------------

type Room struct {
	l     sync.RWMutex
	num   int
	name  string
	conns map[*connection]struct{}
}

func (r *Room) init(name string) *Room {
	r.conns = make(map[*connection]struct{}, 100)
	r.name = name
	return r
}

func (r *Room) Add(c *connection) {
	r.l.Lock()
	defer r.l.Unlock()

	r.num++
	r.conns[c] = struct{}{}
}

func (r *Room) Del(c *connection) bool {
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
	defer r.l.RUnlock()

	for conn, _ := range r.conns {
		conn.Send(msg)
	}
}

//---------------------------

type Rooms struct {
	l     sync.Mutex
	rooms map[string]*Room
}

func (r *Rooms) init() {
	r.rooms = make(map[string]*Room, 100)
}

func (m *Rooms) Get(id string) *Room {
	m.l.Lock()
	defer m.l.Unlock()
	r, ok := m.rooms[id]
	if !ok {
		r = new(Room).init(id)
		m.rooms[id] = r
	}
	return r
}

func (m *Rooms) Add(id string, r *Room) {
	m.l.Lock()
	defer m.l.Unlock()

	m.rooms[id] = r
}

func (m *Rooms) Del(id string) {
	m.l.Lock()
	defer m.l.Unlock()

	delete(m.rooms, id)
}

func (m *Rooms) GC(id string) {
	m.l.Lock()
	defer m.l.Unlock()

	r := m.rooms[id]
	if r != nil {
		if r.num < 1 {
			delete(m.rooms, id)
		}
	}
}
