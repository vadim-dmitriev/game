package main

import (
	"net"
	"sync"
	"time"
)

const (
	// server режим работы в качестве сервера и клиета одновременно
	server = iota

	// client - режим работы в качестве клиента
	client
)

const (
	pingTimeout = 5 * time.Second
)

// Network TODO
type Network struct {
	URI   string
	Users []user
	mu    *sync.RWMutex
}

// NewNetwork создает новый объект сети.
func NewNetwork(uri string) *Network {
	this := &Network{
		URI:   uri,
		Users: make([]user, 0, 3),
		mu:    &sync.RWMutex{},
	}

	return this
}

// Run определяет в каком из двух режимов будет работать сеть и запускает
// сетевое взаимодействие.ы
func (n *Network) Run() {
	// в каком режиме будет запущен мультиплеер
	switch whichMode(n) {

	case server:
		go n.runAsServer()

	case client:
		go n.runAsClient()

	}
}

// whichMode пингует адрес.
// 		если пинг вернул ошибку => значит на том конце нет ответа => значит я сервер
// 		если пинг без ошибок => значит на том конце кто-то есть => я клиент
func whichMode(n *Network) int {
	_, err := net.DialTimeout("tcp", n.URI, pingTimeout)
	if err != nil {
		return server
	}

	return client
}
