package main

import "net"

type MessageHandlerFunc func(buf []byte, size int)
type ErrorHandlerFunc func(err error)

type UDPClient struct {
	udpAddr *net.UDPAddr
	conn    *net.UDPConn
}

func NewUDPClient(serverAddr string) (*UDPClient, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", serverAddr)
	if err != nil {
		return nil, err
	}

	return &UDPClient{
		udpAddr: udpAddr,
	}, nil
}

func (c *UDPClient) Connect(messageHandler MessageHandlerFunc, errorhandler ErrorHandlerFunc) error {
	conn, err := net.DialUDP("udp", nil, c.udpAddr)
	if err != nil {
		return err
	}

	c.conn = conn

	go func() {
		for {
			buf := make([]byte, 1024)
			size, err := conn.Read(buf)
			if err != nil {
				errorhandler(err)
				break
			}
			messageHandler(buf, size)
		}
	}()

	return nil
}

func (c *UDPClient) Send(buf []byte) error {
	_, err := c.conn.Write(buf)
	if err != nil {
		return err
	}

	return nil
}

func (c *UDPClient) Close() {
	c.conn.Close()
}
