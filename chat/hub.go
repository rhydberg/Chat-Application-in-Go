package chat

import (
	// "fmt"
	"log"
	"sync"
	"time"

	
	"gorm.io/gorm"
)


type ClientList map[*Client]bool

type Message struct {
	Type int
	Sender string
	Content string
	CreatedAt time.Time
}

type Hub struct {
	Clients    ClientList
	Register   chan *Client
	Message    chan Message
	Unregister chan *Client
	DB		 *gorm.DB

	sync.RWMutex
}

func NewHub(db *gorm.DB) *Hub {
	return &Hub{
		Clients:    make(ClientList),
		Message:    make(chan Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		DB: db,
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
		client.Conn.Close()
		close(client.Send)
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

func (h *Hub) GetRecentMessages(n int) []Message {


	var messages []Message
	h.DB.Limit(n).Order("ID desc").Find(&messages)
	return messages
}

func (h *Hub) handleMessage(msg Message) {
	m:=Message{}
	if m == msg{ //checking for empty message
		return

	}
	msg.CreatedAt = time.Now()
	// msg = Message{

	h.DB.Create(&msg)

	for client := range h.Clients {
		select {
		case client.Send <- msg:
			log.Printf("sent msg to client")
		default:
			h.Unregister<-client
		}

	}
}
