package main

import "sync"

const clientBufferSize = 128

type client struct {
	ch     chan string
	closed chan struct{}
	once   sync.Once
}

func newClient(buffer int) *client {
	return &client{
		ch:     make(chan string, buffer),
		closed: make(chan struct{}),
	}
}

func (c *client) close() {
	c.once.Do(func() {
		close(c.closed)
	})
}
