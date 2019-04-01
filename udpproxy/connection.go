package udpproxy

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net"
)

type Connection interface {
    SendTo(data []byte) (n int, err error)
	ReadFrom(data []byte) (n int, err error)
}

type ServerConnection struct {
	Connection
	Host string
	Port int
	Addr *net.UDPAddr
	Socket *net.UDPConn
}

func NewServerConnection(host string, port int) *ServerConnection {
	return &ServerConnection{Host: host, Port: port}
}

func (conn *ServerConnection) Listen() (err error) {
	log.Infof("start listening %s: %d", conn.Host, conn.Port)
	addrStr := fmt.Sprintf("%s:%d", conn.Host, conn.Port)
	addr, err := net.ResolveUDPAddr(UDPProtocol, addrStr)
	if err != nil {
		log.Fatalf("Can not listen: %e", err)
		return
	}

	socket, err := net.ListenUDP(UDPProtocol, addr)
	if err != nil {
		log.Fatalf("Can not listen: %e", err)
		return
	}
	conn.Socket = socket
	return nil
}

func (conn *ServerConnection) SendTo(data []byte) (n int, err error) {
	if conn.Addr != nil {
		if n, err = conn.Socket.WriteToUDP(data, conn.Addr); err != nil {
			log.Fatalf("Failed to send data to %d", conn.Port)
		}
		log.Debugf("Send data to %d", conn.Port)
	}
	return 0, nil
}

func (conn *ServerConnection) ReadFrom(data []byte) (n int, err error) {
	if n, conn.Addr, err = conn.Socket.ReadFromUDP(data); err != nil {
		log.Fatalf("Failed to read from %d", conn.Port)
	}
	log.Debugf("Read data from %d", conn.Port)
	return n, err
}

type ClientConnection struct{
	Connection
	Host string
	Port int
	Socket *net.UDPConn
}

func NewClientConnection(host string, port int) *ClientConnection {
	return &ClientConnection{Host: host, Port:port}
}

func (conn *ClientConnection) Connect() (err error) {
	log.Printf("starting connect to %s:%d", conn.Host, conn.Port)
	addrStr := fmt.Sprintf("%s:%d", conn.Host, conn.Port)
	addr, err := net.ResolveUDPAddr(UDPProtocol, addrStr)
	if err != nil {
		log.Fatalf("Can not connect to server: %e", err)
		return
	}

	socket, err := net.DialUDP(UDPProtocol, nil, addr)
	if err != nil {
		log.Fatalf("Can not connect to server: %e", err)
		return
	}
	conn.Socket = socket
	return nil
}

func (conn *ClientConnection) SendTo(data []byte) (n int, err error) {
	log.Debugf("Send data to %d", conn.Port)
	return conn.Socket.Write(data)
}

func (conn *ClientConnection) ReadFrom(data []byte) (n int, err error) {
	log.Debugf("Read data from %d", conn.Port)
	return conn.Socket.Read(data)
}
