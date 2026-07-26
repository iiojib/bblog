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
	"syscall"
	"time"
)

const (
	shutdownSentinel = "\n__BBLOG_SHUTDOWN__"
	defaultHost      = "0.0.0.0"
	defaultPort      = 8088
	appVersion       = "0.3.0"
	shutdownTimeout  = 5 * time.Second
	shutdownWaitTime = 100 * time.Millisecond
)

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

	flag.StringVar(&host, "H", defaultHost, "HTTP listen host")
	flag.IntVar(&port, "P", defaultPort, "HTTP listen port")
	flag.BoolVar(&stripEscape, "S", false, "Strip ANSI escape codes")
	flag.BoolVar(&version, "v", false, "Show version and exit")

	flag.Parse()

	if version {
		fmt.Println("bblog version " + appVersion)
		return
	}

	var inputPath string
	if flag.NArg() > 0 {
		inputPath = flag.Arg(0)
	}

	h := newHub()
	addr := fmt.Sprintf("%s:%d", host, port)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	src, err := newSource(inputPath)
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer src.Close()

	stopClose := context.AfterFunc(ctx, func() {
		src.Close()
	})
	defer stopClose()

	server := &http.Server{Addr: addr}

	http.HandleFunc("OPTIONS /", handleOptions)
	http.HandleFunc("GET /", handleSSE(h))

	go func() {
		scanner := bufio.NewScanner(src)

		for scanner.Scan() {
			raw := scanner.Text()
			msg := translateLine(raw, stripEscape)
			h.broadcast(msg)
		}

		if err := scanner.Err(); err != nil && ctx.Err() == nil {
			log.Printf("%s read error: %v", src, err)
		}

		log.Printf("%s closed, shutting down", src)
		stop()
	}()

	go func() {
		<-ctx.Done()

		h.broadcast(shutdownSentinel)
		time.Sleep(shutdownWaitTime)
		h.closeAll()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("server shutdown error: %v", err)
		}
	}()

	log.Printf("bblog listening on %s", addr)

	server.Handler = http.DefaultServeMux
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
