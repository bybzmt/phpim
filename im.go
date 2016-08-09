package phpim

import (
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type IM struct {
	L            sync.Mutex
	LocalIPs []net.IPNet
	Origins  []string
	MaxSingleIP  int16
	IPCounter    map[net.IP]int16
	CallbackUrl  string
	MaxConn      int32
	conns        Global
	rooms        Rooms
	MaxMsgSize   int
	writeWait    time.Duration
	pongWait     time.Duration
	pingPeriod   time.Duration
	upgrader websocket.Upgrader
}

func NewIM() *IM {
	im := IM{}
	im.IpCounter = make(map[net.IP]int16, 2000)
	im.MaxMsgSize = 512

	im.global.init()
	im.rooms.init()

	im.writeWait = 10 * time.Second
	im.pongWait = 60 * time.Second
	im.pingPeriod = (im.pongWait * 9) / 10

	im.upgrader = websocket.Upgrader{
		ReadBufferSize:   1024,
		WriteBufferSize:  1024,
		HandshakeTimeout: 5 * time.Second,
		CheckOrigin: im.checkOrigin
	}

	return im
}

func (im *IM) addIPCounter(ip net.IP, i int16) int16 {
	im.L.Lock()
	defer im.L.Unlock()

	im.IPCounter[ip] += i

	num := im.IPCounter[ip]
	if num < 1 {
		delete(im.IPCounter, ip)
	}

	return num
}

func (im *IM) checkOrigin(r *http.Request) bool {
	for _, y := range im.Origins {
		if y == "*" {
			return true
		}

		u, err := url.Parse(r.Header.Get("Origin"))
		if err != nil {
			return false
		}

		if u.Host == y {
			return true
		}
	}
	return false
}

