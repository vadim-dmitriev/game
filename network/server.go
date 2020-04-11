package network

import "net"

func (n *Network) runAsServer() {
	listener, err := net.Listen("tcp", n.URI)
	if err != nil {
		// TODO: убрать панику
		panic(err)
	}
	defer listener.Close()

	for {
		_, err := listener.Accept()
		if err != nil {
			// TODO: убрать панику
			panic(err)
		}

		// n.clients.Store(conn.RemoteAddr().String(), object{})
		// n.clients

		// go n.handleConn(conn)
	}
}

func (n *Network) handleConn(conn net.Conn) {
	defer conn.Close()
}
