package main

import (
	"flag"
	"github.com/bybzmt/phpim"
	"log"
	"net"
	"net/http"
	"runtime"
	"strings"
)

var addr = flag.String("addr", ":2000", "listen addr:port")
var callback = flag.String("callback", "", "callback url")
var origin = flag.String("origin", "*", "allow origin domain")
var localIP = flag.String("localIP", "127.0.0.0/8,10.0.0.0/8,172.17.0.0/16,192.168.0.0/16", "allow local ip access")
var maxConn = flag.Int("maxConn", 10000, "max conn num")
var maxSingleIP = flag.Int("maxSingleIP", 5, "max single ip conn num")
var realip = flag.String("realip", "", "real ip header")

func main() {
	flag.Parse()

	runtime.GOMAXPROCS(runtime.NumCPU())

	if *callback == "" {
		log.Fatalln("callback url can not empty.")
	}

	IPNets := make([]net.IPNet, 0, 1)
	for _, ip := range strings.Split(*localIP, ",") {
		_, IPNet, err := net.ParseCIDR(ip)
		if err == nil {
			IPNets = append(IPNets, *IPNet)
		}
	}

	origins := strings.Split(*origin, ",")

	im := phpim.NewIM()
	im.Origins = origins
	im.LocalIPs = IPNets
	im.MaxSingleIP = int16(*maxSingleIP)
	im.MaxConn = int32(*maxConn)
	im.CallbackUrl = *callback
	im.RealIP = *realip

	http.HandleFunc("/actions", im.SendMsg)
	http.HandleFunc("/ws", im.ServeWs)

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
