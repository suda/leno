package main

import (
	"bufio"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/suda/leno/parsers"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

//go:embed all:public
var publicFS embed.FS

type hub struct {
	mu      sync.Mutex
	clients map[chan string]struct{}
}

func newHub() *hub {
	return &hub{clients: make(map[chan string]struct{})}
}

func (h *hub) add(ch chan string) {
	h.mu.Lock()
	h.clients[ch] = struct{}{}
	h.mu.Unlock()
}

func (h *hub) remove(ch chan string) {
	h.mu.Lock()
	delete(h.clients, ch)
	h.mu.Unlock()
	close(ch)
}

func (h *hub) broadcast(line string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for ch := range h.clients {
		select {
		case ch <- line:
		default:
		}
	}
}

func main() {
	port := os.Getenv("LENO_PORT")
	if port == "" {
		port = "3000"
	}

	logFormat := flag.String("log-format", "", "Log format to parse (nginx, logfmt)")
	flag.Parse()

	h := newHub()

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			line := scanner.Text()
			switch *logFormat {
			case "nginx":
				if parsed, ok := parsers.ParseNginx(line); ok {
					line = parsed
				}
			case "logfmt":
				if parsed, ok := parsers.ParseLogfmt(line); ok {
					line = parsed
				}
			}
			h.broadcast(line)
		}
	}()

	sub, err := fs.Sub(publicFS, "public")
	if err != nil {
		log.Fatalf("embed error: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, version)
	})
	mux.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		ch := make(chan string, 16)
		h.add(ch)
		defer h.remove(ch)

		for {
			select {
			case line, ok := <-ch:
				if !ok {
					return
				}
				fmt.Fprintf(w, "data: %s\n\n", line)
				if f, ok := w.(http.Flusher); ok {
					f.Flush()
				}
			case <-r.Context().Done():
				return
			}
		}
	})
	mux.Handle("/", http.FileServer(http.FS(sub)))

	log.Printf("Leno running on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
