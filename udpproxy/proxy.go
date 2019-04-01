package udpproxy

import (
	log "github.com/sirupsen/logrus"
	"math/rand"
)

type Proxy struct {
	ServerPorts       []int
	ClientPorts       []int
	ServerChannel     chan []byte
	ClientChannel     chan []byte
	ServerConnections []*ServerConnection
	ClientConnections []*ClientConnection
}

func NewProxy(serverPort []int, clientPort []int) *Proxy {
	return &Proxy{
		ServerPorts: serverPort, ClientPorts: clientPort,
		ServerChannel: make(chan []byte, UDPBuffer), ClientChannel: make(chan []byte, UDPBuffer),
		ServerConnections: make([]*ServerConnection, 0, len(serverPort)),
		ClientConnections: make([]*ClientConnection, 0, len(clientPort)),
	}
}

func (proxy *Proxy) Init() {
	for _, port := range proxy.ServerPorts {
		connection := NewServerConnection(port)
		if err := connection.Listen(); err != nil {
			log.Fatalf("Can not listen port: %d", connection.Port)
		}
		proxy.ServerConnections = append(proxy.ServerConnections, connection)
	}

	for _, port := range proxy.ClientPorts {
		connection := NewClientConnection(port)
		if err := connection.Connect(); err != nil {
			log.Fatalf("Can not connect to port: %d", connection.Port)
		}
		proxy.ClientConnections = append(proxy.ClientConnections, connection)
	}
}

func (proxy *Proxy) Start() {
	for _, connection := range proxy.ServerConnections {
		go proxy.listenServer(connection)
	}
	for _, connection := range proxy.ClientConnections {
		go proxy.listenClient(connection)
	}
	go proxy.proxy()
}

func (proxy *Proxy) listenServer(connection Connection) {
	for {
		var data [MTU]byte
		if n, err := connection.ReadFrom(data[:]); err == nil {
			proxy.ServerChannel <- data[:n]
		} else {
			log.Fatalf("can not recv from server, %e", err)
		}
	}
}

func (proxy *Proxy) listenClient(connection Connection) {
	for {
		var data [MTU]byte
		if n, err := connection.ReadFrom(data[:]); err == nil {
			proxy.ClientChannel <- data[:n]
		} else {
			log.Fatalf("can not recv from client, %e", err)
		}
	}
}

func (proxy *Proxy) proxy() {
	for {
		select {
		case data := <-proxy.ClientChannel:
			i := rand.Intn(len(proxy.ServerConnections))
			if _, err := proxy.ServerConnections[i].SendTo(data); err != nil {
				log.Fatal("Send error")
			}
		case data := <-proxy.ServerChannel:
			i := rand.Intn(len(proxy.ClientConnections))
			if _, err := proxy.ClientConnections[i].SendTo(data); err != nil {
				log.Fatal("Send error")
			}
		}
	}
}
