package main

import (
	"github.com/elgohr/go-tcp-proxy/proxy"
	"log"
	"net"
	"os"
	"os/signal"
)

func main() {
	localAddr := os.Args[1]
	remoteAddr := os.Args[2]

	localAddress, err := net.ResolveTCPAddr("tcp", localAddr)
	if err != nil {
		log.Fatalln(err)
	}
	remoteAdress, err := net.ResolveTCPAddr("tcp", remoteAddr)
	if err != nil {
		log.Fatalln(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Kill, os.Interrupt)

	p, err := proxy.NewProxy(c, localAddress, remoteAdress)
	p.Run()
}
