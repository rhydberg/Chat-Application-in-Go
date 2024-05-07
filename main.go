package main

import (
	"log"
	"net/http"
	"slices"

	"github.com/gorilla/websocket"
	"github.com/rhydberg/chat-app/chat"

	"github.com/rhydberg/chat-app/db"
	"github.com/rhydberg/chat-app/models"
	// "gorm.io/driver/postgres"
	// "gorm.io/gorm"
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


var hub *chat.Hub 


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


	
	recentMessages := hub.GetRecentMessages(5)
	slices.Reverse(recentMessages)
	log.Println("Recent messages are ", len(recentMessages))
	for _,msg := range recentMessages{
		conn.WriteJSON(msg)
	}

	
	go client.Read()
	go client.Write()
}

func main() {

	// var messages []models.Message

	db := db.Init()
	db.AutoMigrate(&models.Message{})

	// db.Limit(2).Order("ID desc").Find(&messages)
	// log.Println(messages)
	hub = chat.NewHub(db)

	router := http.NewServeMux()
	router.HandleFunc("/room", room)

	go hub.Run()



	log.Fatal(http.ListenAndServe(":8080", router))

}
