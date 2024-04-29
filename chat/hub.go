package chat

import (
	// "fmt"
	"log"
	"sync"
)


type ClientList map[*Client]bool

type Message struct {
	Type int
	Sender string
	Content string
}

type Hub struct {
	Clients    ClientList
	Register   chan *Client
	Message    chan Message
	Unregister chan *Client

	sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(ClientList),
		Message:    make(chan Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (h *Hub) RegisterClient(client *Client){
	// h.Lock()
	// defer h.Unlock()

	if h.Clients[client]{
		log.Println("Client already registered")
	} else{
		h.Clients[client]= true
		log.Println("New Client registered")
		log.Println("Total number of clients: ", len(h.Clients))
	}
}

func (h *Hub) UnregisterClient(client *Client){
	// h.Lock()
	// defer h.Unlock()

	if _, ok:= h.Clients[client]; ok{
		client.conn.Close()
		close(client.send)
		delete(h.Clients, client)

		log.Println("Unregistered Client")
		log.Println("Total number of clients: ", len(h.Clients))
	}
}

func (h *Hub) Run() {
	for {
		select {
		case newClient := <-h.Register:

			log.Println("In register ")
			h.RegisterClient(newClient)

			// fmt.Println("new cleint is ", newClient)
			// if !h.Clients[newClient] {
			// 	h.Clients[newClient] = true
			// 	fmt.Println("registered new client")
			// 	fmt.Printf("number of clients: %d\n", len(h.Clients))
			// }

		case msg := <-h.Message:
			log.Println("Hub received message ", msg)
			h.handleMessage(msg)

		case client := <-h.Unregister:

			delete(h.Clients, client)
			log.Println("unregistered client")
			log.Printf("number of clients: %d\n", len(h.Clients))

		}

	}
}

func (h *Hub) handleMessage(msg Message) {
	m:=Message{}
	if m == msg{
		return

	}
	for client := range h.Clients {
		select {
		case client.send <- msg:
			log.Printf("sent msg to client")
		default:
			h.Unregister<-client
		}

	}
}
