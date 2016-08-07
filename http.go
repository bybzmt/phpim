package phpim

import (
	"log"
	"net"
	"net/http"
	"strings"
)

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

	c := NewConn(ws, im)

	err := im.connectCallback(c)
	if err != nil {
		return;
	}

	defer im.disconnectCallback(ws)

	go c.writePump()

	c.readPump(im)
}

func (im *IM) ServeAction(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if e := recover(); e != nil {
			if b, ok := e.(BadRequest); ok {
				w.WriteHeader(http.StatusBadRequest)
				enc := json.NewEncoder(w)
				enc.Encode(Response{ret: 1, msg: string(b)})
			} else {
				log.Println(e)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}
	}()

	allow = false
	for _, ip := range im.AllowLocalIP {
		if strings.Index(r.RemoteAddr, ip) == 0 {
			allow = true
			return
		}
	}
	if !allow {
		log.Println("Unauthorized Visit:", r.RemoteAddr)
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	var actions []Action
	err := json.NewDecoder(r.Body).Decode(actions)
	if err != nil {
		panic(BadRequest("json decode err."))
	}

	serveAction(actions)

	enc := json.NewEncoder(w)
	enc.Encode(Response{ret: 0, msg: "success."})
}
