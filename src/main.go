package main

const (
	listenAddr   = "localhost:8080"
	redirectAddr = "localhost:8081"
)

func main() {
	// Create server
	server, err := NewUDPServer(listenAddr)
	if err != nil {
		panic(err)
	}

	// Create proxy
	proxy := NewProxy(server, func() (Client, error) {
		return NewUDPClient(redirectAddr)
	})
	defer proxy.Close()

	// Register middlewares
	RegisterMiddlewares(proxy)

	// Start proxy
	proxy.Serve()
}
