package main

import (
	"fmt"
	"net"
)

const (
	TO_CLIENT = 0
	TO_SERVER = 1
)

type ProxyMiddleware func(*Proxy, Client, []byte) bool
type ClientGenerator func() (Client, error)

type Proxy struct {
	server          Server
	clients         map[string]Client
	middlewares     map[uint8][]ProxyMiddleware
	clientGenerator ClientGenerator
}

func NewProxy(server Server, clientGenerator ClientGenerator) *Proxy {
	return &Proxy{
		server:          server,
		clients:         make(map[string]Client),
		middlewares:     make(map[uint8][]ProxyMiddleware),
		clientGenerator: clientGenerator,
	}
}

func (p *Proxy) Serve() {
	fmt.Printf("Proxy started on %s\n", p.server.GetAddr())
	p.server.Listen(p.handleServerMessage(), p.handleServerError())
}

func (p *Proxy) Close() {
	p.server.Close()
	for _, client := range p.clients {
		client.Close()
	}

}

func (p *Proxy) GetClient(address net.Addr) (Client, error) {
	addressStr := address.String()
	client, ok := p.clients[addressStr]

	if !ok {
		// Create a new newclient
		newclient, err := p.clientGenerator()
		if err != nil {
			return nil, err
		}

		client = newclient

		// Connect the client
		go client.Connect(p.handleClientMessage(address), p.handleClientError(address))

		// Add the client to the map
		p.clients[addressStr] = client
		fmt.Printf("New connection from %s\n", address)
	}

	return client, nil
}

func (p *Proxy) CloseClient(address net.Addr) {
	addressStr := address.String()
	client, ok := p.clients[addressStr]
	if ok {
		client.Close()
		delete(p.clients, addressStr)
		fmt.Printf("Closed connection from %s\n", address)
	}
}

func (p *Proxy) AddMiddleware(direction uint8, middleware ProxyMiddleware) {
	p.middlewares[direction] = append(p.middlewares[direction], middleware)
}

// Handlers
func (p *Proxy) executeMiddlewares(direction uint8, client Client, buf []byte) bool {
	for _, middleware := range p.middlewares[direction] {
		if !middleware(p, client, buf) {
			return false
		}
	}

	return true
}

func (p *Proxy) handleServerMessage() ServerMessageHandler {
	return func(address net.Addr, buf []byte) {
		client, err := p.GetClient(address)
		if err != nil {
			fmt.Printf("An error occurred with client %s, error: %v\n", address, err)
			return
		}

		// Execute middlewares
		if !p.executeMiddlewares(TO_SERVER, client, buf) {
			return
		}

		client.Send(buf)
	}
}

func (p *Proxy) handleServerError() ServerErrorHandler {
	return func(err error) {
		fmt.Printf("An error occurred with server, error: %v\n", err)
		p.Close()
	}
}

func (p *Proxy) handleClientMessage(address net.Addr) ClientMessageHandler {
	return func(buf []byte) {
		client, err := p.GetClient(address)
		if err != nil {
			fmt.Printf("An error occurred with client %s, error: %v\n", address, err)
			return
		}

		// Execute middlewares
		if !p.executeMiddlewares(TO_CLIENT, client, buf) {
			return
		}

		p.server.Send(address, buf)
	}
}

func (p *Proxy) handleClientError(address net.Addr) ClientErrorHandler {
	return func(err error) {
		fmt.Printf("An error occurred with client %s, error: %v\n", address, err)
		p.CloseClient(address)
	}
}
