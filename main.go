package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

const shutdownSentinel = "\n__BBLOG_SHUTDOWN__"

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

func translateLine(line string, stripEscape bool) string {
	if stripEscape {
		return stripAnsi(line) + "\n"
	}

	segments := textToStyledSegments(line)

	return segmentsToPayload(segments)
}

func main() {
	var host string
	var port int
	var stripEscape bool
	var version bool

	flag.StringVar(&host, "H", "0.0.0.0", "HTTP listen host")
	flag.IntVar(&port, "P", 8088, "HTTP listen port")
	flag.BoolVar(&stripEscape, "S", false, "Strip ANSI escape codes")
	flag.BoolVar(&version, "v", false, "Show version and exit")

	flag.Parse()

	if version {
		fmt.Println("bblog version 0.2.0")
		return
	}

	h := newHub()
	addr := fmt.Sprintf("%s:%d", host, port)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	server := &http.Server{Addr: addr}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "streaming unsupported", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("X-Accel-Buffering", "no")

		client := newClient(128)
		h.add(client)
		defer h.remove(client)

		log.Printf("client connected: %s", r.RemoteAddr)
		defer log.Printf("client disconnected: %s", r.RemoteAddr)

		fmt.Fprint(w, ": connected\n\n")
		flusher.Flush()

		for {
			select {
			case <-r.Context().Done():
				return
			case <-client.closed:
				return
			case msg := <-client.ch:
				msg = strings.ReplaceAll(msg, "\n", "\ndata: ")
				fmt.Fprintf(w, "data: %s\n\n", msg)
				flusher.Flush()
			}
		}
	})

	go func() {
		scanner := bufio.NewScanner(os.Stdin)

		for scanner.Scan() {
			raw := scanner.Text()
			msg := translateLine(raw, stripEscape)
			h.broadcast(msg)
		}

		if err := scanner.Err(); err != nil {
			log.Printf("stdin read error: %v", err)
		}

		log.Printf("stdin closed, shutting down")
		stop()
	}()

	go func() {
		<-ctx.Done()

		h.broadcast(shutdownSentinel)
		time.Sleep(100 * time.Millisecond)

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("server shutdown error: %v", err)
		}

		h.closeAll()
	}()

	log.Printf("bblog listening on %s", addr)

	server.Handler = http.DefaultServeMux
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
