package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(req *http.Request) bool { return true },
}

var connections map[int]*websocket.Conn
var lastUsedConn int

func handleRoot(resp http.ResponseWriter, req *http.Request) {
	fmt.Fprint(resp, "oioiasdasd")
}

func broadcastMessage(message []byte) {
	for _, v := range connections {
		v.WriteMessage(websocket.TextMessage, message)
	}
}

func handleChat(resp http.ResponseWriter, req *http.Request) {
	conn, err := upgrader.Upgrade(resp, req, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	thisConnID := lastUsedConn
	lastUsedConn++
	connections[thisConnID] = conn

	defer conn.Close()
	defer delete(connections, thisConnID)

	for {
		_, message, err := conn.ReadMessage()

		if err != nil {
			log.Print("read:", err)
			break
		}

		log.Printf("recv: %s", message)
		cmd := strings.Fields(string(message))

		if cmd[0] == "enter" {
			broadcastMessage([]byte("::SERVER:: " + cmd[1] + " entrou!"))
		} else if cmd[0] == "change" {
			broadcastMessage([]byte("::SERVER:: " + cmd[1] + " mudou nick para " + cmd[3] + "!"))
		} else if cmd[0] == "msg" {
			broadcastMessage([]byte(">> " + cmd[1] + " >> " + strings.TrimPrefix(string(message), "msg "+cmd[1]+" ")))
		}

		if err != nil {
			log.Print("write:", err)
			break
		}

	}
}

func main() {

	fmt.Print("Startint Server!")

	connections = make(map[int]*websocket.Conn)

	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/chat", handleChat)

	port := os.Getenv("PORT")

	err := http.ListenAndServe(":"+port, nil)

	if err != nil {
		log.Fatal("Serven listen error: ", err)
	}
}
