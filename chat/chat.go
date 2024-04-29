package chat

import (
	"encoding/json"
	// "fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

type Client struct {
	hub  *Hub
	send chan Message
	conn *websocket.Conn
}

func NewClient(hub *Hub, conn *websocket.Conn) *Client {

	return &Client{
		hub:  hub,
		send: make(chan Message),
		conn: conn,
	}

}

func (c *Client) Write() {
	ticker := time.NewTicker(pingPeriod)
	defer func(){
		c.conn.Close()
		c.hub.Unregister <- c
		ticker.Stop()
	}()


	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				log.Println("Connection closed")
				c.hub.Unregister <- c
				return
			} else {
				err := c.conn.WriteJSON(message)
				if err != nil {
					log.Println("Error: ", err)
					break
				}
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			log.Println("ping sent")
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err!= nil {
				return
			}
		}

	}

}

func (c *Client) Read() {
	
	defer func(){
		c.hub.Unregister <- c
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		msgtype, msg, err := c.conn.ReadMessage()
		m:= new(Message)
		
		if err != nil {
			log.Println("Error reading message")
			c.hub.Unregister <- c
			break
		} else if msgtype != 1 {
			log.Println("Message other than Textmessag received")
		} else {
			// js, err := json.Marshal(msg)
			// if err!=nil{
			// 	log.Println("Error in parsing to json")
			// 	return
			// }
			// log.Printf("Json is %v \n", js)
			json.Unmarshal(msg, m)
			log.Println("Received message ", m)
		}

		c.hub.Message <- *m
	}

}
