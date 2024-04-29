package main

import (
	"github.com/gorilla/websocket"
	"github.com/rhydberg/chat-app/chat"
	"log"
	"net/http"
)

var allowedOrigin = "http://localhost:3000"

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     checkOrigin, 
}

func checkOrigin(r *http.Request) bool{
	return true
	origin := r.Header.Get("Origin")
	log.Println("Origin is ", origin)
	return origin == allowedOrigin
}


var hub *chat.Hub = chat.NewHub()

// func echo(w http.ResponseWriter, r *http.Request) {
// 	conn, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		log.Fatal("Upgrading to websocket failed")
// 	}
// 	go func() {

// 		defer conn.Close()

// 		for {
// 			msgtype, msg, err := conn.ReadMessage()
// 			if err != nil {
// 				fmt.Println("Could not Read message")
// 			}

// 			fmt.Printf("Received %s of msgtype %v\n", msg, msgtype)
// 			response := fmt.Sprintf("%s%s", "Echoing : ", msg)
// 			if err := conn.WriteMessage(msgtype, []byte(response)); err != nil {
// 				fmt.Printf("Error %x\n", err)
// 			}

// 		}
// 	}()

// }

// func hello(w http.ResponseWriter, r *http.Request) {
// 	w.Write([]byte("Hello, World!"))
// }

func room(w http.ResponseWriter, r *http.Request) {
	// fmt.Println("in room")

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("Error in upgrading websocket")
	} else {
		log.Println("Upgraded successfully")
	}

	client := chat.NewClient(hub, conn)
	log.Println("After making client")
	hub.Register <- client

	
	go client.Read()
	go client.Write()
}

func main() {

	router := http.NewServeMux()
	// router.HandleFunc("/hello", hello)
	// router.HandleFunc("/echo", echo)
	router.HandleFunc("/room", room)

	go hub.Run()

	log.Fatal(http.ListenAndServe(":8080", router))

}
