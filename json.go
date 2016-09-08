package phpim

import (
	"encoding/json"
)

type CallbackResponse struct {
	Ret     int
	Id      string
	Actions []Action
}

type Action struct {
	Type  int
	Point json.RawMessage
}

type ActionSendMsg struct {
	Conn string
	Msg  string
}

type ActionCloseConn struct {
	Conn string
}

type ActionBroadcast struct {
	Msg string
}

type ActionRoomBroadcast struct {
	Room string
	Msg  string
}

type ActionConnAddRoom struct {
	Conn string
	Room string
}

type ActionRoomDelConn struct {
	Room string
	Conn string
}

type Response struct {
	Ret int
	Msg string
}
