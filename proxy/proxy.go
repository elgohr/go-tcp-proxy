package proxy

import (
	"github.com/egymgmbh/go-prefix-writer/prefixer"
	"io"
	"log"
	"net"
	"os"
)

type Proxy struct {
	c <-chan os.Signal

	listener *net.TCPListener
	remote   *net.TCPAddr

	localConnection  *net.TCPConn
	remoteConnection *net.TCPConn
}

func NewProxy(c <-chan os.Signal, local *net.TCPAddr, remote *net.TCPAddr) (*Proxy, error) {
	listener, err := net.ListenTCP("tcp", local)
	if err != nil {
		return nil, err
	}
	return &Proxy{
		c:        c,
		listener: listener,
		remote:   remote,
	}, nil
}

func (p Proxy) Run() {
	for {
		select {
		case <-p.c:
			log.Println("exiting...")
			return
		default:
			localConnection, err := p.listener.AcceptTCP()
			if err != nil {
				log.Println(err)
			}

			defer localConnection.Close()
			remoteConnection, err := net.DialTCP("tcp", nil, p.remote)
			if err != nil {
				log.Println(err)
			}
			defer remoteConnection.Close()

			go pipe(localConnection, remoteConnection, true)
			go pipe(remoteConnection, localConnection, false)
		}
	}
}

func pipe(src, dst io.ReadWriter, local bool) {
	var prefix string
	if local {
		prefix = ">>  "
	} else {
		prefix = "<<  "
	}
	if _, err := io.Copy(io.MultiWriter(dst, prefixer.New(os.Stdout, func() string {
		return prefix
	})), src); err != nil {
		log.Println(err)
		return
	}
}
