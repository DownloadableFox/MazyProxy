package main

func main() {
	const (
		listenAddr   = ":1337"
		redirectAddr = "192.168.1.101:1337"
	)

	proxy, err := NewProxy(listenAddr, redirectAddr)
	if err != nil {
		panic(err)
	}

	if err := proxy.Serve(); err != nil {
		panic(err)
	}
}
