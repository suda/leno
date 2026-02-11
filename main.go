package main

import (
	"bufio"
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/websocket"
)

//go:embed all:public
var publicFS embed.FS

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type hub struct {
	mu      sync.Mutex
	clients map[*websocket.Conn]struct{}
}

func newHub() *hub {
	return &hub{clients: make(map[*websocket.Conn]struct{})}
}

func (h *hub) add(c *websocket.Conn) {
	h.mu.Lock()
	h.clients[c] = struct{}{}
	h.mu.Unlock()
}

func (h *hub) remove(c *websocket.Conn) {
	h.mu.Lock()
	delete(h.clients, c)
	h.mu.Unlock()
	c.Close()
}

func (h *hub) broadcast(line string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for c := range h.clients {
		if err := c.WriteMessage(websocket.TextMessage, []byte(line)); err != nil {
			delete(h.clients, c)
			c.Close()
		}
	}
}

func main() {
	port := os.Getenv("LENO_PORT")
	if port == "" {
		port = "3000"
	}

	h := newHub()

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			h.broadcast(scanner.Text())
		}
	}()

	sub, err := fs.Sub(publicFS, "public")
	if err != nil {
		log.Fatalf("embed error: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("ws upgrade error: %v", err)
			return
		}
		h.add(conn)
		go func() {
			defer h.remove(conn)
			for {
				if _, _, err := conn.ReadMessage(); err != nil {
					break
				}
			}
		}()
	})
	mux.Handle("/", http.FileServer(http.FS(sub)))

	log.Printf("Leno running on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
