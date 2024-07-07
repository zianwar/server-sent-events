package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type client struct {
	id        string
	messageCh chan string
}

type Broker struct {
	clients      map[string]client
	addClient    chan client
	removeClient chan string
	messages     chan string
	mu           sync.Mutex
}

func NewBroker() *Broker {
	return &Broker{
		clients:      make(map[string]client),
		addClient:    make(chan client),
		removeClient: make(chan string),
		messages:     make(chan string, 1000),
	}
}

func (b *Broker) Start() {
	for {
		select {
		case client := <-b.addClient:
			b.mu.Lock()
			b.clients[client.id] = client
			b.mu.Unlock()
		case clientId := <-b.removeClient:
			b.mu.Lock()
			delete(b.clients, clientId)
			b.mu.Unlock()
		// Broadcast message to all clients
		case message := <-b.messages:
			for _, c := range b.clients {
				c.messageCh <- message
			}
		}
	}
}

func (b *Broker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	// We got a new connection, register the client
	clientId := r.URL.Query().Get("client_id")
	if clientId == "" {
		http.Error(w, "client_id is missing from the query params", http.StatusBadRequest)
		return
	}

	log.Printf("Client connected: %s\n", clientId)

	// Add the client
	messageCh := make(chan string)
	b.addClient <- client{id: clientId, messageCh: messageCh}

	// Unregister the client when this method exists
	defer func() {
		log.Printf("client disconnected: %s\n", clientId)
		b.removeClient <- clientId
	}()

	// Setup response headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Send welcome message
	fmt.Fprintf(w, "data: Welcome %s!\n\n", clientId)
	flusher.Flush()

	// Listen for messages sent to this client's channel and write it to the connection.
	for {
		select {
		case m := <-messageCh:
			// The response content must start with "data:"
			fmt.Fprintf(w, "data: %s\n\n", m)
			flusher.Flush()
		case <-r.Context().Done():
			return
		}
	}
}

func (b *Broker) SendMessage(m string) {
	b.messages <- m
}

func main() {
	broker := NewBroker()
	go broker.Start()

	// Simulate broadcasting events to all clients
	go func() {
		for {
			time.Sleep(2 * time.Second)
			broker.SendMessage(fmt.Sprintf("%v\n", time.Now().Format(time.RFC1123)))
		}
	}()

	// Handle the server-sent events endpoint
	http.Handle("/events", broker)

	// Serve HTML page that will connect to the SSE stream
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	log.Fatal(http.ListenAndServe(":9000", nil))
}
