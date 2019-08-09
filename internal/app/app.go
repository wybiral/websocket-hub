package app

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/wybiral/websocket-hub/pkg/hub"
)

// WebSocket upgrader
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow requests from any origin
		return true
	},
}

// Command is the JSON object used for all pub/sub commands.
type Command struct {
	Type  string `json:"type"`
	Topic string `json:"topic"`
}

// App is the main Pub/Sub application containing the hub and the server.
type App struct {
	Hub    *hub.Hub
	Server *http.Server
}

// New returns a new App with server ready to listen on addr and serve optional
// public directory dir.
func New(addr, dir string) *App {
	a := &App{}
	a.Hub = hub.New()
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/socket", a.socketHandler).Methods("GET")
	if dir != "" {
		fs := http.FileServer(http.Dir(dir))
		r.PathPrefix("/").Handler(fs)
	}
	a.Server = &http.Server{
		Addr:    addr,
		Handler: r,
	}
	return a
}

// socketHandler handles all WebSocket requests and creates client chan.
func (a *App) socketHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer ws.Close()
	ch := make(chan []byte)
	defer func() {
		a.Hub.UnsubscribeAll(ch)
		close(ch)
	}()
	go readLoop(a.Hub, ws, ch)
	writeLoop(ws, ch)
}

// readLoop reads from WebSocket and processes commands.
func readLoop(h *hub.Hub, ws *websocket.Conn, ch chan []byte) {
	var cmd Command
	for {
		_, b, err := ws.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		err = json.Unmarshal(b, &cmd)
		if err != nil {
			log.Println(err)
			return
		}
		switch cmd.Type {
		case "publish":
			h.Publish(cmd.Topic, b)
		case "subscribe":
			h.Subscribe(cmd.Topic, ch)
		case "unsubscribe":
			h.Unsubscribe(cmd.Topic, ch)
		}
	}
}

// writeLoop writes everything from client chan out to client WebSocket.
func writeLoop(ws *websocket.Conn, ch chan []byte) {
	for b := range ch {
		err := ws.WriteMessage(websocket.TextMessage, b)
		if err != nil {
			log.Println(err)
			return
		}
	}
}
