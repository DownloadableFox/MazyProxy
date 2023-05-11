package main

import (
	"fmt"
	"net"
)

const (
	PROXY_START_EVENT      = 0x0
	PROXY_CONNECTION_EVENT = 0x1
	PROXY_CLOSE_EVENT      = 0x2
	PROXY_SEND_EVENT       = 0x3
	PROXY_RECV_EVENT       = 0x4
)

type ProxyEvent struct {
	proxy  *Proxy
	client *UDPClient
	buffer []byte
	size   int
}

type ProxyEventFunc func(handler ProxyEvent) ([]byte, bool)

type Proxy struct {
	listenAddr   string
	redirectAddr string
	server       net.PacketConn
	clients      map[string]*UDPClient
	middlewares  map[uint8][]ProxyEventFunc
}

func NewProxy(listenAddr, redirectAddr string) (*Proxy, error) {
	server, err := net.ListenPacket("udp", listenAddr)
	if err != nil {
		return nil, err
	}

	return &Proxy{
		listenAddr:   listenAddr,
		redirectAddr: redirectAddr,
		server:       server,
		clients:      make(map[string]*UDPClient),
		middlewares:  make(map[uint8][]ProxyEventFunc),
	}, nil
}

func (p *Proxy) Start() error {
	fmt.Printf("Listening on %s\n", p.server.LocalAddr())
	defer p.Close()

	for {
		// Read from the server
		buffer := make([]byte, 1024)
		size, address, err := p.server.ReadFrom(buffer)
		if err != nil {
			return err
		}

		// Check if we have a client for this address
		client, err := p.getClient(address)
		if err != nil {
			return err
		}

		decryptedBuff := NetworkDecrypt(buffer[:size])
		if decryptedBuff[0] == PACKET_EMOJI {
			fmt.Printf("Recieved %d bytes from %s: %v\n", size, address, decryptedBuff)
		}
		client.Send(buffer[:size])
	}
}

func (p *Proxy) AddEvents(eventType uint8, handlers ...ProxyEventFunc) {
	p.middlewares[eventType] = append(p.middlewares[eventType], handlers...)
}

func (p *Proxy) Close() {
	p.server.Close()
	for _, client := range p.clients {
		client.Close()
	}
}

func (p *Proxy) getClient(address net.Addr) (*UDPClient, error) {
	addressStr := address.String()
	client, ok := p.clients[addressStr]
	if !ok {
		// Create a new newclient
		newclient, err := NewUDPClient(p.redirectAddr)
		if err != nil {
			return nil, err
		}
		client = newclient

		// Connect the client
		if err := client.Connect(p.handleMessage(address), p.handleError(address)); err != nil {
			return nil, err
		}

		// Add the client to the map
		p.clients[addressStr] = client
		fmt.Printf("New connection from %s\n", address)
	}

	return client, nil
}

func (p *Proxy) closeClient(address net.Addr) {
	addressStr := address.String()
	client, ok := p.clients[addressStr]
	if ok {
		client.Close()
		delete(p.clients, addressStr)
		fmt.Printf("Closed connection from %s\n", address)
	}
}

func (p *Proxy) handleMessage(address net.Addr) MessageHandlerFunc {
	return func(buf []byte, size int) {
		p.server.WriteTo(buf[:size], address)
	}
}

func (p *Proxy) handleError(address net.Addr) ErrorHandlerFunc {
	return func(err error) {
		fmt.Printf("An error occurred with client %s, error: %v\n", address, err)
		p.closeClient(address)
	}
}

func (p *Proxy) handleEvent(eventType uint8, handler ProxyEvent) ([]byte, bool) {
	middlewares, ok := p.middlewares[eventType]
	if !ok {
		return nil, true
	}

	for _, middleware := range middlewares {
		buffer, ok := middleware(handler)
		if !ok {
			return buffer, false
		}
	}

	return nil, true
}