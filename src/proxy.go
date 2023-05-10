package main

import (
	"fmt"
	"net"
)

type Proxy struct {
	listenAddr   string
	redirectAddr string
	server       net.PacketConn
	clients      map[string]*UDPClient
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
	}, nil
}

func (p *Proxy) Serve() error {
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

		// Send the message to the client
		client.Send(buffer[:size])
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

func (p *Proxy) Close() {
	p.server.Close()
	for _, client := range p.clients {
		client.Close()
	}
}
