package main

import (
	"flag"
	"log"
	"runtime"
	"strings"
)

var addr = flag.String("addr", ":2000", "listen addr:port")
var callback = flag.String("callback", "", "callback url")
var origin = flag.String("origin", "localhost", "allow origin domain")
var localIP = flag.String("localIP", "127.0.0.0/8,10.0.0.0/8,172.17.0.0/16,192.168.0.0/16", "allow local ip access")
var maxConn = flag.Int("maxConn", 10000, "max conn num")
var maxSingleIP = flag.Int("maxSingleIP", 5, "max single ip conn num")

func main() {
	flag.Parse()

	if *callback == "" {
		log.Falteln("callback url can not empty.")
	}

	IPNets = make([]*net.IPNet, 0, 1)
	for _, ip := range strings.Split(*localIP, ",") {
		_, IPNet, err := net.ParseCIDR(ip)
		if err == nil {
			IPNets = append(IPNets, IPNet)
		}
	}

	origins = strings.Split(*origin, ",")

	im := phpim.NewIM()
	im.Origins = origins
	im.LocalIPs = IPNets
	im.MaxSingleIP = *maxSingleIP
	im.MaxConn = *maxConn

}

