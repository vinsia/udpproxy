package udpproxy

import (
	"fmt"
	"log"
	"net"
)

type Connection interface {
    SendTo(data []byte) (n int, err error)
	ReadFrom(data []byte) (n int, err error)
}

type ServerConnection struct {
	Connection
	Port int
	Addr *net.UDPAddr
	Socket *net.UDPConn
}

func NewServerConnection(port int) *ServerConnection {
	return &ServerConnection{Port: port}
}

func (conn *ServerConnection) Listen() (err error) {
	log.Printf("start listening 127.0.0.1: %d", conn.Port)
	addrStr := fmt.Sprintf("%s:%d", "127.0.0.1", conn.Port)
	addr, err := net.ResolveUDPAddr(UDPProtocol, addrStr)
	if err != nil {
		log.Panicf("can not listen: %e", err)
		return
	}

	socket, err := net.ListenUDP(UDPProtocol, addr)
	if err != nil {
		log.Panicf("can not listen: %e", err)
		return
	}
	conn.Socket = socket
	return nil
}

func (conn *ServerConnection) SendTo(data []byte) (n int, err error) {
	if conn.Addr != nil {
		return conn.Socket.WriteToUDP(data, conn.Addr)
	}
	return 0, nil
}

func (conn *ServerConnection) ReadFrom(data []byte) (n int, err error) {
	n, conn.Addr, err = conn.Socket.ReadFromUDP(data)
	return
}

type ClientConnection struct{
	Connection
	Port int
	Socket *net.UDPConn
}

func NewClientConnection(port int) *ClientConnection {
	return &ClientConnection{Port:port}
}

func (conn *ClientConnection) Connect() (err error) {
	log.Printf("starting connect to :%d", conn.Port)
	addrStr := fmt.Sprintf("127.0.0.1:%d", conn.Port)
	addr, err := net.ResolveUDPAddr(UDPProtocol, addrStr)
	if err != nil {
		log.Panicf("can not connect to server: %e", err)
		return
	}

	socket, err := net.DialUDP(UDPProtocol, nil, addr)
	if err != nil {
		log.Panicf("can not connect to server: %e", err)
		return
	}
	conn.Socket = socket
	return nil
}

func (conn *ClientConnection) SendTo(data []byte) (n int, err error) {
	return conn.Socket.Write(data)
}

func (conn *ClientConnection) ReadFrom(data []byte) (n int, err error) {
	return conn.Socket.Read(data)
}
