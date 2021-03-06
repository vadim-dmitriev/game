package main

import (
	"encoding/json"
	"net"
)

func (n *Network) runAsServer() {
	listener, err := net.Listen("tcp", n.URI)
	if err != nil {
		// TODO: убрать панику
		panic(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			// TODO: убрать панику
			panic(err)
		}

		newUser := user{}
		n.Users = append(n.Users, newUser)

		go n.handleConn(conn, &newUser)
	}
}

func (n *Network) handleConn(conn net.Conn, u *user) {
	defer conn.Close()

	username := readUsername(conn)
	u.Username = username

	type msg []user
	for {
		n.mu.RLock()
		for _, user := range n.Users {
			if user.Username == username {
				continue
			}

		}
		n.mu.RUnlock()

		// json.NewEncoder(conn).Encode()

	}
}

func readUsername(conn net.Conn) string {
	type Username struct {
		Username string
	}

	uname := Username{}
	if err := json.NewDecoder(conn).Decode(&uname); err != nil {
		panic("can`t unmarshal username")
	}

	return uname.Username
}
