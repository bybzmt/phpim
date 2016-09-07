package phpim

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"strings"
	"sync/atomic"
)

// serverWs handles webocket requests from the peer.
func (im *IM) ServeWs(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	if atomic.LoadInt32(&im.conns.num) > im.MaxConn {
		log.Println("Server Is Too Busy!", r.RemoteAddr)
		http.Error(w, "Service Unavailable", 503)
		return
	}

	ip := net.ParseIP(strings.Split(r.RemoteAddr, ":")[0])
	num := im.addIPCounter(ip, 1)
	defer im.addIPCounter(ip, -1)

	if num > im.MaxSingleIP {
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

	err = im.connectCallback(c, r)
	if err != nil {
		log.Println(err)
		return
	}

	defer func() {
		err := im.disconnectCallback(c)
		if err != nil {
			log.Println(err)
		}
	}()

	go c.writePump()

	err = c.readPump(im)
	if err != nil {
		log.Println(err)
	}
}

func (im *IM) SendMsg(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if e := recover(); e != nil {
			if b, ok := e.(BadRequest); ok {
				w.WriteHeader(http.StatusBadRequest)
				enc := json.NewEncoder(w)
				enc.Encode(Response{Ret: 1, Msg: string(b)})
			} else {
				log.Println(e)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}
	}()

	ip := net.ParseIP(strings.Split(r.RemoteAddr, ":")[0])

	allow := false
	for _, ipnet := range im.LocalIPs {
		if ipnet.Contains(ip) {
			allow = true
			break
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

	im.serveAction(actions)

	enc := json.NewEncoder(w)
	enc.Encode(Response{Ret: 0, Msg: "success."})
}
