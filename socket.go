package rocket

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type socket struct {
	route string
	h     func(string)
}

func Socket(path string, handle func(string)) *socket {
	return &socket{
		route: path,
		h:     handle,
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (s *socket) handle(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		s.h(string(p))
		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
	}
}

func (s *socket) GetRoute() string {
	return s.route
}

func (s *socket) Method() string { return "GET" }

func (s *socket) WildcardIndex(i int) {}
