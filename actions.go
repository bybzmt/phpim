package phpim

import (
	"encoding/json"
)

const (
	SendMsg = iota + 1
	CloseConn
	Broadcast
	RoomBroadcast
	ConnAddRoom
	RoomDelConn
)

type BadRequest string

func (e BadRequest) Error() {
	return e
}

func (im *IM) serveAction(actions []Action) {
	for _, action := range actions {
		switch action.Type {
		case SendMsg:
			im.serveSendMsg(action.Point)
		case CloseConn:
			im.serveCloseConn(action.Point)
		case Broadcast:
			im.serveBroadcast(action.Point)
		case RoomBroadcast:
			im.serveRoomBroadcast(action.Point)
		case ConnAddRoom:
			im.serveConnAddRoom(action.Point)
		case RoomDelConn:
			im.serveConnAddRoom(action.Point)
		default:
			panic(BadRequest("action type err."))
		}
	}
}

func (im *IM) serveSendMsg(raw json.RawMessage) {
	ac := ActionSendMsg{}
	err := json.Unmarshal(raw, &ac)
	if err != nil {
		panic(BadRequest("action json decode err"))
	}

	c := im.conns.Get(ac.Conn)
	if c != nil {
		panic(BadRequest("conn not found."))
	}

	c.Send(ac.Msg)
}

func (im *IM) serveCloseConn(raw json.RawMessage) {
	ac := ActionCloseConn{}
	err := json.Unmarshal(raw, &ac)
	if err != nil {
		panic(BadRequest("action json decode err"))
	}

	c := im.conns.Get(ac.Conn)
	if c != nil {
		c.Close()
	}
}

func (im *IM) serveBroadcast(raw json.RawMessage) {
	ac := ActionBroadcast{}
	err := json.Unmarshal(raw, &ac)
	if err != nil {
		panic(BadRequest("action json decode err"))
	}

	im.conns.Send(ac.Msg)
}

func (im *IM) serveRoomBroadcast(raw json.RawMessage) {
	ac := ActionRoomBroadcast{}
	err := json.Unmarshal(raw, &ac)
	if err != nil {
		panic(BadRequest("action json decode err"))
	}

	r := im.rooms.Get(ac.Room)
	if r == nil {
		panic(BadRequest("room not found."))
	}

	r.Send(ac.Msg)
}

func (im *IM) serveConnAddRoom(raw json.RawMessage) {
	ac := ActionConnAddRoom{}
	err := json.Unmarshal(raw, &ac)
	if err != nil {
		panic(BadRequest("action json decode err"))
	}

	r := im.rooms.Get(ac.Room)
	c := im.conns.Get(ac.Conn)
	if r == nil {
		panic(BadRequest("room not found."))
	}

	if c == nil {
		panic(BadRequest("conn not found."))
	}

	r.Add(c)
	c.AddRoom(r)
}

func (im *IM) serveConnAddRoom(raw json.RawMessage) {
	ac := ActionRoomDelConn{}
	err := json.Unmarshal(raw, &ac)
	if err != nil {
		panic(BadRequest("action json decode err"))
	}

	r := im.rooms.Get(ac.Room)
	c := im.conns.Get(ac.Conn)
	if r == nil {
		panic(BadRequest("room not found."))
	}

	if c == nil {
		panic(BadRequest("conn not found."))
	}

	empty := r.Del(c)
	c.delRoom(r)
	if empty {
		im.rooms.GC(ac.Room)
	}
}
