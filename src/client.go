package main

import (
	"net"
)

type ClientMessageHandler func(buf []byte)
type ClientErrorHandler func(err error)

type Client interface {
	Connect(ClientMessageHandler, ClientErrorHandler)
	Send([]byte) error
	Close()
}

type UDPClient struct {
	udpAddr *net.UDPAddr
	conn    *net.UDPConn
	closed  bool
}

func NewUDPClient(serverAddr string) (*UDPClient, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", serverAddr)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return nil, err
	}

	return &UDPClient{
		udpAddr: udpAddr,
		conn:    conn,
	}, nil
}

func (c *UDPClient) Connect(messageHandler ClientMessageHandler, errorhandler ClientErrorHandler) {
	defer c.conn.Close()

	buff := make([]byte, 1024)
	for !c.closed {
		n, err := c.conn.Read(buff)
		if err != nil {
			errorhandler(err)
			break
		}

		messageHandler(buff[:n])
	}
}

func (c *UDPClient) Send(buf []byte) error {
	_, err := c.conn.Write(buf)
	if err != nil {
		return err
	}

	return nil
}

func (c *UDPClient) Close() {
	c.closed = true
}
