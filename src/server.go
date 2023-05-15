package main

import (
	"net"
)

type ServerMessageHandler func(net.Addr, []byte)
type ServerErrorHandler func(error)

type Server interface {
	Listen(ServerMessageHandler, ServerErrorHandler)
	Send(net.Addr, []byte) error
	Close()
	GetAddr() net.Addr
}

type UDPServer struct {
	udpAddr *net.UDPAddr
	conn    *net.UDPConn
	closed  bool
}

func NewUDPServer(listenAddr string) (*UDPServer, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", listenAddr)
	if err != nil {
		return nil, err
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return nil, err
	}

	return &UDPServer{
		udpAddr: udpAddr,
		conn:    conn,
	}, nil
}

func (s *UDPServer) Listen(messageHandler ServerMessageHandler, errorHandler ServerErrorHandler) {
	defer s.conn.Close()

	buff := make([]byte, 1024)

	for !s.closed {
		n, addr, err := s.conn.ReadFrom(buff)
		if err != nil {
			errorHandler(err)
			continue
		}

		messageHandler(addr, buff[:n])
	}
}

func (s *UDPServer) Send(addr net.Addr, buff []byte) error {
	_, err := s.conn.WriteTo(buff, addr)
	if err != nil {
		return err
	}

	return nil
}

func (s *UDPServer) Close() {
	s.closed = true
}

func (s *UDPServer) GetAddr() net.Addr {
	return s.udpAddr
}
