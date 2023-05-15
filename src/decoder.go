package main

import (
	"fmt"
)

const (
	PACKET_HEARTBEAT = 0x3c
	PACKET_EMOJI     = 0x45
)

func RegisterMiddlewares(proxy *Proxy) {
	proxy.AddMiddleware(TO_SERVER, decodeToServerMessage)
	proxy.AddMiddleware(TO_CLIENT, decodeToClientMessage)
}

func decodeToServerMessage(p *Proxy, client Client, buff []byte) bool {
	fmt.Printf("[%s] -> %v\n", client.GetAddr(), NetworkDecrypt(buff))
	return true
}

func decodeToClientMessage(p *Proxy, client Client, buff []byte) bool {
	fmt.Printf("[%s] <- %v\n", client.GetAddr(), NetworkDecrypt(buff))
	return true
}
