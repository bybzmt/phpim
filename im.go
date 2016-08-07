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
	AllowLocalIP []string
	AllowOrigin  []string
	MaxSingleIP  int16
	IPCounter    map[net.IP]int16
	CallbackUrl  string
	ConnNum      int32
	MaxConn      int32
	Conns        Global
	Rooms        Rooms
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

	num := IPCounter[ip]
	if num < 1 {
		delete(im.IPCounter, ip)
	}

	return num
}

func (im *IM) checkOrigin(r *http.Request) bool {
	for _, y := range im.AllowOrigin {
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

// serverWs handles webocket requests from the peer.
func (im *IM) ServeWs(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	now := atomic.AddInt32(&im.ConnNum, 1)

	defer func() {
		atomic.AddInt32(&im.ConnNum, -1)
	}()

	if now > im.MaxConn {
		log.Println("Server Is Too Busy!", r.RemoteAddr)
		http.Error(w, "Service Unavailable", 503)
		return
	}

	ip := net.ParseIP(strings.Split(RemoteAddr, ":")[0])
	num := im.addIPCounter(ip, 1)
	defer im.addIPCounter(ip, -1)

	if num > im.MaxSingleIp {
		log.Println("connections over", im.MaxSingleIP, r.RemoteAddr)
		http.Error(w, "ip not allowed", 403)
		return
	}

	ws, err := im.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err, r.RemoteAddr, r.Header.Get("User-Agent"))
		return
	}

	c := &connection{
		send: make(chan []byte, 256),
		ws: ws,
		writeWait: 10 * time.Second,
		pingPeriod: (pongWait * 9) / 10,
		maxMessageSize: 512
	}

	err := im.connectCallback(c)
	if err != nil {
		return;
	}

	defer im.disconnectCallback(ws)

	go c.writePump()

	c.readPump(im)
}
