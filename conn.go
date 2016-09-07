package phpim

import (
	"github.com/gorilla/websocket"
	"sync"
	"time"
)

type connection struct {
	L  sync.Mutex
	im *IM
	ws *websocket.Conn
	//连接的唯一名字
	Id    string
	send  chan []byte
	rooms []*Room
}

func NewConn(ws *websocket.Conn, im *IM) *connection {
	c := &connection{
		send: make(chan []byte, 256),
		ws:   ws,
		im:   im,
	}
	return c
}

func (c *connection) readPump(im *IM) error {
	defer func() {
		c.ws.Close()
	}()
	c.ws.SetReadLimit(c.im.MaxMsgSize)
	c.ws.SetReadDeadline(time.Now().Add(c.im.pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(c.im.pongWait)); return nil })
	for {
		//直接丢换收到的用户数据
		ty, msg, err := c.ws.ReadMessage()
		if err != nil {
			break
		}

		switch ty {
		case websocket.TextMessage:
			err = im.msgCallback(c, string(msg))
			if err != nil {
				return err
			}
		case websocket.BinaryMessage:
		case websocket.CloseMessage:
			break
		}
	}

	return nil
}

func (c *connection) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(c.im.writeWait))
	return c.ws.WriteMessage(mt, payload)
}

func (c *connection) writePump() {
	ticker := time.NewTicker(c.im.pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.write(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func (c *connection) Send(m []byte) {
	select {
	case c.send <- m:
	default:
		c.ws.Close()
	}
}

func (c *connection) Close() {
	c.ws.Close()
}

func (c *connection) AddRoom(r *Room) {
	c.L.Lock()
	defer c.L.Unlock()

	c.rooms = append(c.rooms, r)
}

func (c *connection) DelRoom(r *Room) {
	c.L.Lock()
	defer c.L.Unlock()

	t := []*Room{}
	for _, o := range c.rooms {
		if o != r {
			t = append(t, o)
		}
	}
	c.rooms = t
}
