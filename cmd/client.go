package main

import (
	"flag"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"runtime"
	"sync/atomic"
	"time"
)

var num = flag.Int("num", 10, "num for clients")
var serverUrl = flag.String("url", "ws://127.0.0.1:2000/ws", "websocket server url")
var origin = flag.String("origin", "http://127.0.0.1/", "websocket origin")

var count int64
var msg []byte

func client() {
	header := http.Header{}
	header.Add("origin", *origin)

	ws, _, err := websocket.DefaultDialer.Dial(*serverUrl, header)

	if err != nil {
		log.Println("a", err)
		return
	}
	defer ws.Close()

	for {
		_, msg, err = ws.ReadMessage()
		if err != nil {
			log.Println("b", err)
			break
		}
		log.Println(string(msg))
		atomic.AddInt64(&count, 1)
	}

}

func main() {
	flag.Parse()

	runtime.GOMAXPROCS(runtime.NumCPU())

	for i := 0; i < *num; i++ {
		time.Sleep(30 * time.Millisecond)

		go client()
	}

	c := time.Tick(3 * time.Second)
	for _ = range c {
		now()
	}
}

func now() {
	num := atomic.LoadInt64(&count)

	log.Println("message Received count", num, "msg:", string(msg))
}
