package phpim

import (
	"encoding/json"
)

const (
	SendMsg = iota + 1
	CloseConn
	Broadcast
	RoomBroadcast
	RoomAddConn
	RoomDelConn
)

type BadRequest string

func (e BadRequest) Error() string {
	return string(e)
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
		case RoomAddConn:
			im.serveRoomAddConn(action.Point)
		case RoomDelConn:
			im.serveRoomDelConn(action.Point)
		default:
			panic(BadRequest("action type err."))
		}
	}
}

func (im *IM) serveSendMsg(raw json.RawMessage) {
	ac := ActionSendMsg{}
	err := json.Unmarshal(raw, &ac)
	if err != nil {
		panic(BadRequest("SendMsg action json decode err:" + err.Error()))
	}

	c := im.conns.Get(ac.Conn)
	if c == nil {
		//panic(BadRequest("SendMsg conn not found."))
		return
	}

	c.Send([]byte(ac.Msg))
}

func (im *IM) serveCloseConn(raw json.RawMessage) {
	ac := ActionCloseConn{}
	err := json.Unmarshal(raw, &ac)
	if err != nil {
		panic(BadRequest("CloseConn action json decode err" + err.Error()))
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
		panic(BadRequest("Broadcast action json decode err:" + err.Error()))
	}

	im.conns.Send([]byte(ac.Msg))
}

func (im *IM) serveRoomBroadcast(raw json.RawMessage) {
	ac := ActionRoomBroadcast{}
	err := json.Unmarshal(raw, &ac)
	if err != nil {
		panic(BadRequest("RoomBroadcast action json decode err" + err.Error()))
	}

	r := im.rooms.Get(ac.Room)
	if r == nil {
		//panic(BadRequest("RoomBroadcast room not found."))
		return
	}

	r.Send([]byte(ac.Msg))

	if r.num < 1 {
		im.rooms.GC(r.name)
	}
}

func (im *IM) serveRoomAddConn(raw json.RawMessage) {
	ac := ActionConnAddRoom{}
	err := json.Unmarshal(raw, &ac)
	if err != nil {
		panic(BadRequest("RoomAddConn action json decode err" + err.Error()))
	}

	r := im.rooms.Get(ac.Room)
	c := im.conns.Get(ac.Conn)
	if r == nil {
		//panic(BadRequest("RoomAddConn room not found."))
		return
	}

	if c == nil {
		//panic(BadRequest("RoomAddConn conn not found."))
		return
	}

	r.Add(c)
	c.AddRoom(r)
}

func (im *IM) serveRoomDelConn(raw json.RawMessage) {
	ac := ActionRoomDelConn{}
	err := json.Unmarshal(raw, &ac)
	if err != nil {
		panic(BadRequest("RoomDelConn action json decode err" + err.Error()))
	}

	r := im.rooms.Get(ac.Room)
	c := im.conns.Get(ac.Conn)
	if r == nil {
		//panic(BadRequest("RoomDelConn room not found."))
		return
	}

	if c == nil {
		//panic(BadRequest("RoomDelConn conn not found."))
		return
	}

	empty := r.Del(c)
	c.DelRoom(r)
	if empty {
		im.rooms.GC(r.name)
	}
}
