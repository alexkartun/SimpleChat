package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var mu sync.Mutex								// mutex for synchronization reader/writers from/to the map
var clients = make(map[*websocket.Conn]bool)	// connected clients
var broadcast = make (chan Message)				// broadcast channel


// configure the upgrader
var upgrader = websocket.Upgrader {
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// define our message object
type Message struct {
	Email 	 string    		`json:"email"`
	Username string    		`json:"username"`
	Message  string 		`json:"message"`
}


func main() {
	// create a simple file server
	fs := http.FileServer(http.Dir("../public"))
	http.Handle("/", fs)

	// configure websocket route
	http.HandleFunc("/ws", HandleConnections)

	// start listening for incoming chat messages
	go HandleMessages()

	// start the server on localhost port 8080 and log any errors
	log.Println("http server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: %v", err)
	}
}

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	// upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Upgrade: %v", err)
	}

	// make sure we close the connections when the function returns
	defer ws.Close()

	// register new client
	mu.Lock()
	clients[ws] = true
	mu.Unlock()

	for {
		var msg Message

		// read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("ReadJSON: %v", err)
			mu.Lock()
			delete (clients, ws)
			mu.Unlock()
			break
		}

		// send  the newly recieved message to the broadcast channel
		broadcast <- msg
	}
}

func HandleMessages() {
	for {
		// grab the next message from the broadcast channel
		msg := <-broadcast

		// send it out to every client that is currently connected
		mu.Lock()
		for client := range clients  {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("WriteJSON: %v", err)
			}
		}
		mu.Unlock()
	}
}