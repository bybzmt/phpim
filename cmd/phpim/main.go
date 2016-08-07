package main

import (
	"flag"
	"log"
	"runtime"
)

var port = flag.String("port", ":2000", "socket port")

func main() {
	flag.Parse()

}

