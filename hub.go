package main

import "sync"

type hub struct {
	mu      sync.RWMutex
	clients map[*client]struct{}
}

func newHub() *hub {
	return &hub{
		clients: make(map[*client]struct{}),
	}
}

func (h *hub) add(c *client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.clients[c] = struct{}{}
}

func (h *hub) remove(c *client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.clients[c]; !ok {
		return
	}

	delete(h.clients, c)
	c.close()
}

func (h *hub) closeAll() {
	h.mu.Lock()
	defer h.mu.Unlock()

	for c := range h.clients {
		c.close()
		delete(h.clients, c)
	}
}

func (h *hub) broadcast(payload string) {
	h.mu.RLock()

	clients := make([]*client, 0, len(h.clients))
	for c := range h.clients {
		clients = append(clients, c)
	}

	h.mu.RUnlock()

	for _, c := range clients {
		select {
		case c.ch <- payload:
		case <-c.closed:
			continue
		}
	}
}
